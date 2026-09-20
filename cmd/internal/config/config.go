package config

import (
	"admin-dashboard/cmd/common/helpers"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName                  string
	AppKey                   string
	AppIp                    string
	AppDomain                string
	AppPort                  string
	AppBaseUrl               string
	AssetIp                  string
	AssetDomain              string
	AssetPort                string
	AssetBaseUrl             string
	DbConnection             string
	DbRecordNotFoundMessage  string
	DbHost                   string
	DbPort                   string
	DbDatabase               string
	DbUsername               string
	DbPassword               string
	MailDriver               string
	MailSenderName           string
	MailSMTPHost             string
	MailSMTPPort             int
	MailSMTPUsername         string
	MailSMTPPassword         string
	MailSMTPFrom             string
	MailSendGridAPIKey       string
	MailSendGridFrom         string
	MailAWSRegion            string
	MailSESFrom              string
	AWSAccessKey             string
	AWSSecretKey             string
	WebCookieName            string
	WebCookieDomain          string
	WebCookiePath            string
	WebCookieSecure          bool
	WebCookieHTTPOnly        bool
	WebCookieSameSite        string
	WebCookieMaxAge          int
	WebCSRFTokenLookup       string
	WebCSRFTokenLength       int
	WebXSSProtection         string
	WebContentTypeNosniff    string
	WebXFrameOptions         string
	WebContentSecurityPolicy string
	WebHSTSMaxAge            int
	WebHSTSExcludeSubdomains bool
	WebHSTSPreloadEnabled    bool
	WebReferrerPolicy        string
	WebCSPReportOnly         bool
	WebRateLimit             int
	VerificationTokenLength  int
	VerificationTokenMinute  int
	PasswordResetTokenLength int
	PasswordResetTokenMinute int
	FileSystem               string
	FileSystemLocation       string
	FileBackupLocation       string
	FileSystemAWSRegion      string
	FileSystemS3Bucket       string
	CacheDriver              string
	CacheTime                int
	RedisHost                string
	RedisPort                string
	RedisPassword            string
	RedisDb                  string
	CachePath                string
	SessionDriver            string
	SessionLifetime          int
	SessionPath              string
	SessionTable             string
	QueueWorkerInterval      int
	QueueWorkerBatchSize     int
	QueueWorkerMaxTries      int
}

