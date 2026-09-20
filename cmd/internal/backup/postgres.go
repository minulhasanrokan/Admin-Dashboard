package backup

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"gorm.io/gorm"
)

type PostgreSQLBackup struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	Db       *gorm.DB
}

func (b *PostgreSQLBackup) Backup(outputPath string) error {

	args := []string{
		"--host=" + b.Host,
		"--port=" + b.Port,
		"--username=" + b.Username,
		"--file=" + outputPath,
		b.Database,
	}

	cmd := exec.Command("pg_dump", args...)

	if b.Password != "" {
		cmd.Env = append(
			os.Environ(),
			"PGPASSWORD="+b.Password,
		)
	} else {
		cmd.Env = os.Environ()
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err == nil {
		return nil
	}

	if b.Db != nil {
		return b.dumpNative(outputPath)
	}

	return fmt.Errorf("pg_dump failed and native database connection not available")
}

func (b *PostgreSQLBackup) dumpNative(outputPath string) error {

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	_, _ = writer.WriteString(fmt.Sprintf("-- PostgreSQL Database Dump\n-- Generated: %s\n\n", time.Now().Format(time.RFC3339)))

	sqlDB, err := b.Db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	rows, err := sqlDB.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' AND table_type = 'BASE TABLE'")
	if err != nil {
		return fmt.Errorf("failed to query tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, tableName)
	}
	_ = rows.Close()

	for _, table := range tables {
		dataRows, err := sqlDB.Query(fmt.Sprintf(`SELECT * FROM "%s"`, table))
		if err != nil {
			continue
		}

		cols, err := dataRows.Columns()
		if err != nil {
			_ = dataRows.Close()
			continue
		}

		colCount := len(cols)
		if colCount == 0 {
			_ = dataRows.Close()
			continue
		}

		var quotedCols []string
		for _, col := range cols {
			quotedCols = append(quotedCols, fmt.Sprintf(`"%s"`, col))
		}
		colList := strings.Join(quotedCols, ", ")

		values := make([]interface{}, colCount)
		valuePtrs := make([]interface{}, colCount)
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		for dataRows.Next() {
			if err := dataRows.Scan(valuePtrs...); err != nil {
				continue
			}

			var valStrs []string
			for _, val := range values {
				if val == nil {
					valStrs = append(valStrs, "NULL")
				} else {
					switch v := val.(type) {
					case []byte:
						valStrs = append(valStrs, fmt.Sprintf("'%s'", escapeSQL(string(v))))
					case string:
						valStrs = append(valStrs, fmt.Sprintf("'%s'", escapeSQL(v)))
					case time.Time:
						valStrs = append(valStrs, fmt.Sprintf("'%s'", v.Format("2006-01-02 15:04:05")))
					default:
						valStrs = append(valStrs, fmt.Sprintf("'%v'", escapeSQL(fmt.Sprintf("%v", v))))
					}
				}
			}
			_, _ = writer.WriteString(fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s);`+"\n", table, colList, strings.Join(valStrs, ", ")))
		}
		_ = dataRows.Close()
		_, _ = writer.WriteString("\n")
	}

	return nil
}

