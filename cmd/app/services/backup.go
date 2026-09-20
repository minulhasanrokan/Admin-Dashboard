package services

import (
	"admin-dashboard/cmd/app/repositories"
	"admin-dashboard/cmd/app/requests"
	"admin-dashboard/cmd/internal/backup"
	"admin-dashboard/cmd/internal/constants"
	"admin-dashboard/cmd/internal/models"
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (s *Service) GenerateDatabaseBackup() (string, error) {

	driver := strings.ToLower(strings.TrimSpace(s.Db.GetDriver()))

	extension := ".sql"

	switch driver {
	case "mysql", "mariadb":
		extension = ".sql"
	case "postgres", "postgresql":
		extension = ".sql"
	case "mssql", "sqlserver":
		extension = ".bak"

	default:
		return "", fmt.Errorf("unsupported database driver: %s", driver)
	}

	fileName := fmt.Sprintf("database_backup_%s%s", time.Now().Format("20060102150405"), extension)

	tempDir, err := os.MkdirTemp("", "database-backup-*")

	if err != nil {
		return "", fmt.Errorf("failed to create temporary backup directory: %w", err)
	}

	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	tempFilePath := filepath.Join(tempDir, fileName)

	dbConfig := backup.Config{
		Driver:   driver,
		Host:     s.Db.GetHost(),
		Port:     s.Db.GetPort(),
		Username: s.Db.GetUsername(),
		Password: s.Db.GetPassword(),
		Database: s.Db.GetDatabase(),
		Db:       s.Db.GetDb(),
	}

	databaseBackup, err := backup.NewDatabaseBackup(dbConfig)

	if err != nil {
		return "", fmt.Errorf("failed to initialize database backup: %w", err)
	}

	if err := databaseBackup.Backup(tempFilePath); err != nil {
		return "", fmt.Errorf("failed to generate database backup: %w", err)
	}

	file, err := os.Open(tempFilePath)

	if err != nil {
		return "", fmt.Errorf("failed to open generated database backup: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	uploadedFileName, err := s.Storage.BackupUpload(fileName, file)

	if err != nil {
		return "", fmt.Errorf("failed to upload database backup: %w", err)
	}

	return uploadedFileName, nil
}

func (s *Service) GenerateFilesBackup(backupHistory *models.BackupHistory) (string, error) {

	if backupHistory == nil {
		return "", errors.New("backup history is empty")
	}

	sourceDir := strings.TrimSpace(s.Storage.GetBasePath())

	if sourceDir == "" {
		return "", errors.New("file storage base path is not configured")
	}

	sourceInfo, err := os.Stat(sourceDir)

	if err != nil {
		return "", fmt.Errorf("failed to access file storage: %w", err)
	}

	if !sourceInfo.IsDir() {
		return "", fmt.Errorf("file storage path is not a directory: %s", sourceDir)
	}

	fileName := fmt.Sprintf("files_backup_%s.zip", time.Now().Format("20060102150405"))

	tempDir, err := os.MkdirTemp("", "files-backup-*")

	if err != nil {
		return "", fmt.Errorf("failed to create temporary backup directory: %w", err)
	}

	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	tempFilePath := filepath.Join(tempDir, fileName)

	zipFile, err := os.Create(tempFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create zip file: %w", err)
	}

	zipWriter := zip.NewWriter(zipFile)

	backupDir := filepath.Join(sourceDir, "backup")

	sourceDirAbs, err := filepath.Abs(sourceDir)

	if err != nil {
		_ = zipWriter.Close()
		_ = zipFile.Close()

		return "", fmt.Errorf("failed to resolve source directory: %w", err)
	}

	backupDirAbs, err := filepath.Abs(backupDir)

	if err != nil {
		_ = zipWriter.Close()
		_ = zipFile.Close()

		return "", fmt.Errorf("failed to resolve backup directory: %w", err)
	}

	err = filepath.Walk(
		sourceDirAbs,
		func(path string, info os.FileInfo, walkErr error) error {

			if walkErr != nil {
				return walkErr
			}

			if info.IsDir() {

				currentPath, err := filepath.Abs(path)

				if err != nil {
					return err
				}

				if currentPath == backupDirAbs {
					return filepath.SkipDir
				}
				return nil
			}

			relativePath, err := filepath.Rel(sourceDirAbs, path)

			if err != nil {
				return fmt.Errorf("failed to get relative path: %w", err)
			}

			relativePath = filepath.ToSlash(relativePath)

			sourceFile, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file %s: %w", path, err)
			}

			writer, err := zipWriter.Create(relativePath)

			if err != nil {
				_ = sourceFile.Close()

				return fmt.Errorf("failed to create zip entry %s: %w", relativePath, err)
			}

			if _, err := io.Copy(writer, sourceFile); err != nil {
				_ = sourceFile.Close()

				return fmt.Errorf("failed to copy file %s: %w", path, err)
			}

			if err := sourceFile.Close(); err != nil {
				return fmt.Errorf("failed to close file %s: %w", path, err)
			}

			return nil
		},
	)

	if err != nil {
		_ = zipWriter.Close()
		_ = zipFile.Close()

		return "", fmt.Errorf("failed to create files backup: %w", err)
	}

	if err := zipWriter.Close(); err != nil {
		_ = zipFile.Close()

		return "", fmt.Errorf("failed to finalize files backup: %w", err)
	}

	if err := zipFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close files backup: %w", err)
	}

	file, err := os.Open(tempFilePath)

	if err != nil {
		return "", fmt.Errorf("failed to open generated files backup: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	uploadedFileName, err := s.Storage.BackupUpload(fileName, file)

	if err != nil {
		return "", fmt.Errorf("failed to upload files backup: %w", err)
	}

	return uploadedFileName, nil
}

func (s *Service) GenerateBackup(backupHistory *models.BackupHistory) (string, error) {

	if backupHistory == nil {
		return "", errors.New("backup history is empty")
	}

	backupType := strings.ToLower(
		strings.TrimSpace(backupHistory.BackupType),
	)

	switch backupType {

	case constants.BackupTypeDatabase:
		return s.GenerateDatabaseBackup()
	case constants.BackupTypeFiles:
		return s.GenerateFilesBackup(backupHistory)
	default:
		return "", fmt.Errorf(
			"unsupported backup type: %s",
			backupHistory.BackupType,
		)
	}
}

func (s *Service) BackupCreate(payload *requests.BackupSettingRequest, createdBy *uint64) (*models.BackupHistory, error) {

	if payload == nil {
		return nil, errors.New("backup request is empty")
	}

	backupType := payload.Backup

	if backupType == "" {
		return nil, errors.New("backup type is required")
	}

	switch backupType {

	case constants.BackupTypeDatabase, constants.BackupTypeFiles:

	default:
		return nil, fmt.Errorf("unsupported backup type: %s", payload.Backup)
	}

	backupHistory := &models.BackupHistory{
		BackupType:   backupType,
		BackupStatus: constants.BackupStatusPending,
		CreatedBy:    createdBy,
	}

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	backupHistory, err := repository.CreateBackupHistory(backupHistory)

	if err != nil {
		return nil, fmt.Errorf("failed to create backup history: %w", err)
	}

	go func(historyID uint64, history *models.BackupHistory) {

		fileName, err := s.GenerateBackup(history)

		if err != nil {
			errorLog := err.Error()
			_ = repository.UpdateBackupStatusFailed(historyID, errorLog)
			return
		}

		var fileSize uint64
		if fileName != "" {
			backupPath := filepath.Join("./storage/backup", fileName)
			if fi, err := os.Stat(backupPath); err == nil {
				fileSize = uint64(fi.Size())
			}
		}

		_ = repository.UpdateBackupStatusSuccess(historyID, fileName, fileSize)

	}(backupHistory.Id, backupHistory)

	return backupHistory, nil
}

func (s *Service) MarkBackupDownloaded(historyID uint64, uuidStr string, downloadedBy uint64, downloadedAt time.Time) error {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	return repository.MarkBackupDownloaded(historyID, uuidStr, downloadedBy, downloadedAt)
}

func (s *Service) GetBackupList(params repositories.BackupFilterParams) (map[string]any, error) {

	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > constants.MaxPageLimit {
		params.Limit = constants.DefaultPageLimit
	}

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	list, total, err := repository.GetBackupHistories(params)
	if err != nil {
		return nil, err
	}

	var items []map[string]any
	for _, item := range list {
		typeName := "Database"
		if item.BackupType == constants.BackupTypeFiles {
			typeName = "Files"
		}

		statusName := "Pending"
		switch item.BackupStatus {
		case constants.BackupStatusRunning:
			statusName = "Running"
		case constants.BackupStatusSuccess:
			statusName = "Completed"
		case constants.BackupStatusFailed:
			statusName = "Failed"
		}

		fileSizeFormatted := "-"
		if item.FileSize > 0 {
			fileSizeFormatted = formatFileSize(item.FileSize)
		} else if item.FileName != "" {
			backupPath := filepath.Join("./storage/backup", item.FileName)
			if fi, err := os.Stat(backupPath); err == nil {
				fileSizeFormatted = formatFileSize(uint64(fi.Size()))
			}
		}

		downloadedByName := "-"
		if item.DownloadedBy != nil && *item.DownloadedBy > 0 {
			if user, err := repository.GetUserByID(*item.DownloadedBy); err == nil && user != nil {
				if user.Name != nil && *user.Name != "" {
					downloadedByName = *user.Name
				} else if user.UserName != nil && *user.UserName != "" {
					downloadedByName = *user.UserName
				}
			}
		}

		downloadDate := "-"
		if item.DownloadedAt != nil && !item.DownloadedAt.IsZero() {
			downloadDate = item.DownloadedAt.Format(constants.ViewDateTimeFormat)
		}

		downloadStatusName := "Not Downloaded"
		if item.DownloadStatus == 1 || (item.DownloadedAt != nil && !item.DownloadedAt.IsZero()) {
			downloadStatusName = "Downloaded"
		}

		backupTime := "-"
		if !item.CreatedAt.IsZero() {
			backupTime = item.CreatedAt.Format(constants.ViewDateTimeFormat)
		}

		items = append(items, map[string]any{
			"id":                   item.Id,
			"uuid":                 item.UUID.String(),
			"file_name":            item.FileName,
			"backup_type":          item.BackupType,
			"backup_type_name":     typeName,
			"backup_status":        item.BackupStatus,
			"backup_status_name":   statusName,
			"download_status":      item.DownloadStatus,
			"download_status_name": downloadStatusName,
			"file_size":            item.FileSize,
			"file_size_formatted":  fileSizeFormatted,
			"error_log":            item.ErrorLog,
			"downloaded_by":        downloadedByName,
			"download_date":        downloadDate,
			"backup_time":          backupTime,
			"created_at":           backupTime,
		})
	}

	if items == nil {
		items = []map[string]any{}
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))
	if totalPages < 1 {
		totalPages = 1
	}

	return map[string]any{
		"items":       items,
		"total":       total,
		"page":        params.Page,
		"limit":       params.Limit,
		"total_pages": totalPages,
	}, nil
}

func (s *Service) GetBackupFile(id uint64) (*models.BackupHistory, string, error) {

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	backupHistory, err := repository.GetBackupHistoryByID(id)
	if err != nil {
		return nil, "", fmt.Errorf("backup record not found: %w", err)
	}

	if backupHistory.FileName == "" || backupHistory.BackupStatus != constants.BackupStatusSuccess {
		return nil, "", errors.New("backup file is not available")
	}

	backupPath := filepath.Join(s.Storage.GetBasePath(), "backup", backupHistory.FileName)
	if _, err := os.Stat(backupPath); err != nil {
		backupPath = filepath.Join("./storage/backup", backupHistory.FileName)
		if _, err := os.Stat(backupPath); err != nil {
			return nil, "", fmt.Errorf("backup file not found on disk: %w", err)
		}
	}

	return backupHistory, backupPath, nil
}

func (s *Service) GetBackupFileByUUID(uuidOrID string) (*models.BackupHistory, string, error) {

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	var backupHistory *models.BackupHistory
	var err error

	backupHistory, err = repository.GetBackupHistoryByUUID(uuidOrID)
	if err != nil {
		if numID, parseErr := strconv.ParseUint(uuidOrID, 10, 64); parseErr == nil {
			backupHistory, err = repository.GetBackupHistoryByID(numID)
		}
	}

	if err != nil || backupHistory == nil {
		return nil, "", fmt.Errorf("backup record not found: %w", err)
	}

	if backupHistory.FileName == "" || backupHistory.BackupStatus != constants.BackupStatusSuccess {
		return nil, "", errors.New("backup file is not available")
	}

	backupPath := filepath.Join(s.Storage.GetBasePath(), "backup", backupHistory.FileName)
	if _, err := os.Stat(backupPath); err != nil {
		backupPath = filepath.Join("./storage/backup", backupHistory.FileName)
		if _, err := os.Stat(backupPath); err != nil {
			return nil, "", fmt.Errorf("backup file not found on disk: %w", err)
		}
	}

	return backupHistory, backupPath, nil
}

func formatFileSize(bytes uint64) string {
	if bytes == 0 {
		return "-"
	}
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/float64(TB))
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

