package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type ProxyCheckAPIResponse struct {
	status      string
	CountryCode string
	CountryName string
	Asn         int
	Isp         string
	Block       Block
}

// THis still needs Work
func proxycheckio(iptocheck string) (ClientAddress, error) {

	var results []map[string]interface{}
	// Create and Format new Request
	requestUrl := fmt.Sprintf("http://proxycheck.io/v2/%s?key=%s&vpn=3&asn=1", iptocheck, os.Getenv("APIKEY_PROXYCHECKIO"))
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return ClientAddress{}, err
	}

	// Send the Request ot the API
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

	// Parse Body Response into Object
	var clientAddress ClientAddress
	err = json.Unmarshal(body, &results)
	if err != nil {
		return ClientAddress{}, err
	}

	log.Printf("[IPHUB] - IP:%s Response:%s \n", iptocheck, body)
	return clientAddress, nil

}
