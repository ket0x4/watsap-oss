package secure

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"watsap/utils/config"
	"watsap/utils/messages"
	"watsap/utils/telegram"
)

func SSLPinning() {
	// Load the certificate
	cert, err := os.ReadFile("telegram-org-chain.pem")
	if err != nil {
		config.Logger("Failed to read certificate file: "+err.Error(), "error")
		telegram.TgSendMsg(messages.CertErrMessage())
		os.Exit(1)
	}

	// Create a new certificate pool and add the loaded certificate
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(cert); !ok {
		config.Logger("Failed to append certificate to pool: invalid certificate", "error")
		telegram.TgSendMsg(messages.CertErrMessage())
		os.Exit(1)
	}

	// Create a custom TLS config with our certificate pool
	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	// Create a HTTPS client with the custom TLS config
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// Make a request to the API
	resp, err := client.Get("https://api.telegram.org")
	if err != nil {
		config.Logger("Failed to make request to API: "+err.Error(), "error")
		telegram.TgSendMsg(messages.CertErrMessage())
	}
	defer resp.Body.Close()

	config.Logger("SSL Pinning successfull", "info")
	telegram.TgSendMsg(messages.CertErrMessage())
}
