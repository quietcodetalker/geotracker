package port

import "context"

// LocationClient TODO: add description
type LocationClient interface {
	GetUserIDByUsername(ctx context.Context, username string) (int, error)
}
