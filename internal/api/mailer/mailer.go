package mailer

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/commonsyllabi/viewer/internal/api/models"
	zero "github.com/commonsyllabi/viewer/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/mailgun/mailgun-go/v4"
)

const DOMAIN = "post.enframed.net"

func HandleMagicLink(c *gin.Context) {
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "The syllabus you're trying to access is invalid.")
		zero.Errorf("error getting id for magic link: %v", err)
		return
	}
	email := c.PostForm("email")

	syll, err := models.GetSyllabus(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error getting syllabus by email: %v", err)
		return
	}

	hasher := sha256.New()
	hasher.Write([]byte(email))
	t := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	if t != syll.Email {
		c.String(http.StatusBadRequest, "The email you've entered doesn't correspond to the one submitted with the syllabus.")
		zero.Warn("the submitted email for magic link doesn't correspond to the associated email.")
		return
	}

	token := models.MagicToken{
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Token:           hasher.Sum(nil),
		SyllabusTokenID: syll.ID,
	}
	token, err = models.AddNewToken(&token)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error setting magic token: %v", err)
		return
	}

	mg := mailgun.NewMailgun(DOMAIN, os.Getenv("MAILGUN_PRIVATE_API_KEY"))
	mg.SetAPIBase("https://api.eu.mailgun.net/v3") //-- rgpd mon amour

	sender := "Common Syllabi <cosyl@post.enframed.net>"
	subject := "Common Syllabi - Your Syllabus"
	body := fmt.Sprintf("Hi!\nHere is your magic link for syllabus: %d - %s\nlink: http://%s:%s/syllabi/edit/%d?token=%s", syll.ID, syll.Title, "localhost", "3046", syll.ID, base64.URLEncoding.EncodeToString(token.Token))
	recipient := email

	message := mg.NewMessage(sender, subject, body, recipient)

	if gin.Mode() == gin.TestMode {
		c.JSON(http.StatusOK, message)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, mg_id, err := mg.Send(ctx, message)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		zero.Errorf("error sending email: %v", err)
		return
	}

	zero.Debugf("ID: %s Resp: %s\n", mg_id, resp)
	c.String(http.StatusOK, "sending email to: %v", email)
}
