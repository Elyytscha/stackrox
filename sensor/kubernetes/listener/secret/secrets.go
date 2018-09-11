package secret

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/cloudflare/cfssl/certinfo"
	pkgV1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/listeners"
	"github.com/stackrox/rox/pkg/logging"
	"github.com/stackrox/rox/pkg/protoconv"
	"github.com/stackrox/rox/sensor/kubernetes/listener/watchlister"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
)

var logger = logging.LoggerForModule()

// WatchLister implements the WatchLister interface
type WatchLister struct {
	watchlister.WatchLister
	eventC chan<- *listeners.EventWrap
}

// NewWatchLister implements the watch for secrets
func NewWatchLister(client rest.Interface, eventC chan<- *listeners.EventWrap, resyncPeriod time.Duration) *WatchLister {
	npwl := &WatchLister{
		WatchLister: watchlister.NewWatchLister(client, resyncPeriod),
		eventC:      eventC,
	}
	npwl.SetupWatch("secrets", &v1.Secret{}, npwl.resourceChanged)
	return npwl
}

var dataTypeMap = map[string]pkgV1.SecretType{
	"-----BEGIN CERTIFICATE-----":              pkgV1.SecretType_PUBLIC_CERTIFICATE,
	"-----BEGIN NEW CERTIFICATE REQUEST-----":  pkgV1.SecretType_CERTIFICATE_REQUEST,
	"-----BEGIN PRIVACY-ENHANCED MESSAGE-----": pkgV1.SecretType_PRIVACY_ENHANCED_MESSAGE,
	"-----BEGIN OPENSSH PRIVATE KEY-----":      pkgV1.SecretType_OPENSSH_PRIVATE_KEY,
	"-----BEGIN PGP PRIVATE KEY BLOCK-----":    pkgV1.SecretType_PGP_PRIVATE_KEY,
	"-----BEGIN EC PRIVATE KEY-----":           pkgV1.SecretType_EC_PRIVATE_KEY,
	"-----BEGIN RSA PRIVATE KEY-----":          pkgV1.SecretType_RSA_PRIVATE_KEY,
	"-----BEGIN DSA PRIVATE KEY-----":          pkgV1.SecretType_DSA_PRIVATE_KEY,
	"-----BEGIN PRIVATE KEY-----":              pkgV1.SecretType_CERT_PRIVATE_KEY,
	"-----BEGIN ENCRYPTED PRIVATE KEY-----":    pkgV1.SecretType_ENCRYPTED_PRIVATE_KEY,
}

func getSecretType(data string) pkgV1.SecretType {
	for dataPrefix, t := range dataTypeMap {
		if strings.HasPrefix(data, dataPrefix) {
			return t
		}
	}
	return pkgV1.SecretType_UNDETERMINED
}

func convertInterfaceSliceToStringSlice(i []interface{}) []string {
	strSlice := make([]string, 0, len(i))
	for _, v := range i {
		strSlice = append(strSlice, fmt.Sprintf("%v", v))
	}
	return strSlice
}

func convertCFSSLName(name certinfo.Name) *pkgV1.CertName {
	return &pkgV1.CertName{
		CommonName:       name.CommonName,
		Country:          name.Country,
		Organization:     name.Organization,
		OrganizationUnit: name.OrganizationalUnit,
		Locality:         name.Locality,
		Province:         name.Province,
		StreetAddress:    name.StreetAddress,
		PostalCode:       name.PostalCode,
		Names:            convertInterfaceSliceToStringSlice(name.Names),
	}
}

func parseCertData(data string) *pkgV1.Cert {
	info, err := certinfo.ParseCertificatePEM([]byte(data))
	if err != nil {
		return nil
	}
	return &pkgV1.Cert{
		Subject:   convertCFSSLName(info.Subject),
		Issuer:    convertCFSSLName(info.Issuer),
		Sans:      info.SANs,
		StartDate: protoconv.ConvertTimeToTimestampOrNil(info.NotBefore),
		EndDate:   protoconv.ConvertTimeToTimestampOrNil(info.NotAfter),
		Algorithm: info.SignatureAlgorithm,
	}
}

func (npwl *WatchLister) populateTypeData(secret *pkgV1.Secret, dataFiles map[string][]byte) {
	for file, rawData := range dataFiles {
		// Try to base64 decode and if it fails then try the raw value
		var secretType pkgV1.SecretType
		var data string
		decoded, err := base64.StdEncoding.DecodeString(string(rawData))
		if err != nil {
			data = string(rawData)
		} else {
			data = string(decoded)
		}
		secretType = getSecretType(data)

		file := &pkgV1.SecretDataFile{
			Name: file,
			Type: secretType,
		}

		switch secretType {
		case pkgV1.SecretType_PUBLIC_CERTIFICATE:
			file.Metadata = &pkgV1.SecretDataFile_Cert{
				Cert: parseCertData(data),
			}
		}
		secret.Files = append(secret.Files, file)
	}
}

func (npwl *WatchLister) resourceChanged(secretObj interface{}, action pkgV1.ResourceAction) {
	secret, ok := secretObj.(*v1.Secret)
	if !ok {
		logger.Errorf("Object %+v is not a valid secret", secretObj)
		return
	}

	// Filter out service account tokens because we have a service account field.
	// Also filter out DockerConfigJson/DockerCfgs because we don't really care about them.
	switch secret.Type {
	case v1.SecretTypeDockerConfigJson, v1.SecretTypeDockercfg, v1.SecretTypeServiceAccountToken:
		return
	}

	protoSecret := &pkgV1.Secret{
		Id:          string(secret.GetUID()),
		Name:        secret.GetName(),
		Namespace:   secret.GetNamespace(),
		Labels:      secret.GetLabels(),
		Annotations: secret.GetAnnotations(),
		CreatedAt:   protoconv.ConvertTimeToTimestamp(secret.GetCreationTimestamp().Time),
	}

	npwl.populateTypeData(protoSecret, secret.Data)

	npwl.eventC <- &listeners.EventWrap{
		SensorEvent: &pkgV1.SensorEvent{
			Id:     string(secret.GetUID()),
			Action: action,
			Resource: &pkgV1.SensorEvent_Secret{
				Secret: protoSecret,
			},
		},
	}
}
