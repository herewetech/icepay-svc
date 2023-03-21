/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file config.go
 * @package runtime
 * @author Dr.NP <np@herewe.tech>
 * @since 02/25/2023
 */

package runtime

import (
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

const (
	AppName   = "svc"
	EnvPrefix = "icepay"
)

type mainConfig struct {
	HTTP struct {
		ListenAddr         string `json:"listen_addr" mapstructure:"listen_addr"`
		Prefork            bool   `json:"prefork" mapstructure:"prefork"`
		LongPollingTimeout int64  `json:"long_polling_timeout" mapstructure:"long_polling_timeout"` // In second
	} `json:"http" mapstructure:"http"`
	Database struct {
		DSN string `json:"dsn" mapstructure:"dsn"`
	} `json:"database" mapstructure:"database"`
	Nats struct {
		URL string `json:"url" mapstructure:"url"`
	} `json:"nats" mapstructure:"nats"`
	Auth struct {
		JWTAccessSecret  string `json:"jwt_access_secret" mapstructure:"jwt_access_secret"`
		JWTRefreshSecret string `json:"jwt_refresh_secret" mapstructure:"jwt_refresh_secret"`
		JWTAccessExpiry  int64  `json:"jwt_access_expiry" mapstructure:"jwt_access_expiry"`   // In minute
		JWTRefreshExpiry int64  `json:"jwt_refresh_expiry" mapstructure:"jwt_refresh_expiry"` // In minute
	} `json:"auth" mapstructure:"auth"`
	Security struct {
		CredentialLifetime int64  `json:"credential_lifetime" mapstructure:"credential_lifetime"` // In minute
		AESKey             string `json:"aes_key" mapstructure:"aes_key"`
	} `json:"security" mapstructure:"security"`
	Firebase struct {
		Credentials struct {
			Type                    string `json:"id" mapstructure:"id"`
			ProjectID               string `json:"project_id" mapstructure:"project_id"`
			PrivateKeyID            string `json:"private_key_id" mapstructure:"private_key_id"`
			PrivateKey              string `json:"private_key" mapstructure:"private_key"`
			ClientEmail             string `json:"client_email" mapstructure:"client_email"`
			ClientID                string `json:"client_id" mapstructure:"client_id"`
			AuthURL                 string `json:"auth_url" mapstructure:"auth_url"`
			TokenURL                string `json:"token_url" mapstructure:"token_url"`
			AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url" mapstructure:"auth_provider_x509_cert_url"`
			ClientX509CertURL       string `json:"client_x509_cert_url" mapstructure:"client_x509_cert_url"`
		} `json:"credentials" mapstructure:"credentials"`
		CredentialsFile string `json:"credentials_file" mapstructure:"credentials_file"`
	} `json:"firebase" mapstructure:"firebase"`
	Debug bool `json:"debug" mapstructure:"debug"`
}

var Config mainConfig

var defaultConfigs = map[string]interface{}{
	"http.listen_addr":                                 ":9900",
	"http.prefork":                                     false,
	"http.long_polling_timeout":                        30,
	"database.dsn":                                     "postgres://icepay@localhost:5432/icepay?sslmode=disable",
	"nats.url":                                         nats.DefaultURL,
	"auth.jwt_access_secret":                           "access_secret",
	"auth.jwt_refresh_secret":                          "refresh_secret",
	"auth.jwt_access_expiry":                           10,
	"auth.jwt_refresh_expiry":                          43200,
	"security.aes_key":                                 "icepay@@20130920",
	"security.credential_lifetime":                     5,
	"firebase.credentials.type":                        "service_account",
	"firebase.credentials.auth_url":                    "https://accounts.google.com/o/oauth2/auth",
	"firebase.credentials.token_url":                   "https://oauth2.googleapis.com/token",
	"firebase.credentials.auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
	"firebase.credentials_file":                        "./firebase.json",
	"debug":                                            true,
}

func LoadConfig() error {
	for cfgKey, cfgVal := range defaultConfigs {
		viper.SetDefault(cfgKey, cfgVal)
	}

	viper.SetConfigName(AppName)
	viper.SetConfigType("json")
	viper.AddConfigPath("/etc/" + AppName)
	viper.AddConfigPath("$HOME/." + AppName)
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		// Logging
		Logger.Warnf("reading config file error : %s", err)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix(EnvPrefix)
	viper.AutomaticEnv()
	err = viper.Unmarshal(&Config)
	if err != nil {
		LoggerRaw.Fatal(err.Error())
	}

	return err
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
