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

func AllAttachments(c *gin.Context) {
	attachments, err := models.GetAllAttachments()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error getting attachments: %v", err)
		return
	}

	bytes, err := json.Marshal(attachments)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error marshalling attachments: %v", err)
		return
	}

	c.JSON(http.StatusOK, string(bytes))
}

func NewAttachment(c *gin.Context) {
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
			zero.Errorf("error reading file into bytes: %v", err)
			return
		}

		attachment := models.Attachment{
			Name: f.Filename,
			File: bytes,
			Type: http.DetectContentType(bytes),
		}

		att, _ := models.AddNewAttachment(&attachment)
		attachments = append(attachments, att)
	}

	a, err := json.Marshal(attachments)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error marshalling attachments: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"attachments": string(a)})
}

func UpdateAttachment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		zero.Errorf("not a valid id %d", id)
		return
	}

	var att models.Attachment
	err = c.Bind(&att)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = models.UpdateAttachment(id, &att)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error updating attachment %d: %v", id, err)
		return
	}

	bytes, err := json.Marshal(att)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error marshalling attachment: %v", err)
		return
	}

	c.JSON(http.StatusOK, string(bytes))
}

func GetAttachment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		zero.Errorf("not a valid id %d", id)
		return
	}

	result, err := models.GetAttachment(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error getting Attachment %d: %v", id, err)
		return
	}

	if c.Query("type") == "file" {
		c.Writer.WriteHeader(http.StatusOK)
		c.Header("Content-Disposition", "attachment; filename="+result.Name)
		mimeType := http.DetectContentType(result.File)
		c.Header("Content-Type", mimeType)
		c.Writer.Write(result.File)
	} else {
		bytes, err := json.Marshal(result)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			zero.Errorf("error marshalling Attachment: %v", err)
			return
		}

		c.JSON(http.StatusOK, string(bytes))
	}
}

func DeleteAttachment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		zero.Errorf("not a valid id %d", id)
		return
	}

	err = models.DeleteAttachment(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error getting Attachment %d: %v", id, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}
