package extensions

import (
	"context"
	"testing"

	pkgErrors "github.com/pkg/errors"
	centralV1Alpha1 "github.com/stackrox/rox/operator/apis/platform/v1alpha1"
	"github.com/stackrox/rox/pkg/uuid"
	"github.com/stretchr/testify/suite"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

type secretReconcilerExtensionTestSuite struct {
	suite.Suite

	centralObj *centralV1Alpha1.Central
	k8sClient  kubernetes.Interface

	reconcileExt *secretReconciliationExtension
}

func TestSecretReconcilerExtension(t *testing.T) {
	suite.Run(t, new(secretReconcilerExtensionTestSuite))
}

func (s *secretReconcilerExtensionTestSuite) SetupTest() {
	s.centralObj = &centralV1Alpha1.Central{
		TypeMeta: metav1.TypeMeta{
			APIVersion: centralV1Alpha1.CentralGVK.GroupVersion().String(),
			Kind:       centralV1Alpha1.CentralGVK.Kind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "stackrox-central-services",
			Namespace: testNamespace,
			UID:       types.UID(uuid.NewV4().String()),
		},
	}

	existingSecret := &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "existing-secret",
			Namespace: testNamespace,
		},
		Data: map[string][]byte{
			"secret-name": []byte("existing-secret"),
			"managed":     []byte("false"),
		},
	}

	existingOwnedSecret := &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "existing-managed-secret",
			Namespace: testNamespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(s.centralObj, centralV1Alpha1.CentralGVK),
			},
		},
		Data: map[string][]byte{
			"secret-name": []byte("existing-managed-secret"),
			"managed":     []byte("true"),
		},
	}

	s.k8sClient = fake.NewSimpleClientset(existingSecret, existingOwnedSecret)

	s.reconcileExt = &secretReconciliationExtension{
		ctx:        context.Background(),
		k8sClient:  s.k8sClient,
		centralObj: s.centralObj,
	}
}

func (s *secretReconcilerExtensionTestSuite) Test_ShouldNotExist_OnNonExisting_ShouldDoNothing() {
	validateFn := func(secretDataMap, bool) error {
		s.Require().Fail("this function should not be called")
		panic("unexpected")
	}
	generateFn := func() (secretDataMap, error) {
		s.Require().Fail("this function should not be called")
		panic("unexpected")
	}

	err := s.reconcileExt.reconcileSecret("absent-secret", false, validateFn, generateFn, false)
	s.Require().NoError(err)

	_, err = s.k8sClient.CoreV1().Secrets(testNamespace).Get(context.Background(), "absent-secret", metav1.GetOptions{})
	s.True(errors.IsNotFound(err))
}

func (s *secretReconcilerExtensionTestSuite) Test_ShouldNotExist_OnExistingManaged_ShouldDelete() {
	validateFn := func(secretDataMap, bool) error {
		s.Require().Fail("this function should not be called")
		panic("unexpected")
	}
	generateFn := func() (secretDataMap, error) {
		s.Require().Fail("this function should not be called")
		panic("unexpected")
	}

	err := s.reconcileExt.reconcileSecret("existing-managed-secret", false, validateFn, generateFn, false)
	s.Require().NoError(err)

	_, err = s.k8sClient.CoreV1().Secrets(testNamespace).Get(context.Background(), "existing-managed-secret", metav1.GetOptions{})
	s.True(errors.IsNotFound(err))
}

func (s *secretReconcilerExtensionTestSuite) Test_ShouldNotExist_OnExistingUnmanaged_ShouldDoNothing() {
	validateFn := func(secretDataMap, bool) error {
		s.Require().Fail("this function should not be called")
		panic("unexpected")
	}
	generateFn := func() (secretDataMap, error) {
		s.Require().Fail("this function should not be called")
		panic("unexpected")
	}

	err := s.reconcileExt.reconcileSecret("existing-secret", false, validateFn, generateFn, false)
	s.Require().NoError(err)

	_, err = s.k8sClient.CoreV1().Secrets(testNamespace).Get(context.Background(), "existing-managed-secret", metav1.GetOptions{})
	s.NoError(err)
}

func (s *secretReconcilerExtensionTestSuite) Test_ShouldExist_OnNonExisting_ShouldCreateSecretWithOwnerRef() {
	validateFn := func(secretDataMap, bool) error {
		s.Require().Fail("this function should not be called")
		panic("unexpected")
	}
	// this ensures that we check for the existence of a unique created secret
	var markerID string
	generateFn := func() (secretDataMap, error) {
		markerID = uuid.NewV4().String()
		return secretDataMap{
			"generated": []byte(markerID),
		}, nil
	}

	err := s.reconcileExt.reconcileSecret("absent-secret", true, validateFn, generateFn, false)
	s.Require().NoError(err)
	s.NotEmpty(markerID, "generate function has not been called")

	secret, err := s.k8sClient.CoreV1().Secrets(testNamespace).Get(context.Background(), "absent-secret", metav1.GetOptions{})
	s.Require().NoError(err)

	s.EqualValues(secret.GetOwnerReferences(), []metav1.OwnerReference{*metav1.NewControllerRef(s.centralObj, centralV1Alpha1.CentralGVK)})

	s.Equal(markerID, string(secret.Data["generated"]))
}

