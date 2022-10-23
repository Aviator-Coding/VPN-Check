package main

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/helmet/v2"
	"github.com/joho/godotenv"
)

type Block int

const (
	Residential Block = iota
	Proxy
	Warning
)

type ClientAddress struct {
	IP          string
	CountryCode string
	CountryName string
	Asn         int
	Isp         string
	Block       Block
}

// Response Send to the Client
type ApiResponse struct {
	Message string
	Date    time.Time
	Success bool
	Error   string
	Data    ClientAddress
}

type IpHubError struct {
	Error string
}

func (s Block) String() string {
	switch s {
	case Residential:
		return "Residential"
	case Proxy:
		return "Proxy"
	case Warning:
		return "Warning"
	}
	return "unknown"
}

var toggle int = 1

func main() {

	envfile := ".env"
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) > 0 {
		envfile = argsWithoutProg[0]
	}
	// Load Enviroment Variabels
	err := godotenv.Load(envfile)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()
	app.Use(helmet.New())

	// Logging for the Requests
	app.Use(
		logger.New(),
	)

	// Favicon Icon
	app.Use(favicon.New(favicon.Config{
		File: "./favicon.ico",
	}))

	// Read the API Limits From the Enviroment Variable
	limit, err := strconv.Atoi(os.Getenv("LIMIT"))
	if err != nil {
		limit = 5
	}

	// Limit the Amount of Connections
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:                    limit,
		Expiration:             30 * time.Second,
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
		LimiterMiddleware:      limiter.SlidingWindow{},
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
	}))

	// Request for the API
	app.Get("/api", func(c *fiber.Ctx) error {
		apiresponse := ApiResponse{
			Date: time.Now(),
		}

		realip := c.Get("X-Real-Ip")
		if realip == "" {
			realip = c.IP()
		}
		// Check if the Sender send the Right Auth with the Request
		if c.Get("X-KEY") != os.Getenv("SECURITYKEY") {
			log.Printf("IP:%s Key:%s \n", realip, c.Get("X-KEY"))
			apiresponse.Error = "has no rights"
			apiresponse.Success = false
			return c.JSON(apiresponse)
		}

		clientAddress, err := getData(realip)
		if err != nil {
			apiresponse.Error = err.Error()
			apiresponse.Success = false
			return c.JSON(apiresponse)
		}

		apiresponse.Message = "success"
		apiresponse.Success = true
		apiresponse.Data = clientAddress
		return c.JSON(apiresponse)
	})
	app.Post("/api", func(c *fiber.Ctx) error {
		apiresponse := ApiResponse{
			Date: time.Now(),
		}

		payload := struct {
			IP string `json:"ip"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			apiresponse.Error = "has no rights"
			apiresponse.Success = false
			return c.JSON(apiresponse)
		}

		realip := payload.IP

		// Check if the Sender send the Right Auth with the Request
		if c.Get("X-KEY") != os.Getenv("SECURITYKEY") {
			log.Printf("IP:%s Key:%s \n", realip, c.Get("X-KEY"))
			apiresponse.Error = "has no rights"
			apiresponse.Success = false
			return c.JSON(apiresponse)
		}

		clientAddress, err := getData(realip)
		if err != nil {
			apiresponse.Error = err.Error()
			apiresponse.Success = false
			return c.JSON(apiresponse)
		}

		apiresponse.Message = "success"
		apiresponse.Success = true
		apiresponse.Data = clientAddress
		return c.JSON(apiresponse)
	})
	// Request for Metrics
	app.Get("/metrics", monitor.New(monitor.Config{Title: "VPN-Check-Server Monitoring"}))

	app.Get("/test", func(c *fiber.Ctx) error {
		apiresponse := ApiResponse{
			Date: time.Now(),
		}

		payload := struct {
			IP string `json:"ip"`
		}{}

		if err := c.BodyParser(&payload); err != nil {
			apiresponse.Error = "has no rights"
			apiresponse.Success = false
			return c.JSON(apiresponse)
		}

		realip := payload.IP

		clientAddress, _ := CheckProxy(realip)

		apiresponse.Message = "success"
		apiresponse.Success = true
		apiresponse.Data = clientAddress
		return c.JSON(apiresponse)
	})

	app.Static("/", "./support.html")

	log.Fatal(app.Listen(os.Getenv("PORT")))
}

func getData(iptocheck string) (ClientAddress, error) {
	switch toggle {
	case 1:
		toggle++
		return iphub(iptocheck)
	case 2:
		toggle = 1
		return getipintel(iptocheck)
		// case 3:
		// 	toggle = 1
		// 	return iphunter(iptocheck)
	}
	return ClientAddress{}, errors.New("no handler found")
}
