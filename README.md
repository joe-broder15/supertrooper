# supertrooper
## THIS PROJECT IS FOR EDUCATIONAL PURPOSES ONLY. DO NOT INSTALL THIS SOFTWARE ON ANY COMPUTERS THAT YOU DO NOT OWN OR DOES NOT BELONG TO A CONSENTING PARTY. MISUSE OF THIS SOFTWARE COULD SUBJECT YOU TO LEGAL PENALTIES. THE DEVELOPER OF THIS SOFTWARE DOES NOT TAKE ANY RESPONSIBILITY FOR MISUSE OF THIS CODE.
botnet written in go and inspired by emotet/kraken. 

## authentication
In order for an agent to connect to the c2 server, a challenge based handshake must be completed in order for the server to issue the client a token. The reason we do this is that supertrooper does NOT use signed certificates for TLS, so additional authentication must be implemented at the application level.

The handshake does the following steps
- client sends nonce to the server
- server responds with signed client nonce and a server nonce
- client responds with signed server nonce
- server issues client a token

## current priorities
- verification of self signed certs
- password / token beneath tls layer for aditional protection
- success sending json between agent and server
- develop messages passed between agents and servers
- develop middleware for agents to process messages
- develop tools for agents
- develop shell to send messages to agents from server