func (s *secretReconcilerExtensionTestSuite) Test_ShouldExist_OnExistingManaged_PassingValidation_ShouldDoNothing() {
	initSecret, err := s.k8sClient.CoreV1().Secrets(testNamespace).Get(context.Background(), "existing-managed-secret", metav1.GetOptions{})
	s.Require().NoError(err)

	validated := false
	validateFn := func(data secretDataMap, managed bool) error {
		s.Equal("existing-managed-secret", string(data["secret-name"]))
		s.True(managed)
		validated = true
		return nil
	}

	generateFn := func() (secretDataMap, error) {
		s.Require().Fail("this function should not be called")
		panic("unexpected")
	}

	err = s.reconcileExt.reconcileSecret("existing-managed-secret", true, validateFn, generateFn, false)
	s.Require().NoError(err)
	s.True(validated)

	secret, err := s.k8sClient.CoreV1().Secrets(testNamespace).Get(context.Background(), "existing-managed-secret", metav1.GetOptions{})
	s.Require().NoError(err)

	s.Equal(initSecret, secret)
}

func (s *secretReconcilerExtensionTestSuite) Test_ShouldExist_OnExistingManaged_FailingValidation_NoFixExisting_ShouldFail() {
	failValidationErr := pkgErrors.New("failed validation")
	validateFn := func(data secretDataMap, managed bool) error {
		s.Equal("existing-managed-secret", string(data["secret-name"]))
		s.True(managed)
		return failValidationErr
	}

	generateFn := func() (secretDataMap, error) {
		return secretDataMap{
			"new-secret-data": []byte("foo"),
		}, nil
	}

	err := s.reconcileExt.reconcileSecret("existing-managed-secret", true, validateFn, generateFn, false)
	s.ErrorIs(err, failValidationErr)

	secret, err := s.k8sClient.CoreV1().Secrets(testNamespace).Get(context.Background(), "existing-managed-secret", metav1.GetOptions{})
	s.Require().NoError(err)

	s.Equal("existing-managed-secret", string(secret.Data["secret-name"]))
}

func (s *secretReconcilerExtensionTestSuite) Test_ShouldExist_OnExistingManaged_FailingValidation_WithFixExisting_ShouldFix() {
	failValidationErr := pkgErrors.New("failed validation")
	validateFn := func(data secretDataMap, managed bool) error {
		s.Equal("existing-managed-secret", string(data["secret-name"]))
		s.True(managed)
		return failValidationErr
	}

	generateFn := func() (secretDataMap, error) {
		return secretDataMap{
			"new-secret-data": []byte("foo"),
		}, nil
	}

	err := s.reconcileExt.reconcileSecret("existing-managed-secret", true, validateFn, generateFn, true)
	s.NoError(err)

	secret, err := s.k8sClient.CoreV1().Secrets(testNamespace).Get(context.Background(), "existing-managed-secret", metav1.GetOptions{})
	s.Require().NoError(err)

	s.Equal("foo", string(secret.Data["new-secret-data"]))
}

func (s *secretReconcilerExtensionTestSuite) Test_ShouldExist_OnExistingUnmanaged_PassingValidation_ShouldDoNothing() {
	initSecret, err := s.k8sClient.CoreV1().Secrets(testNamespace).Get(context.Background(), "existing-secret", metav1.GetOptions{})
	s.Require().NoError(err)

	validated := false
	validateFn := func(data secretDataMap, managed bool) error {
		s.Equal("existing-secret", string(data["secret-name"]))
		s.False(managed)
		validated = true
		return nil
	}

	generateFn := func() (secretDataMap, error) {
		s.Require().Fail("this function should not be called")
		panic("unexpected")
	}

	err = s.reconcileExt.reconcileSecret("existing-secret", true, validateFn, generateFn, false)
	s.Require().NoError(err)
	s.True(validated)

	secret, err := s.k8sClient.CoreV1().Secrets(testNamespace).Get(context.Background(), "existing-secret", metav1.GetOptions{})
	s.Require().NoError(err)

	s.Equal(initSecret, secret)
}

func (s *secretReconcilerExtensionTestSuite) Test_ShouldExist_OnExistingUnmanaged_FailingValidation_ShouldDoNothingAndFail() {
	initSecret, err := s.k8sClient.CoreV1().Secrets(testNamespace).Get(context.Background(), "existing-secret", metav1.GetOptions{})
	s.Require().NoError(err)

	failValidationErr := pkgErrors.New("failed validation")
	validateFn := func(data secretDataMap, managed bool) error {
		s.Equal("existing-secret", string(data["secret-name"]))
		s.False(managed)
		return failValidationErr
	}

	generateFn := func() (secretDataMap, error) {
		s.Require().Fail("this function should not be called")
		panic("unexpected")
	}

	err = s.reconcileExt.reconcileSecret("existing-secret", true, validateFn, generateFn, false)
	s.ErrorIs(err, failValidationErr)

	secret, err := s.k8sClient.CoreV1().Secrets(testNamespace).Get(context.Background(), "existing-secret", metav1.GetOptions{})
	s.Require().NoError(err)

	s.Equal(initSecret, secret)
}
