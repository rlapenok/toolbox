package database

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/rlapenok/toolbox/errors"
)

// BuildTLSConfig Создает конфигурацию TLS
func BuildTLSConfig(caPath, certPath, keyPath string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, errors.New(errors.Internal, err.Error())
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true,
	}

	if caPath != "" {
		caBytes, err := os.ReadFile(caPath)
		if err != nil {
			return nil, errors.New(errors.Internal, err.Error())
		}
		caPool := x509.NewCertPool()
		if !caPool.AppendCertsFromPEM(caBytes) {
			return nil, errors.New(errors.Internal, "invalid CA PEM")
		}
		tlsConfig.RootCAs = caPool
		tlsConfig.InsecureSkipVerify = false
	}

	return tlsConfig, nil
}
