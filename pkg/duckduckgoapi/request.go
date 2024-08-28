package duckduckgoapi

import (
	"fmt"

	"regexp"

	"github.com/charmbracelet/log"
	"github.com/go-resty/resty/v2"
)

func InitiateLogin(username string) error {

	type loginLinkResponse struct {
		C struct {
			Ar    int    `json:"ar"`
			Cp    string `json:"cp"`
			Error bool   `json:"error"`
			Flow  string `json:"flow"`
		} `json:"c"`
		Error string `json:"error"`
	}

	resp, err := resty.New().R().
		SetHeader("User-Agent", UserAgent).
		SetQueryParam("user", username).
		SetResult(&loginLinkResponse{}).
		SetError(&loginLinkResponse{}).
		Get(fmt.Sprintf("%s/auth/loginlink", Endpoint))
	if err != nil {
		return err
	}
	log.Debug("InitiateLogin", "response", resp.String())
	if resp.Error() != nil {
		return fmt.Errorf("error: %s", resp.Error().(*loginLinkResponse).Error)
	}
	return nil
}

func LoginWithOtp(username, otp string) error {

	type loginResponse struct {
		Token string `json:"token"`
	}

	formatOtpRegex := regexp.MustCompile(`\s`)
	formattedOtp := formatOtpRegex.ReplaceAllString(otp, "-")
	resp, err := resty.New().R().
		SetHeader("User-Agent", UserAgent).
		SetQueryParam("user", username).
		SetQueryParam("otp", formattedOtp).
		SetResult(&loginResponse{}).
		Get(fmt.Sprintf("%s/auth/login", Endpoint))
	if err != nil {
		return err
	}
	log.Debug("LoginWithOtp", "response", resp.String())

	return nil
}
