package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/utilitywarehouse/terraform-provider-megaport/megaport/api"
)

const (
	usage = `usage: megaport_token [--reset]`
)

func main() {
	if len(os.Args) > 2 {
		log.Fatalln(usage)
	}
	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "--reset":
			token := os.Getenv("MEGAPORT_TOKEN")
			if token == "" {
				log.Fatal("To reset the token, please export MEGAPORT_TOKEN with your current token")
			}
			c := api.NewClient(api.EndpointStaging)
			c.Token = token
			if err := c.Logout(); err != nil {
				log.Fatal(err)
			}
			log.Print("Current token has been reset. Please login again to fetch a new one.")
		default:
			log.Fatalln(usage)
		}
	}
	scanner := bufio.NewScanner(os.Stdin)
	var username, password, otp string
	fmt.Printf("username: ")
	scanner.Scan()
	username = scanner.Text()
	fmt.Printf("password: ")
	scanner.Scan()
	password = scanner.Text()
	fmt.Printf("otp (leave empty if disabled): ")
	scanner.Scan()
	otp = scanner.Text()
	c := api.NewClient(api.EndpointStaging)
	if err := c.Login(username, password, otp); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("MEGAPORT_TOKEN=%s\n", c.Token)
}
