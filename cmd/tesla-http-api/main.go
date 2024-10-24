package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/L480/tesla-http-api/internal/logger"
	"github.com/L480/tesla-http-api/internal/tesla"
	"github.com/teslamotors/vehicle-command/pkg/protocol"
	"github.com/teslamotors/vehicle-command/pkg/proxy"
)

const (
	EnvRefreshToken    = "TESLA_REFRESH_TOKEN"
	EnvClientId        = "TESLA_CLIENT_ID"
	EnvApiTokenEnabled = "ENABLE_API_TOKEN"
	EnvApiToken        = "API_TOKEN"
	cacheSize          = 10000 // Number of cached vehicle sessions
	addr               = "0.0.0.0:8080"
	timeout            = "30s"
)

var (
	config          tesla.Config
	apiTokenEnabled bool
	apiToken        string
)

func middleware(next http.Handler, teslaAccessToken string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if apiTokenEnabled {
			token := r.Header.Get("Authorization")
			if token != apiToken {
				logger.Info("Request to %s from %s \033[31m(invalid token)\033[0m", r.URL.Path, r.Header.Get("X-Forwarded-For"))
				http.Error(w, http.StatusText(403), http.StatusForbidden)
				return
			}
			r.Header.Del("Authorization")
		}

		r.Header.Add("Authorization", "Bearer "+teslaAccessToken)
		logger.Info("Request to %s from %s", r.URL.Path, r.Header.Get("X-Forwarded-For"))
		next.ServeHTTP(w, r)
	})
}

func main() {
	logger.SetLevel(logger.LevelDebug)
	logger.Info("tesla-http-api")
	err := readFromEnvironment()
	if err != nil {
		logger.Error("Invalid configuration: %s", err)
		os.Exit(1)
	}

	if !apiTokenEnabled {
		logger.Warning("\033[33m%s IS SET TO FALSE. YOUR API IS UNPROTECTED AND CAN BE USED WITHOUT AUTHENTICATION. THIS IS NOT RECOMMENDED.\033[0m", EnvApiTokenEnabled)
	}

	go tesla.RefreshToken(config)

	key, err := protocol.LoadPrivateKey(config.PrivateKeyFile)
	if err != nil {
		logger.Error("Failed to load private key: %s", err)
		os.Exit(1)
	}

	p, err := proxy.New(context.Background(), key, cacheSize)
	if err != nil {
		return
	}

	timeout, _ := time.ParseDuration(timeout)
	p.Timeout = timeout

	accessToken, err := os.ReadFile(config.AccessTokenFile)
	logger.Info("Listening on %s", addr)
	logger.Error("Server stopped: %s", http.ListenAndServe(addr, middleware(p, string(accessToken))))
}

func readFromEnvironment() error {
	config.RefreshTokenFile = "./refresh-token"
	config.AccessTokenFile = "./access-token"
	config.PrivateKeyFile = "./private-key.pem"

	if EnvRefreshTokenValue, ok := os.LookupEnv(EnvRefreshToken); ok {
		config.RefreshToken = EnvRefreshTokenValue
	} else {
		return fmt.Errorf("environment variable %s is missing", EnvRefreshToken)
	}

	if EnvClientIdValue, ok := os.LookupEnv(EnvClientId); ok {
		config.ClientId = EnvClientIdValue
	} else {
		return fmt.Errorf("environment variable %s is missing", EnvClientId)
	}

	var err error
	if EnvApiTokenEnabledValue, ok := os.LookupEnv(EnvApiTokenEnabled); ok {
		apiTokenEnabled, err = strconv.ParseBool(EnvApiTokenEnabledValue)
		if err != nil {
			return fmt.Errorf("invalid value for environment variable %s", EnvApiTokenEnabled)
		}
	} else {
		return fmt.Errorf("environment variable %s is missing", EnvApiTokenEnabled)
	}

	if apiTokenEnabled {
		if EnvApiTokenValue, ok := os.LookupEnv(EnvApiToken); ok {
			apiToken = "Bearer " + EnvApiTokenValue
		} else {
			return fmt.Errorf("environment variable %s is missing", EnvApiToken)
		}
	}

	return nil
}
