package storage

import (
	"io"
	"mime/multipart"
)

type Storage interface {
	Upload(fileName string, file io.Reader) (string, error)
	Delete(fileName string) error
	GetFile(fileName string) (io.ReadCloser, error)
	UploadFile(fileName string, header *multipart.FileHeader) (string, error)
	CreateDirectory(path string) error
	GetBasePath() string
	BackupUpload(fileName string, file io.Reader) (string, error)
	BackupDelete(fileName string) error
	BackupGetFile(fileName string) (io.ReadCloser, error)
	BackupUploadFile(fileName string, header *multipart.FileHeader) (string, error)
	BackupCreateDirectory(path string) error
	BackupGetBasePath() string
}
