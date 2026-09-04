package repo

import (
	"io"
	"os"
	"path/filepath"

	"github.com/Cheasezz/fileService/internal/core"
)

type FileSystem struct {
	path string
}

func New() (*FileSystem, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	storagePath := filepath.Join(homeDir, ".fileService/uploads")

	err = createDir(storagePath)
	if err != nil {
		return nil, err
	}

	return &FileSystem{path: storagePath}, nil
}

func (fs *FileSystem) CreateFile(fileInfo *core.FileInfo) (*os.File, error) {
	var f *os.File

	userDir := fs.path + "/" + fileInfo.UserID.String()

	// If dir already exist return nil
	err := createDir(userDir)
	if err != nil {
		return nil, err
	}

	f, err = os.Create(userDir + "/" + fileInfo.Name)
	if err != nil {
		return f, err
	}

	return f, nil
}

func (fs *FileSystem) OpenFile(fileInfo *core.FileInfo) (io.ReadCloser, error) {
	filePath := filepath.Join(fs.path, fileInfo.UserID.String(), fileInfo.Name)

	f, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, core.ErrFileNotFound
		}
		return nil, err
	}

	return f, nil
}

func createDir(path string) error {
	err := os.MkdirAll(path, 0o755)
	if err != nil {
		return err
	}

	return nil
}
