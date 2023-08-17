package ssl

import (
	"crypto/tls"
)

func LoadCertificate(publicKeyFile string, privateKeyFile string) (*tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(publicKeyFile, privateKeyFile)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}
