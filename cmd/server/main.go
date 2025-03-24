package main

import (
	"flag"
	"fmt"

	"github.com/joe-broder15/supertrooper/internal/server"
)

var banner string = `	
======================
SUPERTROOPER C2 SERVER
======================
`

func main() {

	// set up command line args
	serverCertPtr := flag.String("server-cert", "certs/server/server_cert.pem", "path  server cert file")
	serverKeyPtr := flag.String("server-key", "certs/server/server_private_key.pem", "path server private key file")
	caCertPtr := flag.String("ca-cert", "certs/ca/ca_cert.pem", "path ca cert file")
	configPtr := flag.String("config", "config.toml", "path to the config file")
	flag.Parse()

	fmt.Println(banner)

	server.Start(*serverCertPtr, *serverKeyPtr, *caCertPtr, *configPtr)
}
