package prober

import (
	"crypto/tls"
	"time"
)

func ProbeSSLHost(network, hostWithPort string) (expiry time.Time, sslErr error) {
	conf := &tls.Config{
		InsecureSkipVerify: false,
	}
	conn, sslErr := tls.Dial(network, hostWithPort, conf)
	if sslErr != nil {
		return
	}
	defer conn.Close()
	cert0 := conn.ConnectionState().PeerCertificates[0]
	logger.Debug("\tNotAfter: %v", cert0.NotAfter)
	logger.Debug("\tNotBefore: %v", cert0.NotBefore)
	logger.Debug("\tCommonName: %v", cert0.Subject.CommonName)

	expiry = cert0.NotAfter
	return
}
