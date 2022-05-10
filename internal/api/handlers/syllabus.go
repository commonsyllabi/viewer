package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/commonsyllabi/viewer/internal/api/models"
	zero "github.com/commonsyllabi/viewer/internal/logger"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func AllSyllabi(c *gin.Context) {
	syllabi, err := models.GetAllSyllabi()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error getting syllabi: %v", err)
		return
	}

	bytes, err := json.Marshal(syllabi)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error marshalling syllabus: %v", err)
		return
	}

	c.JSON(http.StatusOK, string(bytes))
}

func NewSyllabus(c *gin.Context) {

	//-- sanitizing
	if c.PostForm("title") == "" || c.PostForm("description") == "" || c.PostForm("email") == "" {
		c.String(http.StatusBadRequest, "Cannot have empty title, description or email")
		zero.Log.Error().Msg("Cannot have empty title, description or email")
		return
	}

	// save the actual syllabus
	var syll models.Syllabus
	err := c.Bind(&syll)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//-- hash the email
	hashed, err := bcrypt.GenerateFromPassword([]byte(syll.Email), bcrypt.DefaultCost)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error hashing email: %v", err)
		return
	}

	syll.Email = string(hashed)

	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "error parsing form %v", err)
		zero.Log.Error().Msgf("error parsing form: %v", err)
		return
	}

	syll, err = models.AddNewSyllabus(&syll)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error creating syllabus: %v", err)
		return
	}

	var attachments []models.Attachment
	files := form.File["attachments[]"]

	zero.Log.Warn().Msgf("%d attachments found on new syllabus", len(files))

	for _, f := range files {
		file, err := f.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			zero.Log.Error().Msgf("error reading file into bytes: %v", err)
			return
		}

		attachment := models.Attachment{
			Name:       f.Filename,
			SyllabusID: syll.ID,
			File:       bytes,
			Type:       http.DetectContentType(bytes),
		}

		att, _ := models.AddNewAttachment(&attachment)
		attachments = append(attachments, att)
	}

	s, err := json.Marshal(syll)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error marshalling syllabus: %v", err)
		return
	}

	a, err := json.Marshal(attachments)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error marshalling attachments: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"syllabus":    string(s),
		"attachments": string(a),
	})
}

func UpdateSyllabus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		zero.Log.Error().Msgf("not a valid id %d", id)
		return
	}

	var syll models.Syllabus
	err = c.Bind(&syll)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = models.UpdateSyllabus(id, &syll)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error updating syllabus %d: %v", id, err)
		return
	}

	bytes, err := json.Marshal(syll)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error marshalling syllabus: %v", err)
		return
	}

	c.JSON(http.StatusOK, string(bytes))
}

func GetSyllabus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		zero.Log.Error().Msgf("not a valid id %d", id)
		return
	}

	syll, err := models.GetSyllabus(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error getting syllabus %d: %v", id, err)
		return
	}

	bytes, err := json.Marshal(syll)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error marshalling syllabus: %v", err)
		return
	}

	// nice separation of response body https://stackoverflow.com/a/56722847/4665412
	if c.GetHeader("Content-Type") == "application/json" {
		c.JSON(http.StatusOK, string(bytes))
	} else {
		c.HTML(http.StatusOK, "cartridge.tmpl", syll)
	}

}

func DeleteSyllabus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		zero.Log.Error().Msgf("not a valid id %d", id)
		return
	}

	err = models.DeleteSyllabus(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error getting syllabus %d: %v", id, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}
