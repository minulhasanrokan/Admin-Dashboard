package s3

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const backupPrefix = "backup/"

type Storage struct {
	Client *s3.Client
	Bucket string
	Region string
}

func New(client *s3.Client, bucket string, region string) *Storage {
	return &Storage{
		Client: client,
		Bucket: bucket,
		Region: region,
	}
}

func (s *Storage) Upload(fileName string, file io.Reader) (string, error) {

	if strings.TrimSpace(fileName) == "" {
		return "", errors.New("file name is empty")
	}

	if file == nil {
		return "", errors.New("nil upload file")
	}

	_, err := s.Client.PutObject(
		context.Background(),
		&s3.PutObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(fileName),
			Body:   file,
		},
	)

	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (s *Storage) Delete(fileName string) error {

	_, err := s.Client.DeleteObject(
		context.Background(),
		&s3.DeleteObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(fileName),
		},
	)

	return err
}

func (s *Storage) GetFile(fileName string) (io.ReadCloser, error) {

	result, err := s.Client.GetObject(
		context.Background(),
		&s3.GetObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(fileName),
		},
	)

	if err != nil {
		return nil, err
	}

	return result.Body, nil
}

func (s *Storage) UploadFile(fileName string, header *multipart.FileHeader) (string, error) {

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

	return s.Upload(fileName, file)
}

func (s *Storage) CreateDirectory(path string) error {

	path = strings.TrimSpace(path)

	if path == "" {
		return errors.New("directory path is empty")
	}

	path = strings.TrimSuffix(path, "/") + "/"

	_, err := s.Client.PutObject(
		context.Background(),
		&s3.PutObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(path),
			Body:   strings.NewReader(""),
		},
	)

	return err
}

func (s *Storage) GetBasePath() string {
	return ""
}

func (s *Storage) BackupUpload(fileName string, file io.Reader) (string, error) {

	if strings.TrimSpace(fileName) == "" {
		return "", errors.New("backup file name is empty")
	}

	if file == nil {
		return "", errors.New("nil backup upload file")
	}

	key := backupKey(fileName)

	_, err := s.Client.PutObject(
		context.Background(),
		&s3.PutObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(key),
			Body:   file,
		},
	)

	if err != nil {
		return "", err
	}

	return key, nil
}

func (s *Storage) BackupDelete(fileName string) error {

	if strings.TrimSpace(fileName) == "" {
		return errors.New("backup file name is empty")
	}

	key := backupKey(fileName)

	_, err := s.Client.DeleteObject(
		context.Background(),
		&s3.DeleteObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(key),
		},
	)

	return err
}

func (s *Storage) BackupGetFile(fileName string) (io.ReadCloser, error) {

	if strings.TrimSpace(fileName) == "" {
		return nil, errors.New("backup file name is empty")
	}

	key := backupKey(fileName)

	result, err := s.Client.GetObject(
		context.Background(),
		&s3.GetObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(key),
		},
	)

	if err != nil {
		return nil, err
	}

	return result.Body, nil
}

func (s *Storage) BackupUploadFile(fileName string, header *multipart.FileHeader) (string, error) {

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

	return s.BackupUpload(
		fileName,
		file,
	)
}

func (s *Storage) BackupCreateDirectory(path string) error {

	path = strings.TrimSpace(path)

	if path == "" {
		return errors.New(
			"backup directory path is empty",
		)
	}

	key := backupKey(path)

	key = strings.TrimSuffix(key, "/") + "/"

	_, err := s.Client.PutObject(
		context.Background(),
		&s3.PutObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(key),
			Body:   strings.NewReader(""),
		},
	)

	return err
}

func (s *Storage) BackupGetBasePath() string {
	return backupPrefix
}

func backupKey(fileName string) string {

	fileName = strings.TrimSpace(fileName)

	fileName = strings.TrimPrefix(fileName, "/")

	fileName = strings.TrimPrefix(fileName, "backup/")

	return backupPrefix + fileName
}
