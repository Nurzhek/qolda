package filestore

import "context"

// Store keeps uploaded photos
type Store interface {
	Save(ctx context.Context, ext string, data []byte) (url string, err error)
}
