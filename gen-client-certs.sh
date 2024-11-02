#!/bin/bash

# generate agent key
openssl ecparam -genkey -name secp384r1 -out certs/agent/agent.key

# generate CSR
openssl req -new -key certs/agent/agent.key -out certs/agent/agent.csr

# generate signed certificate
openssl x509 -req -in certs/agent/agent.csr -CA certs/server/server.crt -CAkey certs/server/server.key -CAcreateserial -out certs/agent/agent.crt -days 365 -sha256
