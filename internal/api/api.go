// API is the package that takes care of loading the configuratio and setting up the router
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/commonsyllabi/commoncartridge"
	zero "github.com/commonsyllabi/viewer/internal/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var conf Config

// StartServer gets his port and debug in the environment, registers the router, and registers the database closing on exit.
func StartServer(port string, debug bool, c Config) error {
	conf = c

	gin.SetMode(gin.ReleaseMode)

	err := os.MkdirAll(filepath.Join(conf.TmpDir, conf.FilesDir), os.ModePerm)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Join(conf.TmpDir, conf.UploadsDir), os.ModePerm)
	if err != nil {
		return err
	}

	err = copySampleFiles()
	if err != nil {
		return err
	}

	router, err := setupRouter()
	if err != nil {
		return err
	}

	s := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// from https://gist.github.com/ivan3bx/b0f14449803ce5b0aa72afaa1dfc75e1
	go func() {
		zero.Infof("server starting on port %s", port)
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	<-ch // block until signal received

	zero.Info("shutting down...")
	return s.Shutdown(context.Background())
}

// setupRouter registers all route groups
func setupRouter() (*gin.Engine, error) {
	router := gin.New()

	router.Use(cors.Default())

	err := os.MkdirAll(conf.TmpDir, os.ModePerm)
	if err != nil {
		return router, err
	}

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

	cwd, _ := os.Getwd()
	publicPath := filepath.Join(cwd, conf.PublicDir)

	router.Use(static.Serve("/", static.LocalFile(publicPath, false)))

	router.GET("/ping", handlePing)
	router.POST("/parse", handleUpload)
	api := router.Group("/api")
	{

		api.GET("/resource/:id", handleResource)
		api.GET("/file/:id", handleFile)
	}

	router.Use(handleNotFound)

	return router, nil
}

func handlePing(c *gin.Context) {
	c.String(200, "pong")
}

func handleNotFound(c *gin.Context) {
	c.HTML(http.StatusOK, "Error", gin.H{
		"msg": "We couldn't find the requested resource, sorry :(.",
	})
}

// handleFile takes a file ID and a given cartridge as query parameter, and returns a file stream
func handleFile(c *gin.Context) {

	id := c.Param("id")
	cartridge := c.Request.FormValue("cartridge")

	zero.Infof("handleFile id: %v cartridge %v", id, cartridge)

	cartridgePath := filepath.Join(conf.TmpDir, conf.UploadsDir, cartridge)
	cc, err := commoncartridge.Load(cartridgePath)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error loading CC from disk: %v", err)
		return
	}

	file, err := cc.FindFile(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error finding finding file in CC: %v", err)
		return
	}

	//convert to PDF
	info, err := file.Stat()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error getting file info: %v", err)
		return
	}

	ext := filepath.Ext(info.Name())
	match, err := regexp.Match(`(doc|docx|odt)`, []byte(ext))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error parsing file extension: %v", err)
		return
	}

	if match {
		file, err = convertToPDF(file, info)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			zero.Errorf("error converting to PDF: %v", err)
			return
		}
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error reading file into bytes: %v", err)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	mimeType := http.DetectContentType(bytes)
	c.Header("Content-Type", mimeType)
	c.Writer.Write(bytes)
}

// convertToPDF writes the original doc/docx/odt file to disk, then converts it to PDF, and returns the converted file
func convertToPDF(file fs.File, info fs.FileInfo) (fs.File, error) {
	var f fs.File
	libreoffice, err := exec.LookPath("libreoffice")

	if err != nil {
		return f, err
	}

	path := filepath.Join(conf.TmpDir, conf.FilesDir, info.Name())
	dst, err := os.Create(path)
	if err != nil {
		return f, err
	}

	_, err = io.Copy(dst, file)
	if err != nil {
		return f, err
	}

	cmd := exec.Command(libreoffice, "--headless", "--convert-to", "pdf", "--outdir", filepath.Join(conf.TmpDir, conf.FilesDir), path)

	err = cmd.Run()
	if err != nil {
		return file, err
	}

	path = strings.TrimSuffix(path, filepath.Ext(path)) + ".pdf"
	file, err = os.Open(path)

	return file, err
}

func handleResource(c *gin.Context) {

	id := c.Param("id")
	cartridge := c.Query("cartridge")
	zero.Infof("GET handleResource id: %v cartridge %v", id, cartridge)

	inputFile := filepath.Join(conf.TmpDir, conf.UploadsDir, cartridge)
	cc, err := commoncartridge.Load(inputFile)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error loading CC from disk: %v", err)
		return
	}

	file, err := cc.Find(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error finding resource in CC: %v", err)
		return
	}

	data, err := json.Marshal(file)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Error("error marshalling to json")
		return
	}

	c.JSON(http.StatusOK, string(data))
}

// handleUpload expects a Common Cartridge-compliant file, saves it to disk, and creates an IMSCC instance from the file to return manifest, items and resources in JSON format.
func handleUpload(c *gin.Context) {
	file, err := c.FormFile("cartridge")
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		zero.Warnf("cannot upload cartridge file: %v", err)
		return
	}

	err = os.MkdirAll(filepath.Join(conf.TmpDir, conf.UploadsDir), os.ModePerm)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	dst := filepath.Join(conf.TmpDir, conf.UploadsDir, file.Filename)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error saving CC to filesystem: %v", err)
		return
	}

	defer os.Remove(dst)

	cc, err := commoncartridge.Load(dst)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error loading CC from filesystem: %v", err)
		return
	}

	obj, err := cc.Manifest()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error parsing manifest into JSON: %v", err)
		return
	}

	fi, err := cc.Items()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error getting items: %v", err)
		return
	}

	fr, err := cc.Resources()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error getting resources: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      obj,
		"items":     fi,
		"resources": fr,
	})
}

func copySampleFiles() error {
	samples := []string{"test_01.imscc", "OpenMed-English-IMSCC1-3-Canvas-sakai-export.imscc", "Falconer_Liz-Computers_Canvas_Community-941267a5971248daa62d3196014d1e65.zip"}

	for _, v := range samples {
		f, err := os.ReadFile(filepath.Join(os.Getenv("COSYLL_VIEWER_SAMPLES_DIR"), v))

		if err != nil {
			return err
		}

		err = os.WriteFile(filepath.Join(conf.TmpDir, conf.UploadsDir, v), f, 0644)

		if err != nil {
			return err
		}
	}

	return nil
}
