package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
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

	// GET /
	app.Get("/", func(c *fiber.Ctx) error {
		apiresponse := ApiResponse{
			Date: time.Now(),
		}

		// Check if the Sender send the Right Auth with the Request
		if c.Get("X-KEY") != os.Getenv("SECURITYKEY") {
			log.Printf("IP:%s Key:%s \n", c.IP(), c.Get("X-KEY"))
			apiresponse.Error = "has no rights"
			apiresponse.Success = false
			return c.JSON(apiresponse)
		}

		// Create and Format new Request
		requestUrl := fmt.Sprintf("http://v2.api.iphub.info/ip/%s", c.IP())
		req, err := http.NewRequest("GET", requestUrl, nil)
		if err != nil {
			log.Printf("IP:%s Error:%s \n", c.IP(), err)
			apiresponse.Error = err.Error()
			apiresponse.Success = false
			return c.JSON(apiresponse)
		}

		// Send the Request ot the API
		req.Header.Add("X-Key", os.Getenv("APIKEY"))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("IP:%s Error:%s \n", c.IP(), err)
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
			log.Printf("IP:%s Error:%s \n", c.IP(), body)
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
		log.Printf("IP:%s Response:%s \n", c.IP(), body)
		apiresponse.Message = "success"
		apiresponse.Success = true
		return c.JSON(clientAddress)
	})

	log.Fatal(app.Listen(os.Getenv("PORT")))
}
