// embed cert.pem
package secure

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"
)

func SSLPinning() {
	// Load the certificate
	cert, err := os.ReadFile("cert.pem")
	//cert, err := config.CERT_PATH, nil
	if err != nil {
		log.Fatalf("Failed to read certificate file: %s", err.Error())
		Imha()
	}

	// Create a new certificate pool and add the loaded certificate
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(cert); !ok {
		log.Fatalf("Failed to append certificate to pool: invalid certificate %s", err.Error())
		Imha()
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
		log.Printf("Failed to make request to API: %s", err.Error())
		Imha()
	}
	defer resp.Body.Close()

	log.Println("SSL Pinning successfull")
}
