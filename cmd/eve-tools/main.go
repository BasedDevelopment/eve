package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/BasedDevelopment/eve/internal/config"
	"github.com/BasedDevelopment/eve/pkg/pki"
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
			Msg("TLS Path does not exist, Creating.")
		if err := os.MkdirAll(tlsPath, 0700); err != nil {
			log.Fatal().
				Err(err).
				Str("path", tlsPath).
				Msg("Failed to create TLS Path")
		}
	}
}

func main() {
	if *makeCAKey {
		log.Info().Msg("Creating CA key")
		b := pki.GenKey()
		writeFile(caPrivPath, b)
		return
	}

	if *makeCA {
		log.Info().Msg("Creating CA certificate")
		caPrivBytes := readFile(caPrivPath)
		caPriv := pki.ReadKey(caPrivBytes)
		caCrt := pki.GenCA(caPriv)
		writeFile(caCrtPath, caCrt)
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
		writeFile(evePrivPath, b)
		return
	}

	if *makeCrt {
		log.Info().Msg("Creating eve certificate")
		// First we make a CSR for eve
		evePrivBytes := readFile(evePrivPath)
		evePriv := pki.ReadKey(evePrivBytes)
		eveCsrBytes := pki.GenCSR(evePriv, config.Config.Hostname)
		eveCsr := pki.ReadCSR(eveCsrBytes)

		// Then we take the CA's key and sign the CSR
		caPrivBytes := readFile(caPrivPath)
		caPriv := pki.ReadKey(caPrivBytes)
		caCrtBytes := readFile(caCrtPath)
		caCrt := pki.ReadCrt(caCrtBytes)
		crt := pki.SignCrt(caCrt, caPriv, eveCsr)

		writeFile(eveCrtPath, crt)

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
		csrBytes := readFile(*signCSR)
		csrSum := pki.PemSum(csrBytes)
		csr := pki.ReadCSR(csrBytes)
		log.Info().
			Str("SHA1", csrSum).
			Str("host", host).
			Msg("CSR checksum")

		// Fetch the CA's key and certificate
		caCrtBytes := readFile(caCrtPath)
		caCrt := pki.ReadCrt(caCrtBytes)
		caPrivBytes := readFile(caPrivPath)
		caPriv := pki.ReadKey(caPrivBytes)

		crtBytes := pki.SignCrt(caCrt, caPriv, csr)
		writeFile(host+".crt", crtBytes)

		crtSum := pki.PemSum(crtBytes)
		log.Info().
			Str("SHA1", crtSum).
			Str("host", host).
			Msg("Certificate signed")
		return
	}

	if *checkSum != "" {
		b := readFile(*checkSum)
		result := pki.PemSum(b)
		log.Info().
			Str("path", *checkSum).
			Str("SHA1", result).
			Msg("Checksum")
		return
	}

	fmt.Println("No action specified, checking PKI")
	checkPKI()
}

// Read a file, exit if err
func readFile(path string) []byte {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal().
			Str("path", path).
			Err(err).
			Msg("Failed to read file")
	}
	return b
}

// Write a file, exit if err
func writeFile(path string, data []byte) {
	if err := ioutil.WriteFile(path, data, 0600); err != nil {
		log.Fatal().
			Str("path", path).
			Err(err).
			Msg("Failed to write file")
	}
}

// Check if a file exists
func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
