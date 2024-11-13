package server

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joe-broder15/supertrooper/internal/crypto"
	"github.com/joe-broder15/supertrooper/internal/messages"
)

type Server struct {
}

func Start(serverPubKeyFile string, serverPrivKeyFile string, agentPubKeyFile string) {

	// Load server certificate and private key
	serverKeyPair, err := tls.LoadX509KeyPair(serverPubKeyFile, serverPrivKeyFile)
	if err != nil {
		log.Fatalf("Failed to load server certificate: %v", err)
	}

	// Load agent public key
	agentPubKeyBytes, err := os.ReadFile(agentPubKeyFile)
	if err != nil {
		log.Fatalf("Failed to load CA certificate: %v", err)
	}
	agentCertPool := x509.NewCertPool()
	agentCertPool.AppendCertsFromPEM(agentPubKeyBytes)

	// Configure TLS with agent certificate verification
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverKeyPair},
		ClientCAs:    agentCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert, // Enforce agent certificate verification
	}

	// Set up the HTTPS server using the above TLS config
	server := &http.Server{
		Addr:      ":443",
		TLSConfig: tlsConfig,
	}

	// Define a simple handler
	http.HandleFunc("/challenge", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodPost {
			log.Printf("challenge request from %v", r.RemoteAddr)

			// read the request body into a slice
			requestBody, err := io.ReadAll(r.Body)
			if err != nil {
				log.Fatalln(err)
			}

			// unmarshal the body into a message
			var agentMessage messages.Message
			err = json.Unmarshal(requestBody, &agentMessage)
			if err != nil {
				return
			}

			// unmarshal the message body into an AgentChallengeRequestBodyu
			var agentMessageBody messages.AgentChallengeRequestBody
			if agentMessage.Type != messages.MessageTypeAgentChallengeRequest {
				return
			}
			err = json.Unmarshal([]byte(agentMessage.Body), &agentMessageBody)
			if err != nil {
				return
			}

			// sign the agent nonce
			signedAgentNonce, err := crypto.RSASignNonce(serverKeyPair, agentMessageBody.AgentNonce)
			if err != nil {
				log.Printf("failed to sign agent nonce: %v", err)
				return
			}

			// get a nonce to send to the agent
			serverNonce, err := crypto.GenerateRandomBytes()
			if err != nil {
				log.Printf("failed to generate server nonce: %v", err)
				return
			}

			// create response
			serverChallengeResponse, err := messages.NewServerChallengeResponse(signedAgentNonce, serverNonce)
			if err != nil {
				log.Printf("failed to generate ServerChallengeResponse message: %v", err)
				return
			}

			_, err = w.Write(serverChallengeResponse)
			if err != nil {
				log.Printf("failed to write response: %v", err)
				return
			}
		}

	})

	// Define a simple handler
	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("got connection")
	})

	// Define a simple handler
	http.HandleFunc("/beacon", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("got connection")
	})

	// Define a simple handler
	http.HandleFunc("/payload", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("got connection")
	})

	// Start the HTTPS server
	log.Printf("Starting HTTPS server on port 443")
	err = server.ListenAndServeTLS("", "") // Certificates are already configured
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
