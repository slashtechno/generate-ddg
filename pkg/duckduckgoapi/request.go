package duckduckgoapi

import (
	"errors"
	"fmt"

	"regexp"

	"github.com/charmbracelet/log"
	"github.com/go-resty/resty/v2"
)

func GetEmail(accessToken string) (string, error) {

	type emailResponse struct {
		Address string `json:"address"`
	}

	resp, err := resty.New().R().
		SetHeader("User-Agent", UserAgent).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
		SetResult(&emailResponse{}).
		Post(fmt.Sprintf("%s/email/addresses", Endpoint))
	if err != nil {
		return "", err
	}

	log.Debug("GetEmail", "response", resp.String())
	if resp.StatusCode() != 201 {
		return "", fmt.Errorf("error: %s; status code: %d", resp.String(), resp.StatusCode())
	}

	return resp.Result().(*emailResponse).Address, nil
}

func GetAccessToken(refreshToken string) (string, error) {
	type accessTokenResponse struct {
		// Invites []any `json:"invites"`
		// Stats struct {
		// AddressesGenerated int `json:"addresses_generated"`
		// } `json:"stats"`
		User struct {
			AccessToken string `json:"access_token"`
			// Email       string `json:"email"`
			// Username    string `json:"username"`
		} `json:"user"`
	}

	resp, err := resty.New().R().
		SetHeader("User-Agent", UserAgent).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", refreshToken)).
		SetResult(&accessTokenResponse{}).
		Get(fmt.Sprintf("%s/email/dashboard", Endpoint))
	if err != nil {
		return "", err
	}

	log.Debug("GetAccessToken", "response", resp.String())
	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("error: %s; status code: %d", resp.String(), resp.StatusCode())
	}

	return resp.Result().(*accessTokenResponse).User.AccessToken, nil
}

func InitiateLogin(username string) error {
	// TODO: Handle the "rc" error better
	type returnedError struct {
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
		SetError(&returnedError{}).
		Get(fmt.Sprintf("%s/auth/loginlink", Endpoint))
	if err != nil {
		return err
	}
	log.Debug("InitiateLogin", "response", resp.String())

	if resp.Error() != nil {
		if resp.Error().(*returnedError).Error == "rc" {
			return errors.New("rate limited; login via the browser/app and pass the OTP via `--otp`")
		}
	} else {
		if resp.StatusCode() != 200 {
			return fmt.Errorf("unable to initiate login; try to login via the browser/app and pass the OTP via `--otp`: %s", resp.String())
		}
	}
	return nil
}

func LoginWithOtp(username, otp string) (string, error) {

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
		return "", err
	}
	log.Debug("LoginWithOtp", "response", resp.String())
	if resp.Error() != nil {
		return "", fmt.Errorf("error: %s; status code: %d", resp.String(), resp.StatusCode())
	}

	return resp.Result().(*loginResponse).Token, nil
}
