#!/bin/bash
export MSYS_NO_PATHCONV=1

# Create directories for storing certificates
mkdir -p certs/ca certs/server certs/agent
mkdir -p cmd/super-agent/embed

echo "Removing existing certificates..."
rm -rf certs/ca/*
rm -rf certs/server/*
rm -rf certs/agent/*
rm -rf cmd/super-agent/embed/*


###########################################################
# 1. Generate Root CA keypair and self-signed certificate  #
###########################################################
echo "Generating Root CA keypair and certificate..."
openssl genpkey -algorithm RSA -out certs/ca/ca_private_key.pem
openssl req -new -x509 -key ./certs/ca/ca_private_key.pem -out ./certs/ca/ca_cert.pem -days 365 -subj "/C=US/ST=CA/L=SanFrancisco/O=MyOrg/OU=CA/CN=localhost"

###########################################################
# 2. Generate Server keypair and sign with Root CA        #
###########################################################
echo "Generating Server keypair and certificate..."
openssl genpkey -algorithm RSA -out certs/server/server_private_key.pem
openssl req -new -key ./certs/server/server_private_key.pem -out ./certs/server/server.csr -subj "/C=US/ST=CA/L=SanFrancisco/O=MyOrg/OU=Server/CN=localhost"

# required for git bash weirdness
echo "subjectAltName=DNS:localhost" > "tmp.txt"
openssl x509 -req -in ./certs/server/server.csr -CA ./certs/ca/ca_cert.pem -CAkey ./certs/ca/ca_private_key.pem -CAcreateserial -out ./certs/server/server_cert.pem -days 365 -extfile "tmp.txt"
rm "tmp.txt"

###########################################################
# 3. Generate Client keypair and sign with Root CA        #
###########################################################
echo "Generating Client keypair and certificate..."
openssl genpkey -algorithm RSA -out certs/agent/agent_private_key.pem
openssl req -new -key ./certs/agent/agent_private_key.pem -out ./certs/agent/agent.csr -subj "/C=US/ST=CA/L=SanFrancisco/O=MyOrg/OU=Client/CN=MyClient"
openssl x509 -req -in ./certs/agent/agent.csr -CA ./certs/ca/ca_cert.pem -CAkey ./certs/ca/ca_private_key.pem -CAcreateserial -out ./certs/agent/agent_cert.pem -days 365

# Copy certificates to embed directory (if required by super-agent)
cp ./certs/server/server_cert.pem cmd/super-agent/embed/server_cert.pem
cp ./certs/agent/agent_cert.pem cmd/super-agent/embed/agent_cert.pem
cp ./certs/agent/agent_private_key.pem cmd/super-agent/embed/agent_private_key.pem
cp ./certs/ca/ca_cert.pem cmd/super-agent/embed/ca_cert.pem

echo "Cleaning up temporary files..."
rm -f ./certs/server/server.csr ./certs/agent/agent.csr ./certs/ca/ca_cert.srl

echo "Certificate generation for mutual TLS complete." 