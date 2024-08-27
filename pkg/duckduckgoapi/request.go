package duckduckgoapi

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/go-resty/resty/v2"
)

func InitiateLogin(username string) error {

	resp, err := resty.New().R().
		SetHeader("User-Agent", Useragent).
		SetQueryParam("user", username).
		Get(fmt.Sprintf("%s/auth/loginlink", Endpoint))
	if err != nil {
		return err
	}
	log.Debug("InitiateLogin", "response", resp.String())
	return nil
}
