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

type MySQLBackup struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	Db       *gorm.DB
}

func (b *MySQLBackup) Backup(outputPath string) error {

	host := b.Host
	if host == "localhost" {
		host = "127.0.0.1"
	}

	args := []string{
		"--host=" + host,
		"--port=" + b.Port,
		"--user=" + b.Username,
		"--column-statistics=0",
		"--single-transaction",
		"--quick",
		b.Database,
	}

	file, err := os.Create(outputPath)
	if err == nil {
		cmd := exec.Command("mysqldump", args...)

		if b.Password != "" {
			cmd.Env = append(
				os.Environ(),
				"MYSQL_PWD="+b.Password,
			)
		} else {
			cmd.Env = os.Environ()
		}

		cmd.Stdout = file

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if runErr := cmd.Run(); runErr == nil {
			_ = file.Close()
			return nil
		}
		_ = file.Close()
	}

	if b.Db != nil {
		return b.dumpNative(outputPath)
	}

	return fmt.Errorf("mysqldump failed and native database connection not available")
}

func (b *MySQLBackup) dumpNative(outputPath string) error {

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	writer := bufio.NewWriter(file)
	defer func(writer *bufio.Writer) {
		err := writer.Flush()
		if err != nil {

		}
	}(writer)

	_, _ = writer.WriteString(fmt.Sprintf("-- MySQL Database Dump\n-- Generated: %s\n\nSET FOREIGN_KEY_CHECKS = 0;\n\n", time.Now().Format(time.RFC3339)))

	sqlDB, err := b.Db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	rows, err := sqlDB.Query("SHOW FULL TABLES WHERE Table_type = 'BASE TABLE'")
	if err != nil {
		return fmt.Errorf("failed to query tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName, tableType string
		if err := rows.Scan(&tableName, &tableType); err != nil {
			return fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, tableName)
	}
	_ = rows.Close()

	for _, table := range tables {
		_, _ = writer.WriteString(fmt.Sprintf("--\n-- Table structure for table `%s`\n--\n\n", table))
		_, _ = writer.WriteString(fmt.Sprintf("DROP TABLE IF EXISTS `%s`;\n", table))

		var tableName, createSQL string
		row := sqlDB.QueryRow(fmt.Sprintf("SHOW CREATE TABLE `%s`", table))
		if err := row.Scan(&tableName, &createSQL); err == nil {
			_, _ = writer.WriteString(createSQL + ";\n\n")
		}

		dataRows, err := sqlDB.Query(fmt.Sprintf("SELECT * FROM `%s`", table))
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

		var insertStatements []string
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
			insertStatements = append(insertStatements, fmt.Sprintf("(%s)", strings.Join(valStrs, ", ")))

			if len(insertStatements) >= 100 {
				_, _ = writer.WriteString(fmt.Sprintf("INSERT INTO `%s` VALUES %s;\n", table, strings.Join(insertStatements, ",\n")))
				insertStatements = insertStatements[:0]
			}
		}
		_ = dataRows.Close()

		if len(insertStatements) > 0 {
			_, _ = writer.WriteString(fmt.Sprintf("INSERT INTO `%s` VALUES %s;\n", table, strings.Join(insertStatements, ",\n")))
		}
		_, _ = writer.WriteString("\n")
	}

	_, _ = writer.WriteString("SET FOREIGN_KEY_CHECKS = 1;\n")
	return nil
}

func escapeSQL(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "'", "\\'")
	s = strings.ReplaceAll(s, "\000", "\\0")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	return s
}
