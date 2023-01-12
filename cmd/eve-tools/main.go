package main

import (
	"flag"
	"os"
	"strings"

	"github.com/BasedDevelopment/eve/internal/config"
	"github.com/BasedDevelopment/eve/pkg/pki"
	"github.com/BasedDevelopment/eve/pkg/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	version = "0.0.1"
)

var (
	configPath = flag.String("config-path", "/etc/eve/config.toml", "Path to TLS key and certificate")
	makeCAKey  = flag.Bool("make-cakey", false, "Create key for certificate authority")
	makeCA     = flag.Bool("make-ca", false, "Create certificate authority")
	makeCrtKey = flag.Bool("make-crtkey", false, "Create key for certificate")
	makeCrt    = flag.Bool("make-crt", false, "Make a certificate")
	signCSR    = flag.String("sign-csr", "", "Sign a CSR, put path to CSR here")
	checkSum   = flag.String("checksum", "", "Check the checksum of a pem encoded file")
)

var (
	caPrivPath  string
	caCrtPath   string
	evePrivPath string
	eveCSRPath  string
	eveCrtPath  string
	csrPath     string
)

func init() {
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("eve-tools, tools to manage the CA on an eve instance")

	if err := config.Load(*configPath); err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	tlsPath := config.Config.TLSPath

	// Ensure TLS path has a slash at the end
	if tlsPath[len(tlsPath)-1:] != "/" {
		tlsPath += "/"
	}

	// Initialize path vars
	caPrivPath = tlsPath + "ca.key"
	caCrtPath = tlsPath + "ca.crt"

	evePrivPath = tlsPath + config.Config.Hostname + ".key"
	eveCSRPath = tlsPath + config.Config.Hostname + ".csr"
	eveCrtPath = tlsPath + config.Config.Hostname + ".crt"

	csrPath = *signCSR

	// Ensure TLS Path exists
	if _, err := os.Stat(tlsPath); os.IsNotExist(err) {
		log.Info().
			Str("path", tlsPath).
			Msg("TLS Path does not exist, creating.")
		if err := os.MkdirAll(tlsPath, 0700); err != nil {
			log.Fatal().
				Err(err).
				Str("path", tlsPath).
				Msg("Failed to create TLS path")
		}
	}
}

func main() {
	if *makeCAKey {
		log.Info().Msg("Creating CA key")
		b := pki.GenKey()
		util.WriteFile(caPrivPath, b)
		log.Info().
			Str("path", caPrivPath).
			Msg("CA key written")
		return
	}

	if *makeCA {
		log.Info().Msg("Creating CA certificate")
		caPrivBytes := util.ReadFile(caPrivPath)
		caPriv := pki.ReadKey(caPrivBytes)
		caCrt := pki.GenCA(caPriv)
		util.WriteFile(caCrtPath, caCrt)
		sum := pki.PemSum(caCrt)
		log.Info().
			Str("path", caCrtPath).
			Str("SHA1", sum).
			Msg("Wrote CA certificate")
		return
	}

	if *makeCrtKey {
		log.Info().Msg("Creating eve key")
		b := pki.GenKey()
		util.WriteFile(evePrivPath, b)
		log.Info().
			Str("path", evePrivPath).
			Msg("eve key written")
		return
	}

	if *makeCrt {
		log.Info().Msg("Creating eve certificate")
		// First we make a CSR for eve
		evePrivBytes := util.ReadFile(evePrivPath)
		evePriv := pki.ReadKey(evePrivBytes)
		eveCsrBytes := pki.GenCSR(evePriv, config.Config.Hostname)
		eveCsr := pki.ReadCSR(eveCsrBytes)

		// Then we take the CA's key and sign the CSR
		caPrivBytes := util.ReadFile(caPrivPath)
		caPriv := pki.ReadKey(caPrivBytes)
		caCrtBytes := util.ReadFile(caCrtPath)
		caCrt := pki.ReadCrt(caCrtBytes)
		crt := pki.SignCrt(caCrt, caPriv, eveCsr)

		util.WriteFile(eveCrtPath, crt)

		sum := pki.PemSum(crt)
		log.Info().
			Str("path", eveCrtPath).
			Str("SHA1", sum).
			Msg("Wrote EVE certificate")

		return
	}

	if *signCSR != "" {
		host := strings.Split(*signCSR, ".csr")[0]
		if host == "" {
			log.Fatal().Msg("Failed to parse hostname from CSR")
		}

		// Read the CSR
		csrBytes := util.ReadFile(*signCSR)
		csrSum := pki.PemSum(csrBytes)
		csr := pki.ReadCSR(csrBytes)
		log.Info().
			Str("SHA1", csrSum).
			Str("host", host).
			Msg("CSR checksum")

		// Fetch the CA's key and certificate
		caCrtBytes := util.ReadFile(caCrtPath)
		caCrt := pki.ReadCrt(caCrtBytes)
		caPrivBytes := util.ReadFile(caPrivPath)
		caPriv := pki.ReadKey(caPrivBytes)

		crtBytes := pki.SignCrt(caCrt, caPriv, csr)
		util.WriteFile(host+".crt", crtBytes)

		crtSum := pki.PemSum(crtBytes)
		log.Info().
			Str("SHA1", crtSum).
			Str("host", host).
			Msg("Certificate signed")
		return
	}

	if *checkSum != "" {
		b := util.ReadFile(*checkSum)
		result := pki.PemSum(b)
		log.Info().
			Str("path", *checkSum).
			Str("SHA1", result).
			Msg("Checksum")
		return
	}

	log.Info().Msg("No action specified, checking PKI")
	checkPKI()
}
