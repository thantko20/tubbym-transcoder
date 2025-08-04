package storage

import (
	"io"
	"os"
	"path/filepath"
)

type Storage interface {
	Upload(src string, dst string) error
	Download(src string, dst string) error
	Delete(path string) error
}

type LocalStorage struct {
	baseDir string
}

func NewLocalStorage(baseDir string) *LocalStorage {
	return &LocalStorage{baseDir: baseDir}
}

func (s *LocalStorage) Upload(src, dst string) error {
	dstPath := filepath.Join(s.baseDir, dst)
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return err
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func (s *LocalStorage) Download(src, dst string) error {
	srcPath := filepath.Join(s.baseDir, src)
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	source, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func (s *LocalStorage) Delete(path string) error {
	return os.Remove(filepath.Join(s.baseDir, path))
}
