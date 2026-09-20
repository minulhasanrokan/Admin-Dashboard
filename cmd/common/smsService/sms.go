package smsService

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type SslSmsRequest struct {
	ApiToken string `json:"api_token"`
	Sid      string `json:"sid"`
	Msisdn   string `json:"msisdn"`
	Sms      string `json:"sms"`
	CsmsId   string `json:"csms_id"`
}

type SslSmsInfoItem struct {
	SmsStatus    string `json:"sms_status"`
	StatusMsisdn string `json:"status_msisdn"`
	CsmsId       string `json:"csms_id"`
	ReferenceId  string `json:"reference_id"`
}

type SslSmsResponse struct {
	Status       string           `json:"status"`
	StatusCode   int              `json:"status_code"`
	ErrorMessage string           `json:"error_message"`
	SmsInfo      []SslSmsInfoItem `json:"smsinfo"`
}

func FormatMsisdn(mobile string) string {
	reg := regexp.MustCompile(`[^0-9+]`)
	cleaned := reg.ReplaceAllString(strings.TrimSpace(mobile), "")

	if strings.HasPrefix(cleaned, "+88") {
		cleaned = cleaned[1:]
	} else if strings.HasPrefix(cleaned, "01") {
		cleaned = "88" + cleaned
	} else if strings.HasPrefix(cleaned, "+") {
		cleaned = cleaned[1:]
	}

	return cleaned
}

func GenerateCsmsId(prefix string) string {
	if prefix == "" {
		prefix = "ADM"
	}
	return fmt.Sprintf("%s%d%04d", prefix, time.Now().Unix(), time.Now().Nanosecond()%10000)
}

func SendSslCommerzSms(apiUrl, apiToken, sid, csmsIdPrefix, mobile, message string) (*SslSmsResponse, error) {

	if strings.TrimSpace(apiUrl) == "" {
		apiUrl = "https://smsplus.sslwireless.com/api/v3/send-sms"
	}
	if strings.TrimSpace(apiToken) == "" {
		return nil, errors.New("SSLCommerz SMS API Token is missing. Please configure it in SMS Settings")
	}
	if strings.TrimSpace(sid) == "" {
		return nil, errors.New("SSLCommerz SMS SID is missing. Please configure it in SMS Settings")
	}

	msisdn := FormatMsisdn(mobile)
	if msisdn == "" {
		return nil, errors.New("Invalid recipient mobile number")
	}

	csmsId := GenerateCsmsId(csmsIdPrefix)

	reqPayload := SslSmsRequest{
		ApiToken: apiToken,
		Sid:      sid,
		Msisdn:   msisdn,
		Sms:      message,
		CsmsId:   csmsId,
	}

	jsonData, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode SMS payload: %w", err)
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	httpReq, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request to SSLCommerz failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response from SSLCommerz: %w", err)
	}

	var sslResp SslSmsResponse
	if err := json.Unmarshal(bodyBytes, &sslResp); err != nil {
		return nil, fmt.Errorf("failed to parse SSLCommerz response (%s): %w", string(bodyBytes), err)
	}

	if sslResp.Status != "SUCCESS" && sslResp.StatusCode != 200 && sslResp.StatusCode != 0 {
		errMsg := sslResp.ErrorMessage
		if errMsg == "" {
			errMsg = string(bodyBytes)
		}
		return &sslResp, fmt.Errorf("SSLCommerz error: %s", errMsg)
	}

	return &sslResp, nil
}
