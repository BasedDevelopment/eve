package auto

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"

	"github.com/BasedDevelopment/eve/internal/config"
	"github.com/BasedDevelopment/eve/pkg/pki"
	"github.com/BasedDevelopment/eve/pkg/util"
	"github.com/rs/zerolog/log"
)

var (
	caPool *x509.CertPool
	cert   *tls.Certificate
)

// The TLS client that we will use to talk to auto
var TLSClient http.Client

func Init() {
	log.Info().Msg("Checking for PKI")

	tlsPath := config.Config.TLSPath

	if tlsPath[len(tlsPath)-1:] != "/" {
		tlsPath += "/"
	}

	// Define the paths to the PKI files
	crtPath := tlsPath + config.Config.Hostname + ".crt"
	keyPath := tlsPath + config.Config.Hostname + ".key"
	caPath := tlsPath + "ca.crt"

	// Ensure these paths exists
	if !util.FileExists(crtPath) {
		log.Fatal().Msg("Certificate not found")
	}

	if !util.FileExists(keyPath) {
		log.Fatal().Msg("Key not found")
	}

	if !util.FileExists(caPath) {
		log.Fatal().Msg("CA not found")
	}

	caBytes := util.ReadFile(caPath)
	ca := pki.ReadCrt(caBytes)

	caPool = x509.NewCertPool()
	caPool.AddCert(ca)

	eveCert, err := tls.LoadX509KeyPair(crtPath, keyPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load keypair")
	}

	cert = &eveCert

	log.Info().Msg("PKI ready")
}

type Auto struct {
	Url    string
	Serial string
}

func (a *Auto) getClient() *http.Client {
	tlsConfig := &tls.Config{
		RootCAs:      caPool,
		Certificates: []tls.Certificate{*cert},
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			// Make sure the serial number is correct
			if verifiedChains[0][0].SerialNumber.String() != a.Serial {
				return fmt.Errorf("serial number mismatch")
			}
			return nil
		},
	}

	TLSClient = http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	return &TLSClient
}
