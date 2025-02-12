package server

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/joe-broder15/supertrooper/internal/messages"
)

func prettyPrintJSON(jsonData []byte) string {
	var obj interface{}
	if err := json.Unmarshal(jsonData, &obj); err != nil {
		log.Printf("server: error unmarshalling json for pretty print: %v", err)
		return string(jsonData)
	}
	pretty, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		log.Printf("server: error pretty printing json: %v", err)
		return string(jsonData)
	}
	return string(pretty)
}

func initTLSConfig(serverCertFile string, serverKeyFile string, caCertFile string) *tls.Config {
	// Load server certificate and private key
	cert, err := tls.LoadX509KeyPair(serverCertFile, serverKeyFile)
	if err != nil {
		log.Fatalf("server: error loading key pair: %v", err)
	}

	// Load CA certificate for verifying client's certificate
	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		log.Fatalf("server: error reading CA certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		log.Fatal("server: failed to append CA certificate")
	}

	// Set up mTLS configuration for the server.
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
		MinVersion:   tls.VersionTLS12,
	}
	return tlsConfig
}

// handleAgentConnection handles the connection from the agent and sends messages to the channel
func handleAgentConnection(tlsConn *tls.Conn) {
	defer tlsConn.Close()

	log.Println("server: handling agent connection from ip:", tlsConn.RemoteAddr())

	// create a decoder to read messages from the agent
	decoder := json.NewDecoder(tlsConn)

	for {
		// read the hello message from the agent
		var message messages.C2MessageBase
		err := decoder.Decode(&message)

		// if the agent connection is closed, return
		if err == io.EOF {
			log.Println("server: agent connection closed")
			return
		}
		if err != nil {
			log.Printf("server: error decoding message: %v", err)
			return
		}

		// pretty print the message
		prettyMessage := prettyPrintJSON(message.C2MessagePayload)
		log.Println(prettyMessage)

		// if the job was a hello, send a response to the agent
		if message.Type == messages.C2MessageTypeHello {
			printPayload := messages.BuildJobPayloadPrint("hello")
			response := messages.BuildJobDispatchMessage("server", "123", messages.JobTypePrint, printPayload)
			encoder := json.NewEncoder(tlsConn)
			encoder.Encode(response)
		}

	}
}

// startTlsListener starts the TLS listener and sends messages to the channel
func startTlsListener(tlsConfig *tls.Config) {
	// Create a TLS listener directly on a TCP socket
	listener, err := tls.Listen("tcp", ":443", tlsConfig)
	if err != nil {
		log.Fatalf("server: failed to listen: %v", err)
	}
	defer listener.Close()
	// Accept connections in a loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: error accepting connection: %v", err)
			continue
		}

		// handle the agent connection
		go handleAgentConnection(conn.(*tls.Conn))
	}
}

func Start(serverCertFile string, serverKeyFile string, caCertFile string, configFile string) {
	fmt.Println("starting server as mTLS server using raw TLS sockets")

	// initialize TLS config
	tlsConfig := initTLSConfig(serverCertFile, serverKeyFile, caCertFile)

	// create a channel to send messages to the agent
	// agentMessageChannel := make(chan messages.C2MessageBase, 512)

	// start the agent listener
	go startTlsListener(tlsConfig)

	// wait for the above function to finish
	select {}

	// // read messages from the agent and pretty print the json
	// for message := range agentMessageChannel {
	// 	prettyMessage, err := json.MarshalIndent(message, "", "    ")
	// 	if err != nil {
	// 		log.Printf("server: error pretty printing AgentHello message: %v", err)
	// 	} else {
	// 		fmt.Println(string(prettyMessage))
	// 	}
	// }
}
