package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// {
//     "status": "success",
//     "data": {
//         "ip": "1.1.1.1",
//         "ip_num": "16843009",
//         "country_code": "HK",
//         "country_name": "Hong Kong",
//         "city": "Hong Kong",
//         "isp": "Research Prefix for APNIC Labs",
//         "domain": "apnic.net",
//         "block": 1
//     }
// }

type IphunterResponse struct {
	Status string
	Data   IphunterData
}

type IphunterData struct {
	IP           string
	Ip_num       string
	Country_code string
	Country_name string
	City         string
	Isp          string
	Block        Block
}

func iphunter(iptocheck string) (ClientAddress, error) {

	// Create and Format new Request
	requestUrl := fmt.Sprintf("https://www.iphunter.info:8082/v1/ip/%s", iptocheck)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return ClientAddress{}, err
	}

	// Send the Request ot the API
	req.Header.Add("X-Key", os.Getenv("APIKEY_IPHUNTER"))
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
	var apiResponse IphunterResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return ClientAddress{}, err
	}

	if apiResponse.Status == "error" {
		return ClientAddress{}, errors.New("api has returned an error")
	}

	clientAddress := ClientAddress{
		IP:          apiResponse.Data.IP,
		CountryCode: apiResponse.Data.Country_code,
		CountryName: apiResponse.Data.Country_name,
		Asn:         0,
		Isp:         apiResponse.Data.Isp,
		Block:       apiResponse.Data.Block,
	}

	log.Printf("[IPHUNTER] - IP:%s Response:%s \n", iptocheck, body)
	return clientAddress, nil

}
