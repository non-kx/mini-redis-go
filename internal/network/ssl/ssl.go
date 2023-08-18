package ssl

import (
	"crypto/tls"
)

func GetServerTlsConfig(publickey string, privatekey string) (*tls.Config, error) {
	cert, err := loadCertificate(publickey, privatekey)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{*cert},
		ClientAuth:   tls.RequireAnyClientCert,
		MinVersion:   tls.VersionTLS12,
		MaxVersion:   tls.VersionTLS13,
	}

	return tlsConfig, nil
}

func GetClientTlsConfig(publickey string, privatekey string) (*tls.Config, error) {
	cert, err := loadCertificate(publickey, privatekey)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{*cert},
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
		MaxVersion:         tls.VersionTLS13,
	}

	return tlsConfig, nil
}

func loadCertificate(publicKeyFile string, privateKeyFile string) (*tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(publicKeyFile, privateKeyFile)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}
