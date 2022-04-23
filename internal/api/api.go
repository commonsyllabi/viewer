package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	zero "github.com/commonsyllabi/viewer/internal/logger"
	"github.com/commonsyllabi/viewer/pkg/commoncartridge"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Port       string `yaml:"port"`
	UploadsDir string `yaml:"uploadsDir"`
	TmpDir     string `yaml:"tmpDir"`
}

func (cc *Config) loadConfig(path string) error {
	var c Config
	cwd, _ := os.Getwd()
	content, err := ioutil.ReadFile(filepath.Join(cwd, path))
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &c)
	zero.Log.Debug().Msgf("%+v", c)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) defaults() {
	c.Port = "2046"
	c.UploadsDir = "uploads"
	c.TmpDir = "tmp"
}

var conf Config

func StartServer() {
	err := conf.loadConfig("internal/api/config.yml")
	if err != nil || conf.Port == "" {
		zero.Log.Warn().Msgf("error loading config: %v", err)
		conf.defaults()
	}

	router := setupRouter()

	server := &http.Server{
		Addr:         ":" + conf.Port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server.ListenAndServe()
}

func setupRouter() *gin.Engine {
	router := gin.New()

	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	router.MaxMultipartMemory = 16 << 20 // 16 MiB for uploads
	router.Use(gin.Recovery())

	router.GET("/ping", handlePing)
	router.POST("/upload", handleUpload)
	router.GET("/resource/:id", handleResource)
	router.GET("/file/:id", handleFile)

	return router
}

func handlePing(c *gin.Context) {
	c.String(200, "pong")
}

func handleFile(c *gin.Context) {

	id := c.Param("id")
	cartridge := c.Request.FormValue("cartridge")

	zero.Log.Info().Msgf("handleFile id: %v cartridge %v", id, cartridge)

	inputFile := filepath.Join(conf.UploadsDir, cartridge)
	cc, err := commoncartridge.Load(inputFile)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error loading CC from disk: %v", err)
		return
	}

	file, err := cc.FindFile(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error finding finding file in CC: %v", err)
		return
	}

	info, err := file.Stat()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error getting file info: %v", err)
		return
	}

	err = os.MkdirAll(conf.TmpDir, os.ModePerm)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	path := filepath.Join(conf.TmpDir, info.Name())
	dst, err := os.Create(path)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error creating dest tmp file: %v", err)
		return
	}

	_, err = io.Copy(dst, file)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error writing file to tmp: %v", err)
		return
	}

	//-- handle doc to pdf conversion
	ext := filepath.Ext(info.Name())
	match, err := regexp.Match(`(doc|docx|odt)`, []byte(ext))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error parsing file extension: %v", err)
		return
	}

	if match {
		libreoffice, err := exec.LookPath("libreoffice")

		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			zero.Log.Error().Msgf("error finding libreoffice: %v", err)
			return
		}

		cmd := exec.Command(libreoffice, "--headless", "--convert-to", "pdf", "--outdir", conf.TmpDir, path)

		err = cmd.Run()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			zero.Log.Error().Msgf("error converting file to pdf: %v", err)
			return
		}

		path = strings.TrimSuffix(path, filepath.Ext(path)) + ".pdf"

	}

	c.JSON(http.StatusOK, gin.H{"path": path})
}

func handleResource(c *gin.Context) {

	id := c.Param("id")
	cartridge := c.Query("cartridge")
	zero.Log.Info().Msgf("GET handleResource id: %v cartridge %v", id, cartridge)

	inputFile := filepath.Join(conf.UploadsDir, cartridge)
	cc, err := commoncartridge.Load(inputFile)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error loading CC from disk: %v", err)
		return
	}

	file, err := cc.Find(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error finding resource in CC: %v", err)
		return
	}

	data, err := json.Marshal(file)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msg("error marshalling to json")
		return
	}

	c.JSON(http.StatusOK, string(data))
}

func handleUpload(c *gin.Context) {

	file, err := c.FormFile("cartridge")
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		zero.Log.Warn().Msgf("cannot upload cartridge file: %v", err)
		return
	}

	err = os.MkdirAll(conf.UploadsDir, os.ModePerm)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	dst := filepath.Join(conf.UploadsDir, file.Filename)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error saving CC to filesystem: %v", err)
		return
	}

	cc, err := commoncartridge.Load(dst)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error loading CC from filesystem: %v", err)
		return
	}

	obj, err := cc.MarshalJSON()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error parsing manifest into JSON: %v", err)
		return
	}

	fi, err := cc.Items()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error getting resources: %v", err)
		return
	}
	sfi, err := json.Marshal(fi)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error getting resources: %v", err)
		return
	}

	fr, err := cc.Resources()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error getting resources: %v", err)
		return
	}
	sfr, err := json.Marshal(fr)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error getting resources: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      string(obj),
		"items":     string(sfi),
		"resources": string(sfr),
	})
}
