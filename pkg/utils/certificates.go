package utils

import (
	"crypto/tls"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"

	"test-task-photo-booth/src/entities"
)

const (
	configTlsCertFilePath = "core.tls.certFilePath"
	configKeyFilePath     = "core.tls.keyFilePath"
)

var (
	ErrNoCertFileFound = errors.New("no certificate file found")
)

func LoadCertificate() (*tls.Certificate, error) {
	wd := viper.GetString(entities.ConfigPathWd)

	certFilePath := viper.GetString(configTlsCertFilePath)
	certPEMBlock, err := GetFileBytes(filepath.Join(wd, certFilePath))
	if err != nil {
		return nil, fmt.Errorf("GetFileBytes() failed: %w", err)
	}

	keyFilePath := viper.GetString(configKeyFilePath)
	keyPEMBlock, err := GetFileBytes(filepath.Join(wd, keyFilePath))
	if err != nil {
		return nil, fmt.Errorf("GetFileBytes() failed: %w", err)
	}

	// Load the certificate and key
	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return nil, fmt.Errorf("tls.X509KeyPair() failed: %w", err)
	}

	return &cert, nil
}
