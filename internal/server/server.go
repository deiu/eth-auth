package server

import (
	"errors"
	"os"

	"github.com/deiu/eth-auth/internal/server/handlers"
	"github.com/deiu/eth-auth/pkg/auth"
	"github.com/deiu/eth-auth/pkg/resolver"
	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	jwtware "github.com/gofiber/jwt"
)

// Config contains the configuration options for the server
type Config struct {
	InfuraURL string
	InfuraKey string
	Origins   []string
	Admins    []string
	JWTSecret string
	Logging   bool
	DBFile    string
}

// Listen handles all the requests
func Listen(port int, conf Config) (*fiber.App, error) {
	// init Web server
	app := fiber.New()
	// app.Settings.DisableStartupMessage = true

	// Add graceful recovery in case of panics
	app.Use(middleware.Recover())

	// Pass a custom logger config
	if conf.Logging == true {
		app.Use(middleware.Logger(middleware.LoggerConfig{
			Format:     "${time} | ${status} ${method} ${path}\n",
			TimeFormat: "Mon, 02 Jan 2006 15:04:05 MST",
			TimeZone:   "Europe/Paris",
			Output:     os.Stdout,
		}))
	}

	// Add CORS settings
	if len(conf.Origins) == 0 {
		return nil, errors.New("You must provide a list of allowed HTTP Origins")
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins: conf.Origins,
	}))

	// JWT Middleware
	if len(conf.JWTSecret) == 0 {
		return nil, errors.New("You must provide the Eth private key for signing JWTs")
	}
	privKey, err := auth.ParseHexKey(conf.JWTSecret)
	if err != nil {
		return nil, err
	}
	jtwConf := jwtware.Config{
		SigningMethod: "ES256",
		SigningKey:    &privKey.PublicKey,
	}
	handlers.EthSignKey = conf.JWTSecret

	// Init Infura
	if len(conf.InfuraURL) > 0 && len(conf.InfuraKey) > 0 {
		err := resolver.InitInfura(conf.InfuraURL, conf.InfuraKey)
		if err != nil {
			return nil, err
		}
		handlers.ResolveENS = true
	}

	// Get a challenge for the user to sign, based on the given address
	app.Get("/login/:addr", handlers.Challenge)

	// Validate a login challenge for the given address
	app.Post("/login/:addr", handlers.Validate)

	// Refresh a token
	app.Get("/refresh", jwtware.New(jtwConf), handlers.Refresh)

	app.Listen(port)
	return app, nil
}
