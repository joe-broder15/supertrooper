#!/bin/bash
openssl genpkey -algorithm RSA -out certs/server/server_private_key.pem
openssl req -new -x509 -key certs/server/server_private_key.pem -out certs/server/server_cert.pem -days 365
cp certs/server/server_cert.pem cmd/super-agent/embed/server_cert.pem