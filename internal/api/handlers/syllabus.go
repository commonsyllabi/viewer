package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/commonsyllabi/viewer/internal/api/models"
	zero "github.com/commonsyllabi/viewer/internal/logger"
	"github.com/gin-gonic/gin"
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

	// save the actual syllabus
	var syll models.Syllabus
	err := c.Bind(&syll)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	syll, err = models.AddNewSyllabus(&syll)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error creating syllabus: %v", err)
		return
	}

	// get the attachments
	var attachments []models.Attachment
	form, _ := c.MultipartForm()
	files := form.File["attachments[]"]
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

	result, err := models.GetSyllabus(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error getting syllabus %d: %v", id, err)
		return
	}

	bytes, err := json.Marshal(result)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error marshalling syllabus: %v", err)
		return
	}

	c.JSON(http.StatusOK, string(bytes))
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
