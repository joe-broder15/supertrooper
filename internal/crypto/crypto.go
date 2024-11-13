package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

func Base64Encode(in []byte) string {
	return base64.StdEncoding.EncodeToString(in)
}

func Base64Decode(in string) ([]byte, error) {
	out, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GenerateRandomBytes generates a base64 encoded string of 256 random bytes
func GenerateRandomBytes() (string, error) {
	// make a buffer
	randomBytes := make([]byte, 256)

	// read from rand into the buffer
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return Base64Encode(randomBytes), nil
}

func RSASignNonce(keyPair tls.Certificate, nonceBase64 string) (string, error) {
	// determine whether the keys are rsa
	rsaPrivateKey, ok := keyPair.PrivateKey.(*rsa.PrivateKey)
	if !ok {
		return "", fmt.Errorf("private key is not of type rsa")
	}

	// Decode nonce from base64
	nonceBytes, err := Base64Decode(nonceBase64)
	if err != nil {
		return "", fmt.Errorf("failed to decode nonce: %v", err)
	}

	digest := sha256.Sum256(nonceBytes)

	// sign the agent nonce and generate a server nonce
	signature, err := rsa.SignPSS(rand.Reader, rsaPrivateKey, crypto.SHA256, digest[:], nil)
	if err != nil {
		return "", err
	}

	return Base64Encode(signature), nil
}

func RSAVerifySignature(pubKey []byte, sigBase64 string, nonceBase64 string) (bool, error) {
	// Decode the PEM block
	block, _ := pem.Decode(pubKey)
	if block == nil || block.Type != "CERTIFICATE" {
		return false, fmt.Errorf("failed to decode pem bytes")
	}

	// Parse the public key
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return false, fmt.Errorf("failed to parse DER-encoded public key: %v", err)
	}

	// Decode signature from base64
	sigBytes, err := Base64Decode(sigBase64)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %v", err)
	}

	// Decode nonce from base64
	nonceBytes, err := Base64Decode(nonceBase64)
	if err != nil {
		return false, fmt.Errorf("failed to decode nonce: %v", err)
	}

	// get digest of nonce
	digest := sha256.Sum256(nonceBytes)

	// verify the signature
	err = rsa.VerifyPSS(cert.PublicKey.(*rsa.PublicKey), crypto.SHA256, digest[:], sigBytes, nil)
	if err != nil {
		return false, nil
	} else {
		return true, nil
	}
}
