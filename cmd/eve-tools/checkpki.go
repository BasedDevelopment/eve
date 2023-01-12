package main

import (
	"github.com/BasedDevelopment/eve/internal/config"
	"github.com/BasedDevelopment/eve/pkg/pki"
	"github.com/BasedDevelopment/eve/pkg/util"
	"github.com/rs/zerolog/log"
)

// Ensure that CA key, CA cert, eve key, eve cert are all present, create them if not. Log the checksum of all the certs.
func checkPKI() {
	if !util.FileExists(caPrivPath) {
		log.Info().
			Str("path", caPrivPath).
			Msg("CA key not found, creating a new one")
		caPriv := pki.GenKey()
		util.WriteFile(caPrivPath, caPriv)
	}
	caPrivBytes := util.ReadFile(caPrivPath)
	caPriv := pki.ReadKey(caPrivBytes)

	if !util.FileExists(caCrtPath) {
		log.Info().
			Str("path", caCrtPath).
			Msg("CA cert not found, creating a new one")
		caCrt := pki.GenCA(caPriv)
		util.WriteFile(caCrtPath, caCrt)
	}
	caCrtBytes := util.ReadFile(caCrtPath)
	caCrt := pki.ReadCrt(caCrtBytes)

	caCrtChecksum := pki.PemSum(caCrtBytes)
	log.Info().
		Str("path", caCrtPath).
		Str("SHA1", caCrtChecksum).
		Msg("CA cert")

	if !util.FileExists(evePrivPath) {
		log.Info().
			Str("path", evePrivPath).
			Msg("Eve key not found, creating a new one")
		evePriv := pki.GenKey()
		util.WriteFile(evePrivPath, evePriv)
	}
	evePrivBytes := util.ReadFile(evePrivPath)
	evePriv := pki.ReadKey(evePrivBytes)

	if !util.FileExists(eveCrtPath) {
		log.Info().
			Str("path", eveCrtPath).
			Msg("Eve cert not found, creating a new one")
		eveCSRBytes := pki.GenCSR(evePriv, config.Config.Hostname)
		util.WriteFile(eveCSRPath, eveCSRBytes)
		eveCSR := pki.ReadCSR(eveCSRBytes)
		eveCrt := pki.SignCrt(caCrt, caPriv, eveCSR)
		util.WriteFile(eveCrtPath, eveCrt)
	}
	eveCrt := util.ReadFile(eveCrtPath)
	eveCrtSum := pki.PemSum(eveCrt)
	log.Info().
		Str("path", eveCrtPath).
		Str("SHA1", eveCrtSum).
		Msg("Eve cert")
}
