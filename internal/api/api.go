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

	"github.com/commonsyllabi/viewer/internal/api/handlers"
	"github.com/commonsyllabi/viewer/internal/api/models"
	zero "github.com/commonsyllabi/viewer/internal/logger"
	"github.com/commonsyllabi/viewer/pkg/commoncartridge"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func StartServer() error {

	err := conf.load("./internal/api/config.yml")

	if err != nil || conf.port == "" {
		zero.Log.Warn().Msgf("error loading config: %v", err)
		conf.defaults()
	}

	router, err := setupRouter(true)
	if err != nil {
		return err
	}

	s := &http.Server{
		Addr:         ":" + conf.port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// from https://gist.github.com/ivan3bx/b0f14449803ce5b0aa72afaa1dfc75e1
	go func() {
		zero.Log.Info().Msgf("server starting on port %s", conf.port)
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c // block until signal received

	zero.Log.Info().Msg("shutting down...")
	s.Shutdown(context.Background())
	err = models.Shutdown()

	return err
}

func setupRouter(debug bool) (*gin.Engine, error) {
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))

	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	err := os.MkdirAll(conf.tmpDir, os.ModePerm)
	if err != nil {
		return router, err
	}

	err = os.MkdirAll(filepath.Join(conf.tmpDir, conf.filesDir), os.ModePerm)
	if err != nil {
		return router, err
	}

	err = os.MkdirAll(filepath.Join(conf.tmpDir, conf.uploadsDir), os.ModePerm)
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
	publicPath := filepath.Join(cwd, conf.publicDir)

	router.Use(static.Serve("/", static.LocalFile(publicPath, false)))

	router.GET("/ping", handlePing)

	api := router.Group("/api")
	{
		api.POST("/upload", handleUpload)
		api.GET("/resource/:id", handleResource)
		api.GET("/file/:id", handleFile)
	}

	syllabi := router.Group("/syllabi")
	{
		syllabi.GET("/", handlers.AllSyllabi)
		syllabi.POST("/", handlers.NewSyllabus)
		syllabi.POST("/:id", handlers.UpdateSyllabus)
		syllabi.GET("/:id", handlers.GetSyllabus)
		syllabi.DELETE("/:id", handlers.DeleteSyllabus)
	}

	attachments := router.Group("/attachments")
	{
		attachments.GET("/", handlers.AllAttachments)
		attachments.POST("/", handlers.NewAttachment)
		attachments.POST("/:id", handlers.UpdateAttachment)
		attachments.GET("/:id", handlers.GetAttachment)
		attachments.DELETE("/:id", handlers.DeleteAttachment)
	}

	return router, nil
}

func handlePing(c *gin.Context) {
	c.String(200, "pong")
}

func handleFile(c *gin.Context) {

	id := c.Param("id")
	cartridge := c.Request.FormValue("cartridge")

	zero.Log.Info().Msgf("handleFile id: %v cartridge %v", id, cartridge)

	inputFile := filepath.Join(conf.tmpDir, conf.uploadsDir, cartridge)
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

	//convert to PDF
	info, err := file.Stat()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error getting file info: %v", err)
		return
	}

	ext := filepath.Ext(info.Name())
	match, err := regexp.Match(`(doc|docx|odt)`, []byte(ext))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error parsing file extension: %v", err)
		return
	}

	if match {
		file, err = convertToPDF(file, info)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			zero.Log.Error().Msgf("error converting to PDF: %v", err)
		}
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error reading file into bytes: %v", err)
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

	path := filepath.Join(conf.tmpDir, conf.filesDir, info.Name())
	dst, err := os.Create(path)
	if err != nil {
		return f, err
	}

	_, err = io.Copy(dst, file)
	if err != nil {
		return f, err
	}

	cmd := exec.Command(libreoffice, "--headless", "--convert-to", "pdf", "--outdir", filepath.Join(conf.tmpDir, conf.filesDir), path)

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
	zero.Log.Info().Msgf("GET handleResource id: %v cartridge %v", id, cartridge)

	inputFile := filepath.Join(conf.tmpDir, conf.uploadsDir, cartridge)
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

	err = os.MkdirAll(filepath.Join(conf.tmpDir, conf.uploadsDir), os.ModePerm)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	dst := filepath.Join(conf.tmpDir, conf.uploadsDir, file.Filename)
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
