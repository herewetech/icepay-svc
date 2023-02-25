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
	/*
		// Check database table
		if model.CheckRecordTable() != nil {
			// InitDB
			runtime.Logger.Warn("Database does not initialized, try to create tables")
			actionInitdb(c)
		}
	*/
	handler.InitMisc()
	//handler.InitCalc()

	return runtime.Serve()
}

func actionInitdb(c *cli.Context) error {
	//model.InitRecord()

	return nil
}

// Portal
func main() {
	runtime.LoadConfig()
	runtime.InitLogger()
	runtime.InitServer()
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
		runtime.Logger.Fatal(err.Error())
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
