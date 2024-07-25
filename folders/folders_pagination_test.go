package folders_test

import (
	"testing"

	"github.com/georgechieng-sc/interns-2022/folders"
	"github.com/stretchr/testify/assert"
	"github.com/gofrs/uuid"
)

func Test_GetAllFoldersPaginated(t *testing.T) {
	t.Run("returns error when limit is non-positive", func(t *testing.T) {
		orgID := uuid.FromStringOrNil(folders.DefaultOrgID)
		res, err := folders.GetAllFoldersPaginated(&folders.PaginatedFetchRequest{
			OrgID:  orgID,
			Limit:  -1,
			Cursor: "",
		})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("returns error for nil request", func(t *testing.T) {
		res, err := folders.GetAllFoldersPaginated(nil)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("returns error for invalid cursor token", func(t *testing.T) {
		orgID := uuid.FromStringOrNil(folders.DefaultOrgID)
		res, err := folders.GetAllFoldersPaginated(&folders.PaginatedFetchRequest{
			OrgID:  orgID,
			Limit:  5,
			Cursor: "invalidToken",
		})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("fetches the first 5 folders", func(t *testing.T) {
		orgID := uuid.FromStringOrNil(folders.DefaultOrgID)
		res, err := folders.GetAllFoldersPaginated(&folders.PaginatedFetchRequest{
			OrgID:  orgID,
			Limit:  5,
			Cursor: "",
		})

		assert.NoError(t, err)

		expected, _ := folders.FetchAllFoldersByOrgID(orgID)
		assert.Equal(t, expected[0:5], res.Folders)
	})

	t.Run("fetches the first 5 folders and then the next 5", func(t *testing.T) {
		orgID := uuid.FromStringOrNil(folders.DefaultOrgID)
		firstBatch, _ := folders.GetAllFoldersPaginated(&folders.PaginatedFetchRequest{
			OrgID:  orgID,
			Limit:  5,
			Cursor: "",
		})

		secondBatch, err := folders.GetAllFoldersPaginated(&folders.PaginatedFetchRequest{
			OrgID:  orgID,
			Limit:  5,
			Cursor: firstBatch.NextCursor,
		})

		expected, _ := folders.FetchAllFoldersByOrgID(orgID)
		assert.NoError(t, err)
		assert.Equal(t, expected[5:10], secondBatch.Folders)
	})

	t.Run("fetches near the end of the folder list", func(t *testing.T) {
		orgID := uuid.FromStringOrNil(folders.DefaultOrgID)

		expected, _ := folders.FetchAllFoldersByOrgID(orgID)
		nextCursor := folders.EncodeCursor(len(expected) - 3)

		res, err := folders.GetAllFoldersPaginated(&folders.PaginatedFetchRequest{
			OrgID:  orgID,
			Limit:  5,
			Cursor: nextCursor,
		})

		assert.NoError(t, err)
		assert.Equal(t, len(expected[len(expected)-3:]), len(res.Folders))
		assert.Equal(t, expected[len(expected)-3:], res.Folders)
	})
}

func Test_EncodeCursor(t *testing.T) {
	t.Run("encode and decode", func(t *testing.T) {
		original := 5
		encoded := folders.EncodeCursor(original)
		decoded, err := folders.DecodeCursor(encoded)

		assert.NoError(t, err)
		assert.Equal(t, original, decoded)
	})

	t.Run("decode invalid cursor", func(t *testing.T) {
		_, err := folders.DecodeCursor("invalid_cursor")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "illegal base64 data")
	})

	t.Run("decode empty cursor", func(t *testing.T) {
		decoded, err := folders.DecodeCursor("")

		assert.NoError(t, err)
		assert.Equal(t, 0, decoded)
	})

	t.Run("decode invalid base64", func(t *testing.T) {
		_, err := folders.DecodeCursor("ThisIsNotBase64!")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "illegal base64 data")
	})
}