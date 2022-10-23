package main

import (
	"fmt"
	"strconv"

	"github.com/ip2location/ip2proxy-go/v3"
)

func CheckProxy(iptocheck string) (ClientAddress, error) {
	db, err := ip2proxy.OpenDB("./data/IP2PROXY-LITE-PX10.BIN")

	if err != nil {
		return ClientAddress{}, err
	}
	all, err := db.GetAll(iptocheck)

	if err != nil {
		fmt.Print(err)
		return ClientAddress{}, err
	}

	fmt.Printf("ModuleVersion: %s\n", ip2proxy.ModuleVersion())
	fmt.Printf("PackageVersion: %s\n", db.PackageVersion())
	fmt.Printf("DatabaseVersion: %s\n", db.DatabaseVersion())

	fmt.Printf("isProxy: %s\n", all["isProxy"])
	fmt.Printf("ProxyType: %s\n", all["ProxyType"])
	fmt.Printf("CountryShort: %s\n", all["CountryShort"])
	fmt.Printf("CountryLong: %s\n", all["CountryLong"])
	fmt.Printf("Region: %s\n", all["Region"])
	fmt.Printf("City: %s\n", all["City"])
	fmt.Printf("ISP: %s\n", all["ISP"])
	fmt.Printf("Domain: %s\n", all["Domain"])
	fmt.Printf("UsageType: %s\n", all["UsageType"])
	fmt.Printf("ASN: %s\n", all["ASN"])
	fmt.Printf("AS: %s\n", all["AS"])
	fmt.Printf("LastSeen: %s\n", all["LastSeen"])
	fmt.Printf("Threat: %s\n", all["Threat"])
	fmt.Printf("Provider: %s\n", all["Provider"])

	db.Close()

	asn, err := strconv.Atoi(all["ASN"])
	if err != nil {

	}

	proxy, err := strconv.Atoi(all["isProxy"])
	if err != nil {

	}

	clientAddress := ClientAddress{
		IP:          iptocheck,
		CountryCode: all["CountryShort"],
		CountryName: all["CountryLong"],
		Asn:         asn,
		Isp:         all["ISP"],
		Block:       Block(proxy),
	}

	return clientAddress, nil
}
