package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

// {
// 	"status": "success",
// 	"result": "1",
// 	"queryIP": "1.1.1.0",
// 	"queryFlags": null,
// 	"queryOFlags": "c",
// 	"queryFormat": "json",
// 	"contact": "info@aviator-coding.de",
// 	"Country": "AU"
//   }

type GetipintelApiResponse struct {
	Status      string
	Result      string
	QueryIP     string
	QueryFlags  string
	QueryOFlags string
	QueryFormat string
	Contact     string
	Country     string
}

func getipintel(iptocheck string) (ClientAddress, error) {

	// Create and Format new Request
	requestUrl := fmt.Sprintf("http://check.getipintel.net/check.php?ip=%s&contact=%s&format=json&oflags=c", iptocheck, os.Getenv("CONTACT_EMAIL"))
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
	var apiResponse GetipintelApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return ClientAddress{}, err
	}

	if apiResponse.Status == "error" {
		return ClientAddress{}, errors.New("api has returned an error")
	}

	result, err := strconv.Atoi(apiResponse.Result)
	if err != nil {
		return ClientAddress{}, errors.New("result is numeric")
	}

	clientAddress := ClientAddress{
		IP:          apiResponse.QueryIP,
		CountryCode: apiResponse.Country,
		CountryName: apiResponse.Country,
		Asn:         0,
		Isp:         "",
		Block:       Block(result),
	}

	log.Printf("[IPHUNTER] - IP:%s Response:%s \n", iptocheck, body)
	return clientAddress, nil

}
