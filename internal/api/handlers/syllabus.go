package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/commonsyllabi/viewer/internal/api/models"
	"github.com/commonsyllabi/viewer/internal/db"
	zero "github.com/commonsyllabi/viewer/internal/logger"
	"github.com/gin-gonic/gin"
)

func NewSyllabus(c *gin.Context) {
	syll := &models.Syllabus{Id: 667, Title: "Others", Description: "They will always be there."}

	_, err := db.AddNewSyllabus(syll)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error creating syllabus: %v", err)
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
