package extensions

import (
	"bytes"
	"testing"

	"github.com/stackrox/rox/operator/api/central/v1alpha1"
	common "github.com/stackrox/rox/operator/api/common/v1alpha1"
	"github.com/stackrox/rox/pkg/auth/htpasswd"
	"github.com/stackrox/rox/pkg/grpc/authn/basic"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestReconcileAdminPassword(t *testing.T) {
	hf := htpasswd.New()
	require.NoError(t, hf.Set(basic.DefaultUsername, "foobar"))
	var buf bytes.Buffer
	require.NoError(t, hf.Write(&buf))

	htpasswdWithSomePassword := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "central-htpasswd",
			Namespace: testNamespace,
		},
		Data: map[string][]byte{
			"htpasswd": buf.Bytes(),
		},
	}

	htpasswdWithNoPassword := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "central-htpasswd",
			Namespace: testNamespace,
		},
		Data: map[string][]byte{},
	}

	plaintextPasswordSecret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-password",
			Namespace: testNamespace,
		},
		Data: map[string][]byte{
			"password": []byte("foobarbaz\n"),
		},
	}

	cases := map[string]secretReconciliationTestCase{
		"If no central-htpasswd secret exists and no plaintext secret reference was specified, a password should be automatically generated": {
			ExpectedCreatedSecrets: map[string]secretVerifyFunc{
				"central-htpasswd": func(t *testing.T, data secretDataMap) {
					plaintextPW := string(data[adminPasswordKey])
					require.NotEmpty(t, plaintextPW)

					htpasswdBytes := data[htpasswdKey]
					hf, err := htpasswd.ReadHashFile(bytes.NewReader(htpasswdBytes))
					require.NoError(t, err)

					assert.True(t, hf.Check(basic.DefaultUsername, plaintextPW))
				},
			},
			VerifyStatus: func(t *testing.T, status *v1alpha1.CentralStatus) {
				require.NotNil(t, status.Central)
				require.NotNil(t, status.Central.AdminPassword)
				assert.Contains(t, status.Central.AdminPassword.Info, "A password for the 'admin' user has been automatically generated and stored")
			},
		},
		"If a central-htpasswd secret with a password exists, no password should be generated": {
			Existing: []*v1.Secret{htpasswdWithSomePassword},
			VerifyStatus: func(t *testing.T, status *v1alpha1.CentralStatus) {
				require.NotNil(t, status.Central)
				require.NotNil(t, status.Central.AdminPassword)
				assert.Contains(t, status.Central.AdminPassword.Info, "A user-defined central-htpasswd secret was found, containing htpasswd-encoded credentials.")
			},
		},
		"If a central-htpasswd secret with no password exists, no password should be generated and the user should be informed that basic auth is disabled": {
			Existing: []*v1.Secret{htpasswdWithNoPassword},
			VerifyStatus: func(t *testing.T, status *v1alpha1.CentralStatus) {
				require.NotNil(t, status.Central)
				require.NotNil(t, status.Central.AdminPassword)
				assert.Contains(t, status.Central.AdminPassword.Info, "Login with username/password has been disabled")
			},
		},
		"If a secret with a plaintext password is referenced, a central-htpasswd secret should be created accordingly": {
			Spec: v1alpha1.CentralSpec{
				Central: &v1alpha1.CentralComponentSpec{
					AdminPasswordSecret: &common.LocalSecretReference{
						Name: plaintextPasswordSecret.Name,
					},
				},
			},
			Existing: []*v1.Secret{plaintextPasswordSecret},
			ExpectedCreatedSecrets: map[string]secretVerifyFunc{
				"central-htpasswd": func(t *testing.T, data secretDataMap) {
					htpasswdBytes := data[htpasswdKey]
					hf, err := htpasswd.ReadHashFile(bytes.NewReader(htpasswdBytes))
					require.NoError(t, err)

					assert.True(t, hf.Check(basic.DefaultUsername, "foobarbaz"))
				},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			testSecretReconciliation(t, reconcileAdminPassword, c)
		})
	}
}
