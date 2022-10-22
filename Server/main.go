package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/helmet/v2"
	"github.com/joho/godotenv"
)

type Block int64

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
		// for key, value := range c.GetReqHeaders() {
		// 	log.Printf("%s=%s \n", key, value)
		// }
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

		// Create and Format new Request
		requestUrl := fmt.Sprintf("http://v2.api.iphub.info/ip/%s", realip)
		req, err := http.NewRequest("GET", requestUrl, nil)
		if err != nil {
			log.Printf("IP:%s Error:%s \n", realip, err)
			apiresponse.Error = err.Error()
			apiresponse.Success = false
			return c.JSON(apiresponse)
		}

		// Send the Request ot the API
		req.Header.Add("X-Key", os.Getenv("APIKEY"))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("IP:%s Error:%s \n", realip, err)
			apiresponse.Error = err.Error()
			apiresponse.Success = false
			return c.JSON(apiresponse)
		}

		// Read the Body Information from Request
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			apiresponse.Error = err.Error()
			apiresponse.Success = false
			return c.JSON(apiresponse)
		}

		// Check if the Api Request has returned any Error
		if strings.Contains(string(body), "\"error\"") {
			log.Printf("IP:%s Error:%s \n", realip, body)
			// Try to Parse the Error Message
			var ipHubError IpHubError
			err := json.Unmarshal(body, &ipHubError)
			if err != nil {
				apiresponse.Error = err.Error()
				apiresponse.Success = false
				return c.JSON(apiresponse)
			}
			apiresponse.Error = ipHubError.Error
			apiresponse.Success = false
			return c.JSON(apiresponse)

		}

		// Parse Body Response into Object
		var clientAddress ClientAddress
		err = json.Unmarshal(body, &clientAddress)
		if err != nil {
			apiresponse.Error = err.Error()
			apiresponse.Success = false
			return c.JSON(apiresponse)
		}

		// Send Response back to User
		log.Printf("IP:%s Response:%s \n", realip, body)
		apiresponse.Message = "success"
		apiresponse.Success = true
		apiresponse.Data = clientAddress
		return c.JSON(apiresponse)
	})
	// Request for Metrics
	app.Get("/metrics", monitor.New(monitor.Config{Title: "VPN-Check-Server Monitoring"}))

	app.Static("/", "./support.html")

	log.Fatal(app.Listen(os.Getenv("PORT")))
}
