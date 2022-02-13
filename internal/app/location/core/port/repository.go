package port

// Repository is a composite interface that embeds all repository interfaces.
type Repository interface {
	UserRepository
	LocationRepository
}
