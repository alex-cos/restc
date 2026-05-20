package restc_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alex-cos/restc"
	"github.com/stretchr/testify/assert"
)

func generateTestCerts(t *testing.T) ([]byte, []byte, []byte) {
	t.Helper()

	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	caTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "Test CA"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour),
		IsCA:         true,
	}
	caCertDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	assert.NoError(t, err)
	caCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caCertDER})

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)
	clientTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject:      pkix.Name{CommonName: "Test Client"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour),
	}
	clientCertDER, err := x509.CreateCertificate(rand.Reader, clientTemplate, caTemplate, &key.PublicKey, caKey)
	assert.NoError(t, err)
	clientCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: clientCertDER})
	clientKeyDER, err := x509.MarshalECPrivateKey(key)
	assert.NoError(t, err)
	clientKey := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: clientKeyDER})

	return caCert, clientCert, clientKey
}

func TestWithTLSConfig(t *testing.T) {
	t.Parallel()

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS13,
		ServerName: "api.test.com",
	}

	httpClient := &http.Client{}
	restc.NewWithClient("https://api.test.com",
		httpClient,
		restc.WithTLSConfig(tlsConfig),
	)

	transport := httpClient.Transport.(*http.Transport)
	assert.NotNil(t, transport.TLSClientConfig)
	assert.Equal(t, uint16(tls.VersionTLS13), transport.TLSClientConfig.MinVersion)
	assert.Equal(t, "api.test.com", transport.TLSClientConfig.ServerName)
}

func TestWithMTLS(t *testing.T) {
	t.Parallel()

	caCertPEM, clientCertPEM, clientKeyPEM := generateTestCerts(t)
	tmpDir := t.TempDir()

	caCertFile := filepath.Join(tmpDir, "ca.pem")
	clientCertFile := filepath.Join(tmpDir, "client.pem")
	clientKeyFile := filepath.Join(tmpDir, "client-key.pem")

	err := os.WriteFile(caCertFile, caCertPEM, 0644)
	assert.NoError(t, err)
	err = os.WriteFile(clientCertFile, clientCertPEM, 0644)
	assert.NoError(t, err)
	err = os.WriteFile(clientKeyFile, clientKeyPEM, 0644)
	assert.NoError(t, err)

	httpClient := &http.Client{}
	restc.NewWithClient("https://api.test.com",
		httpClient,
		restc.WithMTLS(caCertFile, clientCertFile, clientKeyFile),
	)

	transport := httpClient.Transport.(*http.Transport)
	assert.NotNil(t, transport.TLSClientConfig)
	assert.Len(t, transport.TLSClientConfig.Certificates, 1)
	assert.NotNil(t, transport.TLSClientConfig.RootCAs)
	assert.Equal(t, uint16(tls.VersionTLS12), transport.TLSClientConfig.MinVersion)
}

func TestWithMTLS_InvalidFiles(t *testing.T) {
	t.Parallel()

	httpClient := &http.Client{}
	restc.NewWithClient("https://api.test.com",
		httpClient,
		restc.WithMTLS("/nonexistent/ca.pem", "/nonexistent/client.pem", "/nonexistent/key.pem"),
	)

	assert.Nil(t, httpClient.Transport)
}
