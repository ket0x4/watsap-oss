package secure

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"embed"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"net/http"
	"os"
	"strings"
	"watsap/utils/config"
	"watsap/utils/logger"
)

//go:embed cert.pem
var certpem embed.FS

// Live Public Key (SPKI) SHA-256 fingerprint for api.telegram.org
const telegramLiveSPKI = "020c82993cac14e23a690092c9027e4085e99c69f4bfb9fe0fe9afea3580b507"

func SSLPinning() {
	var certBytes []byte
	var err error

	// Try reading from file first if path is specified, fallback to embedded certificate
	if config.CertPath != "" {
		certBytes, err = os.ReadFile(config.CertPath)
	}

	if err != nil || len(certBytes) == 0 {
		certBytes, _ = certpem.ReadFile("cert.pem")
	}

	var fileSPKIHash []byte
	// Extract the SPKI hash from the certificate file if loaded
	if len(certBytes) > 0 {
		block, _ := pem.Decode(certBytes)
		if block != nil && block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err == nil {
				hasher := sha256.New()
				hasher.Write(cert.RawSubjectPublicKeyInfo)
				fileSPKIHash = hasher.Sum(nil)
			}
		}
	}

	liveSPKIBytes, err := hex.DecodeString(telegramLiveSPKI)
	if err != nil {
		logger.Error("SSLPinning", "Failed to decode hardcoded Telegram SPKI: %s", err.Error())
		Imha()
		return
	}

	// Custom TLS Config with VerifyConnection callback (SPKI Pinning)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Custom SPKI Pinning via VerifyConnection handles server verification
		VerifyConnection: func(cs tls.ConnectionState) error {
			if len(cs.PeerCertificates) == 0 {
				return errors.New("server presented no certificates")
			}
			serverLeafCert := cs.PeerCertificates[0]

			hasher := sha256.New()
			hasher.Write(serverLeafCert.RawSubjectPublicKeyInfo)
			actualSPKIHash := hasher.Sum(nil)

			// Match against hardcoded live hash or the loaded file's hash
			if bytes.Equal(actualSPKIHash, liveSPKIBytes) {
				return nil
			}

			if fileSPKIHash != nil && bytes.Equal(actualSPKIHash, fileSPKIHash) {
				return nil
			}

			return errors.New("SSL Pinning failed: Public Key (SPKI) mismatch")
		},
	}

	// Apply tlsConfig globally to DefaultTransport so all http.Client requests are verified
	if transport, ok := http.DefaultTransport.(*http.Transport); ok {
		transport.TLSClientConfig = tlsConfig
	}

	// Create a HTTPS client with the custom TLS config
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// Make a request to the API to verify
	resp, err := client.Get("https://api.telegram.org")
	if err != nil {
		logger.Warn("SSLPinning", "Failed to make request to API: %s", err.Error())
		if strings.Contains(err.Error(), "SSL Pinning failed") {
			logger.Error("SSLPinning", "Certificate verification failed! Self-destructing...")
			Imha()
		} else {
			logger.Warn("SSLPinning", "Network connection unavailable/failed. Continuing execution.")
		}
		return
	}
	defer resp.Body.Close()

	logger.Info("SSLPinning", "SSL Pinning successful (SPKI Pinning verified)")
}
