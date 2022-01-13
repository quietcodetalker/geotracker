package pagination

import (
	"encoding/base64"
	"errors"
	"fmt"
)

var (
	ErrInvalidCursor = errors.New("invalid cursor")
)

// EncodeCursor encodes page token and page size into opaque page cursor.
func EncodeCursor(pageToken int, pageSize int) string {
	rawStr := fmt.Sprintf("%d %d", pageToken, pageSize)
	return base64.URLEncoding.EncodeToString([]byte(rawStr))
}

// DecodeCursor decodes opaque page cursor into page token and page size.
func DecodeCursor(cursor string) (int, int, error) {
	rawStr, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return 0, 0, ErrInvalidCursor
	}

	var pageToken, pageSize int
	_, err = fmt.Sscanf(string(rawStr), "%d %d", &pageToken, &pageSize)
	if err != nil {
		return 0, 0, ErrInvalidCursor
	}

	return pageToken, pageSize, nil
}
