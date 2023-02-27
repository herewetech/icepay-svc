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

	"github.com/spf13/viper"
)

const (
	AppName   = "svc"
	EnvPrefix = "icepay"
)

type mainConfig struct {
	HTTP struct {
		ListenAddr string `json:"listen_addr" mapstructure:"listen_addr"`
		Prefork    bool   `json:"prefork" mapstructure:"prefork"`
	} `json:"http" mapstructure:"http"`
	Database struct {
		DSN string `json:"dsn" mapstructure:"dsn"`
	} `json:"database" mapstructure:"database"`
	Auth struct {
		JWTSecret string `json:"jwt_secret" mapstructure:"jwt_secret"`
		JWTExpiry int64  `json:"jwt_expiry" mapstructure:"jwt_expiry"`
	} `json:"auth" mapstructure:"auth"`
	Debug bool `json:"debug" mapstructure:"debug"`
}

var Config mainConfig

var defaultConfigs = map[string]interface{}{
	"http.listen_addr": ":9900",
	"http.prefork":     false,
	"database.dsn":     "postgres://icepay@localhost:5432/icepay?sslmode=disable",
	"auth.jwt_secret":  "secret",
	"auth.jwt_expiry":  1440,
	"debug":            true,
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
