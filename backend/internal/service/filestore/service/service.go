package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// UrlPrefix is the path under which the app serves dir
const UrlPrefix = "/files/"

// Service stores photos on local disk (a volume in docker).
// ponytail: local disk, single instance; S3/MinIO with signed URLs once there is a second instance
type Service struct {
	dir string
}

func New(dir string) (*Service, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("os.MkdirAll: %w", err)
	}
	return &Service{dir: dir}, nil
}

// Save writes the file under a random name (128 bits) and returns its URL
func (s *Service) Save(_ context.Context, ext string, data []byte) (string, error) {
	name := strings.ToLower(rand.Text()) + "." + ext

	if err := os.WriteFile(filepath.Join(s.dir, name), data, 0o644); err != nil {
		return "", fmt.Errorf("os.WriteFile: %w", err)
	}
	return UrlPrefix + name, nil
}
