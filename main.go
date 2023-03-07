/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file main.go
 * @package main
 * @author Dr.NP <np@herewe.tech>
 * @since 02/25/2023
 */

package main

import (
	"icepay-svc/handler"
	"icepay-svc/runtime"
	"os"

	"github.com/urfave/cli/v2"
)

func actionServe(c *cli.Context) error {
	handler.InitMisc()
	handler.InitClient()
	handler.InitTenant()
	handler.InitCreditCard()
	handler.InitPayment()

	return runtime.Serve()
}

func actionInitdb(c *cli.Context) error {
	// We do not initialize database here
	// Run php bin/console doctrine:schema:update --force --complete in icepay-admin
	return nil
}

// Portal

// @title icePay Demo API
// @version 0.0.1
// @description icePay Demo API
// @contact.name HereweTech CO.LTD
// @contact.url https://herewe.tech
// @contact.email support@herewetech.com

// @host api.icepay.herewe.tech
// @BasePath /
func main() {
	runtime.LoadConfig()
	runtime.InitLogger()
	runtime.InitServer()
	runtime.InitNats()
	runtime.InitDB()

	app := &cli.App{
		Name: runtime.AppName,
		Commands: []*cli.Command{
			{
				Name:   "serve",
				Usage:  "Run service",
				Action: actionServe,
			},
			{
				Name:   "initdb",
				Usage:  "Initialize database tables",
				Action: actionInitdb,
			},
		},
		DefaultCommand: "serve",
	}

	// Startup app
	if err := app.Run(os.Args); err != nil {
		runtime.Logger.Fatal(err)
	}
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
