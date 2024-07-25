package folders

import (
	"errors"
	"strings"

	"github.com/gofrs/uuid"
	"encoding/base64"
	"strconv"
)

/**
 * Idea behind the solution:
 * 
 * I chose to use a cursor-based method for pagination because it's more efficient than using offsets.
 * Instead of searching through all the data to find the starting point, we can jump straight to the 
 * next set of data. A user sends a request with three main parts: the organization ID (orgID), a limit 
 * on how many folders to fetch (this should be a positive number), and a cursor token (encoded in base64) 
 * which tells us where to start fetching the next set of folders. This token can be empty if it's the 
 * first request. When you call the PaginatedFetchRequest function, it gives you a list of folders 
 * starting from the given index up to the limit. It also gives you a new cursor that points to where the 
 * next set of folders should start. If the cursor token is empty, we start from the beginning. This 
 * function keeps getting called until there are no more folders to fetch, at which point it returns an 
 * empty cursor.
 * 
 */

type PaginatedFetchRequest struct {
	Limit  int
	Cursor string
	OrgID  uuid.UUID
}

type PaginatedFetchResponse struct {
	NextCursor string
	Folders    []*Folder
}

/**
 * GetAllFoldersPaginated fetches a paginated list of folders based on the request parameters.
 * 
 * @param req - The request containing the organization ID, limit, and cursor.
 * @return PaginatedFetchResponse - The response containing the list of folders and the next cursor.
 * @return error - Any error encountered during the process.
 * 
 * The function performs the following steps:
 * 1. validates the request parameters
 * 2. decodes the cursor to determine the starting index
 * 3. fetches all folders for the specified organization ID
 * 4. calculates the end index for the current page of results
 * 5. generates the next cursor for pagination
 * 6. returns the paginated response containing the folders and the next cursor
 */
func GetAllFoldersPaginated(req *PaginatedFetchRequest) (*PaginatedFetchResponse, error) {
	if err := validateRequest(req); err != nil {
		return nil, err
	}

	start, err := getStartIndex(req.Cursor)
	if err != nil {
		return nil, err
	}

	folders, err := FetchAllFoldersByOrgID(req.OrgID)
	if err != nil {
		return nil, err
	}

	end := calculateEndIndex(start, req.Limit, len(folders))

	nextCursor := generateNextCursor(end, len(folders))

	return &PaginatedFetchResponse{
		Folders:    folders[start:end],
		NextCursor: nextCursor,
	}, nil
}

/**
 * validateRequest checks if the request parameters are valid.
 */
func validateRequest(req *PaginatedFetchRequest) error {
	if req == nil {
		return errors.New("request invalid, cannot be nil")
	}
	if req.Limit <= 0 {
		return errors.New("limit has to be greater than 0")
	}
	return nil
}

/**
 * getStartIndex decodes the cursor to determine the starting index.
 */
func getStartIndex(cursor string) (int, error) {
	if cursor == "" {
		return 0, nil
	}
	return DecodeCursor(cursor)
}

/**
 * calculateEndIndex calculates the end index for the current page of results.
 */
func calculateEndIndex(start, limit, total int) int {
	end := start + limit
	if end > total {
		end = total
	}
	return end
}

/**
 * generateNextCursor generates the next cursor for pagination.
 */
func generateNextCursor(end, total int) string {
	if end == total {
		return ""
	}
	return EncodeCursor(end)
}

/**
 * EncodeCursor encodes an index into a base64 encoded cursor string.
 */
func EncodeCursor(index int) string {
	return base64.StdEncoding.EncodeToString([]byte("next cursor:" + strconv.Itoa(index)))
}

/**
 * DecodeCursor decodes a base64 encoded cursor string into an index.
 */
func DecodeCursor(encodedCursor string) (int, error) {
	if encodedCursor == "" {
		return 0, nil
	}

	decoded, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return 0, err
	}
	parts := strings.Split(string(decoded), ":")
	if len(parts) < 2 {
		return 0, errors.New("invalid cursor format")
	}
	return strconv.Atoi(parts[1])
}
