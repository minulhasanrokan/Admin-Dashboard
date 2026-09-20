package local

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type Local struct {
	BasePath       string
	BaseURL        string
	BackupBasePath string
}

func New(path string, url string, backupPath string) *Local {
	return &Local{
		BasePath:       path,
		BaseURL:        url,
		BackupBasePath: backupPath,
	}
}

func (l *Local) GetBasePath() string {
	return l.BasePath
}

func (l *Local) BackupGetBasePath() string {
	return l.BackupBasePath
}

func (l *Local) Upload(fileName string, file io.Reader) (string, error) {

	if strings.TrimSpace(fileName) == "" {
		return "", errors.New("file name is empty")
	}

	if file == nil {
		return "", errors.New("nil upload file")
	}

	cleanFileName := filepath.Clean(fileName)

	if filepath.IsAbs(cleanFileName) ||
		cleanFileName == ".." ||
		strings.HasPrefix(
			cleanFileName,
			".."+string(os.PathSeparator),
		) {
		return "", errors.New("invalid file name")
	}

	if err := os.MkdirAll(l.BasePath, os.ModePerm); err != nil {
		return "", err
	}

	fullPath := filepath.Join(
		l.BasePath,
		cleanFileName,
	)

	if err := os.MkdirAll(
		filepath.Dir(fullPath),
		os.ModePerm,
	); err != nil {
		return "", err
	}

	out, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = out.Close()
	}()

	if _, err := io.Copy(out, file); err != nil {
		return "", err
	}

	return cleanFileName, nil
}

func (l *Local) Delete(fileName string) error {
	return os.Remove(
		filepath.Join(l.BasePath, fileName),
	)
}

func (l *Local) GetFile(fileName string) (io.ReadCloser, error) {

	path := filepath.Join(
		l.BasePath,
		fileName,
	)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (l *Local) UploadFile(fileName string, header *multipart.FileHeader) (string, error) {

	if fileName == "" {
		return "", errors.New("file name is empty")
	}

	if header == nil {
		return "", errors.New("file header is empty")
	}

	file, err := header.Open()
	if err != nil {
		return "", err
	}

	defer func() {
		_ = file.Close()
	}()

	return l.Upload(fileName, file)
}

func (l *Local) CreateDirectory(path string) error {

	if strings.TrimSpace(path) == "" {
		return errors.New("directory path is empty")
	}

	fullPath := filepath.Join(
		l.BasePath,
		path,
	)

	return os.MkdirAll(
		fullPath,
		os.ModePerm,
	)
}

func (l *Local) BackupUpload(fileName string, file io.Reader) (string, error) {

	if strings.TrimSpace(fileName) == "" {
		return "", errors.New("backup file name is empty")
	}

	if file == nil {
		return "", errors.New("nil backup upload file")
	}

	if strings.TrimSpace(l.BackupBasePath) == "" {
		return "", errors.New("backup storage path is empty")
	}

	cleanFileName := filepath.Clean(fileName)

	if filepath.IsAbs(cleanFileName) ||
		cleanFileName == ".." ||
		strings.HasPrefix(
			cleanFileName,
			".."+string(os.PathSeparator),
		) {
		return "", errors.New("invalid backup file name")
	}

	fullPath := filepath.Join(
		l.BackupBasePath,
		cleanFileName,
	)

	if err := os.MkdirAll(
		filepath.Dir(fullPath),
		os.ModePerm,
	); err != nil {
		return "", err
	}

	out, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = out.Close()
	}()

	if _, err := io.Copy(out, file); err != nil {
		return "", err
	}

	return cleanFileName, nil
}

func (l *Local) BackupDelete(fileName string) error {

	if strings.TrimSpace(fileName) == "" {
		return errors.New("backup file name is empty")
	}

	return os.Remove(
		filepath.Join(
			l.BackupBasePath,
			fileName,
		),
	)
}

func (l *Local) BackupGetFile(fileName string) (io.ReadCloser, error) {

	if strings.TrimSpace(fileName) == "" {
		return nil, errors.New("backup file name is empty")
	}

	path := filepath.Join(
		l.BackupBasePath,
		fileName,
	)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (l *Local) BackupUploadFile(fileName string, header *multipart.FileHeader) (string, error) {

	if fileName == "" {
		return "", errors.New("backup file name is empty")
	}

	if header == nil {
		return "", errors.New("backup file header is empty")
	}

	file, err := header.Open()
	if err != nil {
		return "", err
	}

	defer func() {
		_ = file.Close()
	}()

	return l.BackupUpload(
		fileName,
		file,
	)
}

func (l *Local) BackupCreateDirectory(path string) error {

	if strings.TrimSpace(path) == "" {
		return errors.New("backup directory path is empty")
	}

	if strings.TrimSpace(l.BackupBasePath) == "" {
		return errors.New("backup storage path is empty")
	}

	fullPath := filepath.Join(
		l.BackupBasePath,
		path,
	)

	return os.MkdirAll(
		fullPath,
		os.ModePerm,
	)
}
