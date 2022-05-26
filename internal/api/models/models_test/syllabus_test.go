package models_test

import (
	"testing"
	"time"

	"github.com/commonsyllabi/viewer/internal/api/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var syllabusID int64

func TestSyllabusModel(t *testing.T) {
	teardown := setup(t)
	defer teardown(t)

	t.Run("Test get all syllabi", func(t *testing.T) {
		syll, err := models.GetAllSyllabi()
		require.Nil(t, err)
		assert.NotEqual(t, len(syll), 0)
	})

	t.Run("Test add syllabus", func(t *testing.T) {
		syll := models.Syllabus{
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       "Test Title 2",
			Description: "Test Description 2",
		}
		s, err := models.AddNewSyllabus(&syll)
		syllabusID = s.ID
		assert.Nil(t, err)
	})

	t.Run("Test get syllabus", func(t *testing.T) {
		syll, err := models.GetSyllabus(int(syllabusID))
		require.Nil(t, err)
		assert.Equal(t, syll.ID, syllabusID)
	})

	t.Run("Test update syllabus", func(t *testing.T) {
		syll := models.Syllabus{
			UpdatedAt:   time.Now(),
			Title:       "Test Title 1 (updated)",
			Description: "Test Description 1 (updated",
		}
		updated, err := models.UpdateSyllabus(1, &syll)
		require.Nil(t, err)
		assert.Equal(t, updated.Description, syll.Description)
	})

	t.Run("Test delete syllabus", func(t *testing.T) {
		err := models.DeleteSyllabus(1)
		assert.Nil(t, err)
	})
}
