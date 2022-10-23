package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func iphub(iptocheck string) (ClientAddress, error) {

	// Create and Format new Request
	requestUrl := fmt.Sprintf("http://v2.api.iphub.info/ip/%s", iptocheck)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return ClientAddress{}, err
	}

	// Send the Request ot the API
	req.Header.Add("X-Key", os.Getenv("APIKEY_IPHUB"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ClientAddress{}, err
	}

	// Read the Body Information from Request
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ClientAddress{}, err
	}

	// Check if the Api Request has returned any Error
	if strings.Contains(string(body), "\"error\"") {
		log.Printf("[IPHUB] - IP:%s Error:%s \n", iptocheck, body)
		// Try to Parse the Error Message
		var ipHubError IpHubError
		err := json.Unmarshal(body, &ipHubError)
		if err != nil {
			return ClientAddress{}, err
		}
		return ClientAddress{}, err
	}

	// Parse Body Response into Object
	var clientAddress ClientAddress
	err = json.Unmarshal(body, &clientAddress)
	if err != nil {
		return ClientAddress{}, err
	}

	log.Printf("[IPHUB] - IP:%s Response:%s \n", iptocheck, body)
	return clientAddress, nil

}
