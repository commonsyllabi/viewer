package handlers

import (
	"net/http"

	"github.com/commonsyllabi/viewer/internal/api/models"
	"github.com/commonsyllabi/viewer/internal/db"
	zero "github.com/commonsyllabi/viewer/internal/logger"
	"github.com/gin-gonic/gin"
)

func NewSyllabus(c *gin.Context) {
	syll := models.Syllabus{Id: 0, Title: "Others", Description: "They will always be there."}

	err := db.AddNewSyllabus(syll)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Log.Error().Msgf("error loading CC from disk: %v", err)
		return
	}

}
