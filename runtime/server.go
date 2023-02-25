/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file server.go
 * @package runtime
 * @author Dr.NP <np@herewe.tech>
 * @since 11/17/2022
 */

package runtime

import (
	"os"

	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/pkg/errors"
)

var Server *fiber.App

func InitServer() error {
	app := fiber.New(fiber.Config{
		ServerHeader:          AppName,
		DisableKeepalive:      false,
		AppName:               AppName,
		Prefork:               Config.HTTP.Prefork,
		DisableStartupMessage: true,
	})
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: LoggerRaw,
	}))
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(requestid.New())

	Server = app

	return nil
}

func Serve() error {
	if Server == nil {
		// Not initialized
		return errors.New("Server not initialized")
	}

	Logger.Infof("starting HTTP server on [%s]", Config.HTTP.ListenAddr)

	return Server.Listen(Config.HTTP.ListenAddr)
}

func Exit() {
	// TODO: Pure runtime
	os.Exit(-1)
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
