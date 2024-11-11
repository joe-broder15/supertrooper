#!/bin/bash
openssl genpkey -algorithm RSA -out certs/agent/agent_private_key.pem
openssl req -new -x509 -key certs/agent/agent_private_key.pem -out certs/agent/agent_cert.pem -days 365
cp certs/agent/agent_cert.pem cmd/super-agent/embed/agent_cert.pem
cp certs/agent/agent_private_key.pem cmd/super-agent/embed/agent_private_key.pem