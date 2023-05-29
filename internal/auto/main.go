package auto

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/BasedDevelopment/eve/internal/config"
	"github.com/BasedDevelopment/eve/pkg/pki"
	"github.com/BasedDevelopment/eve/pkg/util"
	"github.com/rs/zerolog/log"
)

var (
	caPool  *x509.CertPool
	crtPair *tls.Certificate
)

// The TLS client that we will use to talk to auto
var TLSClient http.Client

func Init() {
	log.Info().Msg("Loading PKI")

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

	eveCrtPair, err := tls.LoadX509KeyPair(crtPath, keyPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load keypair")
	}

	crtPair = &eveCrtPair

	// Loading just the cert to get serial
	eveCrtBytes := util.ReadFile(crtPath)
	eveCrt := pki.ReadCrt(eveCrtBytes)
	serial := eveCrt.SerialNumber.String()
	eveCrtHash := pki.PemSum(eveCrtBytes)

	log.Info().
		Str("crt serial", serial).
		Str("crt hash", eveCrtHash).
		Msg("PKI loaded")
}

type Auto struct {
	Url    string
	Serial string
}

func (a *Auto) getTLSConfig() *tls.Config {
	return &tls.Config{
		RootCAs:      caPool,
		Certificates: []tls.Certificate{*crtPair},
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			// Make sure the serial number is correct
			if verifiedChains[0][0].SerialNumber.String() != a.Serial {
				return fmt.Errorf("serial number mismatch")
			}
			return nil
		},
	}
}

func (a *Auto) getHttpsClient() *http.Client {
	tlsConfig := a.getTLSConfig()

	TLSClient = http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	return &TLSClient
}

func (a *Auto) WSProxy(wsUrl *url.URL, w http.ResponseWriter, r *http.Request) {
	tlsConfig := a.getTLSConfig()

	proxy := httputil.NewSingleHostReverseProxy(wsUrl)
	proxy.Transport = &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	r.URL.Path = ""

	proxy.ServeHTTP(w, r)
}
