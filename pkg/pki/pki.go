package pki

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"math/big"
	mrand "math/rand"
	"time"

	"github.com/rs/zerolog/log"
)

var commonName string = "eve ca"

// TODO: Make this proper
var serialNumber *big.Int = big.NewInt((time.Now().Unix() + int64(mrand.Intn(4096))))
var today time.Time = time.Now()
var CaExpirey time.Time = time.Now().AddDate(10, 0, 0)
var CertExpirey time.Time = time.Now().AddDate(3, 0, 0)

// Take a PEM encoded SEC1,ASN1 DER private key and return the *ed25519.PrivateKey object
func ReadKey(keyPEMBytes []byte) ed25519.PrivateKey {
	pemBlock, _ := pem.Decode(keyPEMBytes)
	priv, err := x509.ParsePKCS8PrivateKey(pemBlock.Bytes)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse private key")
	}
	privKey := priv.(ed25519.PrivateKey)
	return privKey
}

// Take a PEM encoded DER encoded cert and return the *x509.Certificate object
func ReadCrt(certPEMBytes []byte) *x509.Certificate {
	pemBlock, _ := pem.Decode(certPEMBytes)
	cert, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse certificate")
	}
	return cert
}

// Output a pem and der encoded ed25519 private key
func GenKey() []byte {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to generate key")
	}

	privDER, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to marshal private key")
	}

	privPem := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privDER,
	})

	return privPem
}

// Take a PEM encoded DER encoded csr and return the *x509.CertificateRequest object
func ReadCSR(certPEMBytes []byte) *x509.CertificateRequest {
	pemBlock, _ := pem.Decode(certPEMBytes)
	csr, err := x509.ParseCertificateRequest(pemBlock.Bytes)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse certificate request")
	}
	return csr
}

// Take a privat key object and a hostname and return a PEM encoded CA with the commonname of the hostname
func GenCA(caKey ed25519.PrivateKey) []byte {
	caTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: commonName,
		},
		NotBefore:             today,
		NotAfter:              CaExpirey,
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	caBytes, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, caKey.Public(), caKey)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create CA")
	}
	caPem := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	return caPem
}

// Take a private key, a hostname, and a path, create a CSR with the CN and the DNSNames as the hostname and return a PEM encoded CSR
func GenCSR(privKey any, hostname string) []byte {
	csrTemplate := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName: hostname,
		},
		DNSNames: []string{hostname},
	}
	csr, err := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, privKey)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create CSR")
	}
	csrPem := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csr,
	})

	return csrPem
}

// Take a CA cert object, private key of the CA, a CSR, and a path to write the cert
// Generate a cert from the CSA and sign with the CA supplied
func SignCrt(caCert *x509.Certificate, caPriv ed25519.PrivateKey, csr *x509.CertificateRequest) []byte {
	certBytes, err := x509.CreateCertificate(rand.Reader, &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: csr.Subject.CommonName,
		},
		NotBefore:   today,
		NotAfter:    CertExpirey,
		DNSNames:    csr.DNSNames,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}, caCert, csr.PublicKey, caPriv)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create certificate")
	}

	certPem := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	return certPem
}

// Verify if a cert was signed by the CA
func VerifyCrt(caCert *x509.Certificate, cert *x509.Certificate) error {
	caPool := x509.NewCertPool()
	caPool.AddCert(caCert)

	opts := x509.VerifyOptions{
		Roots:     caPool,
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	if _, err := cert.Verify(opts); err != nil {
		return err
	}

	return nil
}

// Return the sum of any PEM encoded object
func PemSum(pemBytes []byte) string {
	pem, _ := pem.Decode(pemBytes)
	shaBytes := sha1.Sum(pem.Bytes)
	return hex.EncodeToString(shaBytes[:])
}
