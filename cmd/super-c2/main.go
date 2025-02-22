package main

import (
	"flag"
	"fmt"

	"github.com/joe-broder15/supertrooper/internal/server"
)

var banner string = `	
 _______ _______ ______ _______ ______ _______ ______ _______ _______ ______ _______ ______ 
|     __|   |   |   __ \    ___|   __ \_     _|   __ \       |       |   __ \    ___|   __ \
|__     |   |   |    __/    ___|      < |   | |      <   -   |   -   |    __/    ___|      <
|_______|_______|___|  |_______|___|__| |___| |___|__|_______|_______|___|  |_______|___|__|
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
