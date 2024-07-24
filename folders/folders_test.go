package folders_test

import (
	"testing"
	// "github.com/georgechieng-sc/interns-2022/folders"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

// TestGetAllFolders tests the GetAllFolders function for various scenarios.
func TestGetAllFolders(t *testing.T) {
	
	t.Run("Invalid OrgID", func(t *testing.T) {
		// Invalid UUID
		testReq := &FetchFolderRequest{
			OrgID: uuid.Nil,
		}

		result, err := GetAllFolders(testReq)
		assert.Error(t, err, "Expected error for invalid OrgID")
		assert.Nil(t, result)
	})

	t.Run("Valid OrgID with Matching Folders", func(t *testing.T) {
		testOrgID := uuid.Must(uuid.NewV4())
		testReq := &FetchFolderRequest{
			OrgID: testOrgID,
		}

		// Mock the GetSampleData function to return test data
		GetSampleData = func() []*Folder {
			return []*Folder{
				{OrgId: testOrgID, Name: "Test Folder"},
			}
		}

		result, err := GetAllFolders(testReq)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Folders, 1, "Expected one folder in result")
	})

	t.Run("Valid OrgID with No Matching Folders", func(t *testing.T) {
		testOrgID := uuid.Must(uuid.NewV4())
		
		// Different OrgID
		otherOrgID := uuid.Must(uuid.NewV4())
		testReq := &FetchFolderRequest{
			OrgID: testOrgID,
		}

		// Mock the GetSampleData function to return test data with different OrgID
		GetSampleData = func() []*Folder {
			return []*Folder{
				{OrgId: otherOrgID, Name: "Unrelated Folder"},
			}
		}

		result, err := GetAllFolders(testReq)
		assert.Error(t, err, "Expected error for no matching folders")
		assert.Nil(t, result)
	})
}
