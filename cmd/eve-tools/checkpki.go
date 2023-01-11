package main

import (
	"github.com/BasedDevelopment/eve/internal/config"
	"github.com/BasedDevelopment/eve/pkg/pki"
	"github.com/rs/zerolog/log"
)

// Ensure that CA key, CA cert, eve key, eve cert are all present, create them if not. Log the checksum of all the certs.
func checkPKI() {
	if !fileExists(caPrivPath) {
		log.Info().
			Str("path", caPrivPath).
			Msg("CA key not found, creating a new one")
		caPriv := pki.GenKey()
		writeFile(caPrivPath, caPriv)
	}
	caPrivBytes := readFile(caPrivPath)
	caPriv := pki.ReadKey(caPrivBytes)

	if !fileExists(caCrtPath) {
		log.Info().
			Str("path", caCrtPath).
			Msg("CA cert not found, creating a new one")
		caCrt := pki.GenCA(caPriv)
		writeFile(caCrtPath, caCrt)
	}
	caCrtBytes := readFile(caCrtPath)
	caCrt := pki.ReadCrt(caCrtBytes)

	caCrtChecksum := pki.PemSum(caCrtBytes)
	log.Info().
		Str("path", caCrtPath).
		Str("checksum", caCrtChecksum).
		Msg("CA cert")

	if !fileExists(evePrivPath) {
		log.Info().
			Str("path", evePrivPath).
			Msg("Eve key not found, creating a new one")
		evePriv := pki.GenKey()
		writeFile(evePrivPath, evePriv)
	}
	evePrivBytes := readFile(evePrivPath)
	evePriv := pki.ReadKey(evePrivBytes)

	if !fileExists(eveCrtPath) {
		log.Info().
			Str("path", eveCrtPath).
			Msg("Eve cert not found, creating a new one")
		eveCSRBytes := pki.GenCSR(evePriv, config.Config.Hostname)
		writeFile(eveCSRPath, eveCSRBytes)
		eveCSR := pki.ReadCSR(eveCSRBytes)
		eveCrt := pki.SignCrt(caCrt, caPriv, eveCSR)
		writeFile(eveCrtPath, eveCrt)
	}
	eveCrt := readFile(eveCrtPath)
	eveCrtSum := pki.PemSum(eveCrt)
	log.Info().
		Str("path", eveCrtPath).
		Str("checksum", eveCrtSum).
		Msg("Eve cert")
}
