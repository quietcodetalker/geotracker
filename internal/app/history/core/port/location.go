//go:generate mockgen -destination=mock/mock_location.go -package=mock . LocationClient

package port

import "context"

// LocationClient TODO: add description
type LocationClient interface {
	GetUserIDByUsername(ctx context.Context, username string) (int, error)
}
