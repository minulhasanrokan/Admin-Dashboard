package backup

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"gorm.io/gorm"
)

type SQLServerBackup struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	Db       *gorm.DB
}

func (b *SQLServerBackup) Backup(outputPath string) error {

	backupSQL := fmt.Sprintf(
		"BACKUP DATABASE [%s] TO DISK = N'%s' WITH INIT, FORMAT;",
		b.Database,
		outputPath,
	)

	args := []string{
		"-S",
		b.Host + "," + b.Port,
		"-U",
		b.Username,
		"-P",
		b.Password,
		"-Q",
		backupSQL,
	}

	cmd := exec.Command("sqlcmd", args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {

		if stderr.Len() > 0 {
			return fmt.Errorf(
				"sql server backup failed: %s",
				strings.TrimSpace(stderr.String()),
			)
		}

		return err
	}

	return nil
}