func Init() (*Config, error) {
	err := godotenv.Load()

	if err != nil {

		return nil, err
	}

	mailSMTPPort, err := helpers.StrToInt(os.Getenv("MAIL_SMTP_PORT"))

	if err != nil {
		return nil, err
	}

	webMaxAge, err := helpers.StrToInt(os.Getenv("WEB_COOKIE_MAX_AGE"))

	if err != nil {
		return nil, err
	}

	webTokenLength, err := helpers.StrToInt(os.Getenv("WEB_CSRF_TOKEN_LENGTH"))

	if err != nil {
		return nil, err
	}

	webHstsMaxAge, err := helpers.StrToInt(os.Getenv("WEB_SECURE_HSTS_MAX_AGE"))

	if err != nil {
		return nil, err
	}

	webHstsExclude, err := helpers.StrToBool(os.Getenv("WEB_SECURE_HSTS_SUBDOMAIN"))

	if err != nil {
		return nil, err
	}

	webHstsPreload, err := helpers.StrToBool(os.Getenv("WEB_SECURE_HSTS_PRELOAD"))

	if err != nil {
		return nil, err
	}

	webRateLimit, err := helpers.StrToInt(os.Getenv("WEB_RATE_LIMIT"))

	if err != nil {
		return nil, err
	}

	webCookieSecure, err := helpers.StrToBool(os.Getenv("WEB_COOKIE_SECURE"))

	if err != nil {
		return nil, err
	}

	webCookieHTTPOnly, err := helpers.StrToBool(os.Getenv("WEB_COOKIE_HTTP_ONLY"))

	if err != nil {
		return nil, err
	}

	webCSPReportOnly, err := helpers.StrToBool(os.Getenv("WEB_CSP_REPORT_ONLY"))

	if err != nil {
		return nil, err
	}

	verificationTokenLength, err := helpers.StrToInt(os.Getenv("VERIFICATION_TOKEN_LENGTH"))

	if err != nil {
		return nil, err
	}

	verificationTokenMinute, err := helpers.StrToInt(os.Getenv("VERIFICATION_TOKEN_EXPIRE_MINUTE"))

	if err != nil {
		return nil, err
	}

	passwordResetTokenLength, err := helpers.StrToInt(os.Getenv("PASSWORD_RESET_TOKEN_LENGTH"))

	if err != nil {
		return nil, err
	}

	passwordResetTokenMinute, err := helpers.StrToInt(os.Getenv("PASSWORD_RESET_TOKEN_EXPIRE_MINUTE"))

	if err != nil {
		return nil, err
	}

	cacheTime, err := helpers.StrToInt(os.Getenv("CACHE_TIME"))
	if err != nil {
		cacheTime = 0
	}

	sessionLifetime, err := helpers.StrToInt(os.Getenv("SESSION_LIFETIME"))
	if err != nil || sessionLifetime <= 0 {
		sessionLifetime = 120
	}

	sessionDriver := os.Getenv("SESSION_DRIVER")
	if sessionDriver == "" {
		sessionDriver = "file"
	}

	sessionPath := os.Getenv("SESSION_PATH")
	if sessionPath == "" {
		sessionPath = "storage/sessions"
	}

	sessionTable := os.Getenv("SESSION_TABLE")
	if sessionTable == "" {
		sessionTable = "sessions"
	}

	queueWorkerInterval, err := helpers.StrToInt(os.Getenv("QUEUE_WORKER_INTERVAL"))
	if err != nil || queueWorkerInterval <= 0 {
		queueWorkerInterval = 5
	}

	queueWorkerBatchSize, err := helpers.StrToInt(os.Getenv("QUEUE_WORKER_BATCH_SIZE"))
	if err != nil || queueWorkerBatchSize <= 0 {
		queueWorkerBatchSize = 50
	}

	queueWorkerMaxTries, err := helpers.StrToInt(os.Getenv("QUEUE_WORKER_MAX_TRIES"))
	if err != nil || queueWorkerMaxTries <= 0 {
		queueWorkerMaxTries = 5
	}

	cfg := &Config{
		AppKey:                   os.Getenv("APP_KEY"),
		AppName:                  os.Getenv("APP_NAME"),
		AppIp:                    os.Getenv("APP_IP"),
		AppDomain:                os.Getenv("APP_DOMAIN"),
		AppPort:                  os.Getenv("APP_PORT"),
		AppBaseUrl:               os.Getenv("APP_BASE_URL"),
		AssetIp:                  os.Getenv("ASSET_IP"),
		AssetDomain:              os.Getenv("ASSET_DOMAIN"),
		AssetPort:                os.Getenv("ASSET_PORT"),
		AssetBaseUrl:             os.Getenv("ASSET_BASE_URL"),
		DbConnection:             os.Getenv("DB_CONNECTION"),
		DbRecordNotFoundMessage:  os.Getenv("DB_RECORD_NOT_FOUND_MESSAGE"),
		DbHost:                   os.Getenv("DB_HOST"),
		DbPort:                   os.Getenv("DB_PORT"),
		DbDatabase:               os.Getenv("DB_DATABASE"),
		DbUsername:               os.Getenv("DB_USERNAME"),
		DbPassword:               os.Getenv("DB_PASSWORD"),
		MailDriver:               os.Getenv("MAIL_DRIVER"),
		MailSenderName:           os.Getenv("MAIL_SENDER_NAME"),
		MailSMTPHost:             os.Getenv("MAIL_SMTP_HOST"),
		MailSMTPPort:             mailSMTPPort,
		MailSMTPUsername:         os.Getenv("MAIL_SMTP_USERNAME"),
		MailSMTPPassword:         os.Getenv("MAIL_SMTP_PASSWORD"),
		MailSMTPFrom:             os.Getenv("MAIL_SMTP_FROM"),
		MailSendGridAPIKey:       os.Getenv("MAIL_SENDGRID_API_KEY"),
		MailSendGridFrom:         os.Getenv("MAIL_SENDGRID_FROM"),
		MailAWSRegion:            os.Getenv("MAIL_AWS_REGION"),
		MailSESFrom:              os.Getenv("MAIL_SES_FROM"),
		AWSAccessKey:             os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretKey:             os.Getenv("AWS_SECRET_ACCESS_KEY"),
		WebCookieName:            os.Getenv("WEB_COOKIE_NAME"),
		WebCookieDomain:          os.Getenv("WEB_COOKIE_DOMAIN"),
		WebCookiePath:            os.Getenv("WEB_COOKIE_PATH"),
		WebCookieSecure:          webCookieSecure,
		WebCookieHTTPOnly:        webCookieHTTPOnly,
		WebCookieSameSite:        os.Getenv("WEB_COOKIE_SAME_SITE"),
		WebCookieMaxAge:          webMaxAge,
		WebCSRFTokenLookup:       os.Getenv("WEB_CSRF_TOKEN_LOOKUP"),
		WebCSRFTokenLength:       webTokenLength,
		WebXSSProtection:         os.Getenv("WEB_SECURE_XSS"),
		WebContentTypeNosniff:    os.Getenv("WEB_SECURE_NOSNIFF"),
		WebXFrameOptions:         os.Getenv("WEB_SECURE_FRAME_OPTIONS"),
		WebContentSecurityPolicy: os.Getenv("WEB_CONTENT_SECURITY_POLICY"),
		WebHSTSMaxAge:            webHstsMaxAge,
		WebHSTSExcludeSubdomains: webHstsExclude,
		WebHSTSPreloadEnabled:    webHstsPreload,
		WebReferrerPolicy:        os.Getenv("WEB_SECURE_REFERRER_POLICY"),
		WebCSPReportOnly:         webCSPReportOnly,
		WebRateLimit:             webRateLimit,
		VerificationTokenLength:  verificationTokenLength,
		VerificationTokenMinute:  verificationTokenMinute,
		PasswordResetTokenLength: passwordResetTokenLength,
		PasswordResetTokenMinute: passwordResetTokenMinute,
		FileSystem:               os.Getenv("FILE_SYSTEM"),
		FileSystemLocation:       os.Getenv("FILE_SYSTEM_LOCATION"),
		FileBackupLocation:       os.Getenv("FILE_BACKUP_LOCATION"),
		FileSystemAWSRegion:      os.Getenv("FILE_SYSTEM_AWS_REGION"),
		FileSystemS3Bucket:       os.Getenv("FILE_SYSTEM_S3_BUCKET"),
		CacheDriver:              os.Getenv("CACHE_DRIVER"),
		CacheTime:                cacheTime,
		RedisHost:                os.Getenv("REDIS_HOST"),
		RedisPort:                os.Getenv("REDIS_PORT"),
		RedisPassword:            os.Getenv("REDIS_PASSWORD"),
		RedisDb:                  os.Getenv("REDIS_DB"),
		CachePath:                os.Getenv("CACHE_PATH"),
		SessionDriver:            sessionDriver,
		SessionLifetime:          sessionLifetime,
		SessionPath:              sessionPath,
		SessionTable:             sessionTable,
		QueueWorkerInterval:      queueWorkerInterval,
		QueueWorkerBatchSize:     queueWorkerBatchSize,
		QueueWorkerMaxTries:      queueWorkerMaxTries,
	}

	return cfg, nil
}
