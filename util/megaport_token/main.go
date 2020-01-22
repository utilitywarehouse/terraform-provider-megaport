package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

const (
	usage = `usage: megaport_token [--reset]`
)

func main() {
	var (
		reset                          = false
		endpoint                       = api.EndpointStaging
		username, password, totpSecret string
	)
	if len(os.Args) > 2 {
		log.Fatalln(usage)
	}
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "--reset":
			reset = true
			log.Println("The token will be reset to retrieve a new one")
		default:
			log.Fatalln(usage)
		}
	}
	if v := os.Getenv("MEGAPORT_ENDPOINT"); v != "" {
		endpoint = v
	}
	log.Printf("Endpoint: %s\n", endpoint)
	if username = os.Getenv("MEGAPORT_USERNAME"); username == "" {
		log.Fatalln("MEGAPORT_USERNAME is empty")
	}
	if password = os.Getenv("MEGAPORT_PASSWORD"); password == "" {
		log.Fatalln("MEGAPORT_PASSWORD is empty")
	}
	totpSecret = os.Getenv("MEGAPORT_TOTP_SECRET")
	token, err := getToken(endpoint, username, password, totpSecret, reset)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("MEGAPORT_TOKEN=%s\n", token)
}

func getToken(endpoint, username, password, totpSecret string, reset bool) (string, error) {
	otp := ""
	if totpSecret != "" {
		v, err := totp.GenerateCode(totpSecret, time.Now())
		if err != nil {
			return "", err
		}
		otp = v
	}
	c := api.NewClient(endpoint)
	if reset {
		if err := c.Login(username, password, otp); err != nil {
			return "", err
		}
		if err := c.Logout(); err != nil {
			return "", err
		}
	}
	if err := c.Login(username, password, otp); err != nil {
		return "", err
	}
	return c.Token, nil
}
