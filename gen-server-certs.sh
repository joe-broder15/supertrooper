#!/bin/bash
openssl ecparam -genkey -name secp384r1 -out certs/server/server.key
openssl req -new -x509 -sha256 -key certs/server/server.key -out certs/server/server.crt 