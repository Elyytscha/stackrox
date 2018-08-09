package mtls

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/cloudflare/cfssl/config"
	cfcsr "github.com/cloudflare/cfssl/csr"
	cfsigner "github.com/cloudflare/cfssl/signer"
	"github.com/cloudflare/cfssl/signer/local"
	siStore "github.com/stackrox/rox/central/serviceidentities/store"
	"github.com/stackrox/rox/generated/api/v1"
)

const (
	certsPrefix = "/run/secrets/stackrox.io/certs/"
	// caCertFilePath is where the certificate is stored.
	caCertFilePath = certsPrefix + "ca.pem"
	// caKeyFilePath is where the key is stored.
	caKeyFilePath = certsPrefix + "ca-key.pem"

	// certFilePath is where the certificate is stored.
	certFilePath = certsPrefix + "cert.pem"
	// keyFilePath is where the key is stored.
	keyFilePath = certsPrefix + "key.pem"

	// To account for clock skew, set certificates to be valid some time in the past.
	beforeGracePeriod = 1 * time.Hour

	certLifetime = 365 * 24 * time.Hour
)

var (
	// CentralCN is the Common Name used in certificates for Central.
	CentralCN = CommonName{ServiceType: v1.ServiceType_CENTRAL_SERVICE, Identifier: "Central"}

	readCAOnce sync.Once
	caCert     *x509.Certificate
	caCertDER  []byte
	caCertErr  error
)

// LeafCertificateFromFile reads a tls.Certificate (including private key and cert)
// from the canonical locations on non-central services.
func LeafCertificateFromFile() (tls.Certificate, error) {
	return tls.LoadX509KeyPair(certFilePath, keyFilePath)
}

// CACertDER reads the PEM-decoded bytes of the cert from the local file system.
func CACertDER() ([]byte, error) {
	b, err := ioutil.ReadFile(caCertFilePath)
	if err != nil {
		return nil, fmt.Errorf("file access: %s", err)
	}
	decoded, _ := pem.Decode(b)
	if decoded == nil {
		return nil, fmt.Errorf("invalid PEM")
	}
	return decoded.Bytes, nil
}

// CACertPEM returns the PEM-encoded CA certificate.
func CACertPEM() ([]byte, error) {
	_, caDER, err := CACert()
	if err != nil {
		return nil, fmt.Errorf("CA cert loading: %s", err)
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caDER,
	}), nil
}

// CACert reads the cert from the local file system and returns the cert and the DER encoding.
func CACert() (*x509.Certificate, []byte, error) {
	readCAOnce.Do(func() {
		der, err := CACertDER()
		if err != nil {
			caCertErr = fmt.Errorf("CA cert could not be decoded: %s", err)
			return
		}
		cert, err := x509.ParseCertificate(der)
		if err != nil {
			caCertErr = fmt.Errorf("CA cert could not be parsed: %s", err)
			return
		}
		caCert = cert
		caCertDER = der
	})
	return caCert, caCertDER, caCertErr
}

func signer() (cfsigner.Signer, error) {
	return local.NewSignerFromFile(caCertFilePath, caKeyFilePath, signingPolicy())
}

func signingPolicy() *config.Signing {
	return &config.Signing{
		Default: &config.SigningProfile{
			Usage:    []string{"signing", "key encipherment", "server auth", "client auth"},
			Expiry:   certLifetime + beforeGracePeriod,
			Backdate: beforeGracePeriod,
			CSRWhitelist: &config.CSRWhitelist{
				PublicKey:          true,
				PublicKeyAlgorithm: true,
				SignatureAlgorithm: true,
			},
		},
	}
}

// IssueNewCert generates a new key and certificate chain for a sensor.
func IssueNewCert(cn CommonName, storage siStore.Store) (certPEM, keyPEM []byte, identity *v1.ServiceIdentity, err error) {
	returnErr := func(err error, prefix string) ([]byte, []byte, *v1.ServiceIdentity, error) {
		return nil, nil, nil, fmt.Errorf("%s: %s", prefix, err)
	}

	s, err := signer()
	if err != nil {
		return returnErr(err, "signer creation")
	}

	serial, err := randomSerial()
	if err != nil {
		return returnErr(err, "serial generation")
	}
	csr := &cfcsr.CertificateRequest{
		KeyRequest: cfcsr.NewBasicKeyRequest(),
	}
	csrBytes, keyBytes, err := cfcsr.ParseRequest(csr)
	if err != nil {
		return returnErr(err, "request parsing")
	}
	req := cfsigner.SignRequest{
		Request: string(csrBytes),
		Subject: &cfsigner.Subject{
			CN:           cn.String(),
			Names:        []cfcsr.Name{{OU: ou(cn.ServiceType)}},
			SerialNumber: strconv.FormatInt(serial, 10),
		},
	}
	certBytes, err := s.Sign(req)
	if err != nil {
		return returnErr(err, "signing")
	}

	certPEM = certBytes
	keyPEM = keyBytes

	id := generateIdentity(cn, serial)
	if storage != nil {
		err = storage.AddServiceIdentity(id)
		if err != nil {
			return returnErr(err, "identity storage")
		}
	}

	return certPEM, keyPEM, id, nil
}

func randomSerial() (int64, error) {
	serial, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return 0, fmt.Errorf("serial number generation: %s", err)
	}
	return serial.Int64(), nil
}

func ou(t v1.ServiceType) string {
	return t.String()
}

func generateIdentity(cn CommonName, serial int64) *v1.ServiceIdentity {
	return &v1.ServiceIdentity{
		Id:     cn.Identifier,
		Type:   cn.ServiceType,
		Serial: serial,
	}
}
