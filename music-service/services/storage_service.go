package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"music-service/config"
)

type StorageService interface {
	SaveAudio(ctx context.Context, fileHeader *multipart.FileHeader) (relativePath string, err error)
	SaveCover(ctx context.Context, fileHeader *multipart.FileHeader) (relativePath string, err error)
}

type storageService struct {
	cfg *config.Config
}

func NewStorageService(cfg *config.Config) StorageService {
	return &storageService{cfg: cfg}
}

func (s *storageService) SaveAudio(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	return s.save(ctx, fileHeader, s.cfg.StorageAudioDir, []string{".mp3", ".wav", ".flac", ".m4a", ".ogg"})
}

func (s *storageService) SaveCover(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	return s.save(ctx, fileHeader, s.cfg.StorageCoversDir, []string{".jpg", ".jpeg", ".png", ".webp"})
}

func (s *storageService) save(ctx context.Context, fh *multipart.FileHeader, dir string, allowedExt []string) (string, error) {
	if s.cfg == nil {
		return "", errors.New("storage config is nil")
	}
	if fh == nil {
		return "", errors.New("file is required")
	}
	if dir == "" {
		return "", errors.New("storage dir is empty")
	}

	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if !isAllowedExt(ext, allowedExt) {
		return "", errors.New("file type is not allowed")
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	name, err := randomName(16)
	if err != nil {
		return "", err
	}

	dstRelative := filepath.Join(dir, name+ext)
	dstAbsolute := dstRelative

	src, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(dstAbsolute)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	buf := make([]byte, 32*1024)
	for {
		select {
		case <-ctx.Done():
			_ = os.Remove(dstAbsolute)
			return "", ctx.Err()
		default:
		}

		n, rerr := src.Read(buf)
		if n > 0 {
			if _, werr := dst.Write(buf[:n]); werr != nil {
				_ = os.Remove(dstAbsolute)
				return "", werr
			}
		}
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			_ = os.Remove(dstAbsolute)
			return "", rerr
		}
	}

	// Возвращаем относительный путь внутри проекта
	return dstRelative, nil
}

func isAllowedExt(ext string, allowed []string) bool {
	for _, a := range allowed {
		if ext == a {
			return true
		}
	}
	return false
}

func randomName(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
