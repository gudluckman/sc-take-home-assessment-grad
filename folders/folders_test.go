package folders

import (
	"testing"
	// "github.com/georgechieng-sc/interns-2022/folders"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func setupOrgID(orgID string) uuid.UUID {
	return uuid.FromStringOrNil(orgID)
}

func TestGetAllFolders(t *testing.T) {
	validOrgID := setupOrgID("6591e16c-c257-4366-bf6d-650c68f71284")
	nonexistentOrgID := uuid.Must(uuid.NewV4())

	tests := []struct {
		name      string
		orgID     uuid.UUID
		expectErr bool
		errMsg    string
		nonEmpty  bool
	}{
		{"Valid OrgID - Existing Folders", validOrgID, false, "", true},
		{"Nonexistent OrgID", nonexistentOrgID, true, "no folders found", false},
		{"Nil UUID", uuid.Nil, true, "Invalid ORG ID, cannot be nil", false},
		{"Performance Check", validOrgID, false, "", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &FetchFolderRequest{OrgID: tc.orgID}
			res, err := GetAllFolders(req)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errMsg)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				if tc.nonEmpty {
					assert.NotEmpty(t, res.Folders)
					for _, folder := range res.Folders {
						assert.NotEqual(t, uuid.Nil, folder.Id, "Folder ID must not be the nil UUID; it should be a properly set unique identifier.")
						assert.Equal(t, tc.orgID, folder.OrgId, "The organization ID of the folder must match the organization ID provided in the test case.")
						assert.NotEmpty(t, folder.Name, "The name of the folder must not be empty; it should contain valid text.")

					}
				}
			}
		})
	}
}

func TestFetchAllFoldersByOrgID(t *testing.T) {
	validOrgID := setupOrgID("6591e16c-c257-4366-bf6d-650c68f71284")
	nonexistentOrgID := uuid.Must(uuid.NewV4())

	tests := []struct {
		name      string
		orgID     uuid.UUID
		expectErr bool
		errMsg    string
		nonEmpty  bool
	}{
		{"Valid OrgID - Existing Folders", validOrgID, false, "", true},
		{"Nonexistent OrgID", nonexistentOrgID, true, "no folders found", false},
		{"Nil UUID", uuid.Nil, true, "Invalid ORG ID", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			folders, err := FetchAllFoldersByOrgID(tc.orgID)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errMsg)
				assert.Nil(t, folders)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, folders)
				for _, folder := range folders {
					assert.NotEqual(t, uuid.Nil, folder.Id, "Folder ID must not be the nil UUID; it should be a properly set unique identifier.")
					assert.Equal(t, tc.orgID, folder.OrgId, "The organization ID of the folder must match the organization ID provided in the test case.")
					assert.NotEmpty(t, folder.Name, "The name of the folder must not be empty; it should contain valid text.")
				}
			}
		})
	}
}