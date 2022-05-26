package models_test

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/commonsyllabi/viewer/internal/api/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const singleTestFile = "../../../../pkg/commoncartridge/test_files/test_01.imscc"

var attachmentID int64

func TestAttachmentModel(t *testing.T) {
	teardown := setup(t)
	defer teardown(t)

	t.Run("Test get all attachments", func(t *testing.T) {
		att, err := models.GetAllAttachments()
		require.Nil(t, err)
		assert.GreaterOrEqual(t, len(att), 1)
	})

	t.Run("Test add attachment", func(t *testing.T) {
		bytes, err := ioutil.ReadFile(singleTestFile)
		require.Nil(t, err)
		att := models.Attachment{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      "test_01.imscc",
			File:      bytes,
			Type:      "zip",
		}

		a, err := models.AddNewAttachment(&att)
		attachmentID = a.ID
		assert.Nil(t, err)
	})

	t.Run("Test get attachment", func(t *testing.T) {
		att, err := models.GetAttachment(int(attachmentID))
		require.Nil(t, err)

		assert.Equal(t, att.ID, attachmentID)
	})

	t.Run("Test update attachment", func(t *testing.T) {
		bytes, err := ioutil.ReadFile(singleTestFile)
		if err != nil {
			t.Error(err)
		}
		att := models.Attachment{
			UpdatedAt: time.Now(),
			Name:      "test_01.imscc (updated)",
			File:      bytes,
			Type:      "zip",
		}

		updated, err := models.UpdateAttachment(1, &att)
		require.Nil(t, err)
		assert.Equal(t, updated.Name, att.Name)
	})

	t.Run("Test delete attachment", func(t *testing.T) {
		err := models.DeleteAttachment(1)
		assert.Nil(t, err)
	})
}
