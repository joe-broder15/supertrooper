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
	agentCertPtr := flag.String("agent-cert", "certs/agent/agent_cert.pem", "path server private key file")
	flag.Parse()

	fmt.Println(banner)

	server.Start(*serverCertPtr, *serverKeyPtr, *agentCertPtr)
}
