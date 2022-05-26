package models_test

import (
	"testing"
	"time"

	"github.com/commonsyllabi/viewer/internal/api/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//-- to not recreate the whole table everytime, one can also do all transactions, and then rollback rather than commit at the end of each test
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
		_, err := models.AddNewSyllabus(&syll)
		assert.Nil(t, err)
	})

	t.Run("Test get syllabus", func(t *testing.T) {
		syll, err := models.GetSyllabus(1)
		require.Nil(t, err)
		assert.Equal(t, syll.ID, int64(1))
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
