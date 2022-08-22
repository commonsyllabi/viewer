package handlers

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/commonsyllabi/viewer/internal/api/models"
	zero "github.com/commonsyllabi/viewer/internal/logger"
	"github.com/gin-gonic/gin"
)

func AllSyllabi(c *gin.Context) {
	syllabi, err := models.GetAllSyllabi()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error getting syllabi: %v", err)
		return
	}

	bytes, err := json.Marshal(syllabi)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error marshalling syllabus: %v", err)
		return
	}

	c.JSON(http.StatusOK, string(bytes))
}

func NewSyllabus(c *gin.Context) {

	err := sanitizeInput(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		zero.Error(err.Error())
		return
	}

	// save the actual syllabus
	var syll models.Syllabus
	err = c.Bind(&syll)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//-- hash the email
	hasher := sha256.New()
	hasher.Write([]byte(syll.Email))
	syll.Email = base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "error parsing form %v", err)
		zero.Errorf("error parsing form: %v", err)
		return
	}

	syll.CreatedAt = time.Now()
	syll.UpdatedAt = time.Now()

	syll, err = models.AddNewSyllabus(&syll)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error creating syllabus: %v", err)
		return
	}

	var attachments []models.Attachment
	files := form.File["attachments[]"]

	zero.Warnf("%d attachments found on new syllabus", len(files))

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
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
			Name:               f.Filename,
			SyllabusAttachedID: syll.ID,
			File:               bytes,
			Type:               http.DetectContentType(bytes),
		}

		att, err := models.AddNewAttachment(&attachment)
		if err != nil {
			zero.Warnf("error adding attachment: %s", err)
		}
		attachments = append(attachments, att)
	}

	s, err := json.Marshal(syll)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error marshalling syllabus: %v", err)
		return
	}

	a, err := json.Marshal(attachments)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error marshalling attachments: %v", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"syllabus":    string(s),
		"attachments": string(a),
	})
}

func UpdateSyllabus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		zero.Errorf("not a valid id %d", id)
		return
	}

	var syll models.Syllabus
	err = c.Bind(&syll)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	syll.UpdatedAt = time.Now()

	_, err = models.UpdateSyllabus(id, &syll)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error updating syllabus %d: %v", id, err)
		return
	}

	bytes, err := json.Marshal(syll)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error marshalling syllabus: %v", err)
		return
	}

	c.JSON(http.StatusOK, string(bytes))
}

func GetSyllabus(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		zero.Errorf("not a valid id %d", id)
		return
	}

	if !strings.Contains(c.GetHeader("Content-Type"), "application/json") && gin.Mode() != gin.TestMode {

		c.HTML(http.StatusOK, "Syllabus", id)
		return
	}

	syll, err := models.GetSyllabus(id)
	if err != nil {
		zero.Errorf("error getting syllabus %v: %s", id, err)
		c.HTML(http.StatusOK, "Error", gin.H{
			"msg": "We couldn't find the syllabus.",
		})

		return
	}

	bytes, err := json.Marshal(syll)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error marshalling syllabus: %v", err)
		return
	}

	// nice separation of response body https://stackoverflow.com/a/56722847/4665412
	if c.GetHeader("Content-Type") == "application/json" {
		c.JSON(http.StatusOK, string(bytes))
	} else if gin.Mode() == gin.TestMode {
		c.JSON(http.StatusOK, string(bytes))
		return
	}

}

func DisplayMagicLink(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "Error", gin.H{
			"msg": "The ID of the resource you're asking for is invalid.",
		})
		zero.Errorf("error parsing the ID param: %v", err)
		return
	}

	token := c.Query("token")
	if id == 0 || token == "" {
		c.HTML(http.StatusBadRequest, "Error", gin.H{
			"msg": "The ID of the resource you're asking for is invalid.",
		})
		zero.Errorf("invalid values for ID (%v) or token (%v)", id, token)
		return
	}

	magic_token, err := models.GetTokenSyllabus(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "Error", gin.H{
			"msg": "The link you're trying to access has expired.",
		})
		zero.Warnf("error getting magic token: %v", err.Error())
		return
	}

	if token == base64.URLEncoding.EncodeToString(magic_token.Token) {
		if gin.Mode() == gin.TestMode {
			c.JSON(http.StatusOK, token)
			return
		}

		syll, err := models.GetSyllabus(id)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "Error", gin.H{
				"msg": "We could not find the syllabus you are looking for.",
			})
			zero.Warnf("error getting syllabus: %v", err)
			return
		}

		c.HTML(http.StatusOK, "MagicLink", syll)
	} else {
		c.HTML(http.StatusForbidden, "Error", gin.H{
			"msg": "You're trying to access a protected resource.",
		})
	}
}

func DeleteSyllabus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		zero.Errorf("not a valid id %d", id)
		return
	}

	err = models.DeleteSyllabus(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error getting syllabus %d: %v", id, err)
		return
	}

	//-- TODO delete any associated token

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func sanitizeInput(c *gin.Context) error {
	var err error

	if c.PostForm("title") == "" || c.PostForm("description") == "" || c.PostForm("email") == "" {
		zero.Error("Cannot have empty title, description or email")
		err = fmt.Errorf("cannot have empty title, description or email")
		return err
	}

	if len(c.PostForm("title")) < 10 && len(c.PostForm("title")) > 200 {
		zero.Errorf("the title of the syllabus should be between 10 and 200 characters: %d", len(c.PostForm("title")))
		err = fmt.Errorf("the title of the syllabus should be between 10 and 200 characters: %d", len(c.PostForm("title")))
		return err
	}

	if len(c.PostForm("description")) < 10 && len(c.PostForm("description")) > 500 {
		zero.Errorf("the description of the syllabus should be between 10 and 500 characters: %d", len(c.PostForm("description")))
		err = fmt.Errorf("the description of the syllabus should be between 10 and 500 characters: %d", len(c.PostForm("description")))
		return err
	}

	if len(c.PostForm("email")) < 10 && len(c.PostForm("email")) > 50 {
		zero.Errorf("the email of the syllabus should be between 10 and 50 characters: %d", len(c.PostForm("email")))
		err = fmt.Errorf("the email of the syllabus should be between 10 and 50 characters: %d", len(c.PostForm("email")))
		return err
	}

	_, err = mail.ParseAddress(c.PostForm("email"))
	if err != nil {
		return err
	}

	return nil
}
