package constants

type SelectOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

const (
	MailDriverSendgrid = "sendgrid"
	MailDriverSmtp     = "smtp"
	MailDriverSes      = "ses"
)

var AwsRegions = []SelectOption{
	{Value: "us-east-1", Label: "US East (N. Virginia) us-east-1"},
	{Value: "us-east-2", Label: "US East (Ohio) us-east-2"},
	{Value: "us-west-1", Label: "US West (N. California) us-west-1"},
	{Value: "us-west-2", Label: "US West (Oregon) us-west-2"},
	{Value: "ap-southeast-1", Label: "Asia Pacific (Singapore) ap-southeast-1"},
	{Value: "ap-southeast-2", Label: "Asia Pacific (Sydney) ap-southeast-2"},
	{Value: "ap-northeast-1", Label: "Asia Pacific (Tokyo) ap-northeast-1"},
	{Value: "ap-south-1", Label: "Asia Pacific (Mumbai) ap-south-1"},
	{Value: "eu-west-1", Label: "Europe (Ireland) eu-west-1"},
	{Value: "eu-central-1", Label: "Europe (Frankfurt) eu-central-1"},
}

var SmtpEncryptionProtocols = []SelectOption{
	{Value: "tls", Label: "TLS (Port 587 - Recommended)"},
	{Value: "ssl", Label: "SSL (Port 465)"},
	{Value: "starttls", Label: "STARTTLS"},
	{Value: "none", Label: "None (Plain / Port 25)"},
}

func GetAwsRegions() []SelectOption {
	return AwsRegions
}

func GetSmtpEncryptionProtocols() []SelectOption {
	return SmtpEncryptionProtocols
}
