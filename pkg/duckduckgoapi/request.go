package duckduckgoapi

import (
	"fmt"

	"regexp"

	"github.com/charmbracelet/log"
	"github.com/go-resty/resty/v2"
)

func InitiateLogin(username string) error {

	resp, err := resty.New().R().
		SetHeader("User-Agent", UserAgent).
		SetQueryParam("user", username).
		Get(fmt.Sprintf("%s/auth/loginlink", Endpoint))
	if err != nil {
		return err
	}
	log.Debug("InitiateLogin", "response", resp.String())
	return nil
}

func LoginWithOtp(username, otp string) error {

	formatOtpRegex := regexp.MustCompile(`\s`)
	formattedOtp := formatOtpRegex.ReplaceAllString(otp, "-")
	resp, err := resty.New().R().
		SetHeader("User-Agent", UserAgent).
		SetQueryParam("user", username).
		SetQueryParam("otp", formattedOtp).
		Get(fmt.Sprintf("%s/auth/login", Endpoint))
	if err != nil {
		return err
	}
	log.Debug("LoginWithOtp", "response", resp.String())
	return nil
}
