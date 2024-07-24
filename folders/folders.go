package folders

import (
	"github.com/gofrs/uuid"
	"errors"
)
/**
 * Refactoring suggestions:
 * - Delete instances of unused variable declarations to make it less confusing within the code
 * - Include error checking to handle potential issues 
 * - Use a more descriptive name like folders instead of just 'f'
 * - Directly manipulate pointers rather than copying values into a new slice only to take their addresses later
 * - Eliminate unnecessary loops and operations that copy data without modifying or filtering it, thereby directly using the returned data structure.
 */

/**
 * Function retrieves all folders for a given organization ID and returns a structured response.
 *
 * @param req - The fetch folder request containing the organization ID.
 * @returns A response containing all folders associated with the organization ID or an error if the request is invalid or folders cannot be fetched.
 */
func GetAllFolders(req *FetchFolderRequest) (*FetchFolderResponse, error) {
    if req == nil {
        return nil, errors.New("request cannot be nil")
    }

    folders, err := FetchAllFoldersByOrgID(req.OrgID)
    if err != nil {
        return nil, err
    }

    return &FetchFolderResponse{Folders: folders}, nil
}

/**
 * Refactoring suggestions:
 * - Consider returning an error if no folders are found
 */

/**
 * Function retrieves a list of folders associated with a specific organization ID.
 *
 * @param orgID - The unique identifier of the organization to fetch folders for.
 * @returns A slice of folder pointers that match the given organization ID, or an error if no folders match.
 */
func FetchAllFoldersByOrgID(orgID uuid.UUID) ([]*Folder, error) {
	allFolders := GetSampleData()

	var matchedFolders []*Folder
	for _, folder := range allFolders {
		if folder.OrgId == orgID {
			matchedFolders = append(matchedFolders, folder)
		}
	}

	if len(matchedFolders) == 0 {
		return nil, errors.New("no folders found for the specified orgID")
	}

	return matchedFolders, nil
}
