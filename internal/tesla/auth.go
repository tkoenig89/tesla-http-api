package tesla

import (
	"encoding/json"
	"io"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/L480/tesla-http-api/internal/logger"
	"github.com/L480/tesla-http-api/internal/request"
)

type Config struct {
	PrivateKeyFilePath   string
	RefreshTokenFilePath string
	ClientId             string
	RefreshToken         string
}

var (
	Healthy     bool
	AccessToken string
)

func RefreshToken(c Config) {
	tokenTimer := time.NewTimer(0)
	refreshInterval := 6 * time.Hour
	retryInterval := 5 * time.Minute
	var latestRefreshToken string

	type response struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		IDToken      string `json:"id_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
	}

	for {
		<-tokenTimer.C
		logger.Info("Refreshing access token ...")
		file, err := os.ReadFile(c.RefreshTokenFilePath)

		if err != nil {
			latestRefreshToken = c.RefreshToken
		} else {
			latestRefreshToken = string(file)
		}

		form := url.Values{}
		form.Add("grant_type", "refresh_token")
		form.Add("client_id", c.ClientId)
		form.Add("refresh_token", latestRefreshToken)
		tokenEndpoint := request.Endpoint{
			Url:                "https://auth.tesla.com/oauth2/v3/token",
			Method:             "POST",
			Headers:            map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
			Body:               form.Encode(),
			InsecureSkipVerify: false,
		}

		resp, err := request.Connect(tokenEndpoint)
		if err != nil {
			Healthy = false
			logger.Error("Failed to connect to token endpoint: %s", err)
			tokenTimer.Reset(retryInterval)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			Healthy = false
			logger.Error("Failed to retrieve response body: %s", err)
			tokenTimer.Reset(retryInterval)
			continue
		}

		if resp.StatusCode == 200 {
			Healthy = true
			logger.Info("Access token refresh successful")
			logger.Info("Next access token refresh scheduled in %s hours", strconv.FormatFloat(refreshInterval.Hours(), 'g', 2, 64))
			tokenTimer.Reset(refreshInterval)
		} else {
			Healthy = false
			logger.Error("Refresh failed: %s", string(body))
			logger.Info("Retrying access token refresh in %s minutes", strconv.FormatFloat(retryInterval.Minutes(), 'g', 2, 64))
			tokenTimer.Reset(retryInterval)
			continue
		}

		var jsonData response
		json.Unmarshal(body, &jsonData)
		AccessToken = jsonData.AccessToken
		RefreshTokenFilePath, err := os.Create(c.RefreshTokenFilePath)
		if err != nil {
			Healthy = false
			logger.Error("Failed to save refresh token: %s", err)
			tokenTimer.Reset(retryInterval)
			continue
		}
		RefreshTokenFilePath.WriteString(jsonData.RefreshToken)
		RefreshTokenFilePath.Close()
	}
}
