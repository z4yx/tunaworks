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
		logger.Debug("%v", sslErr)
		return
	}
	defer conn.Close()
	cert0 := conn.ConnectionState().PeerCertificates[0]
	logger.Debug("%v", cert0.NotAfter)
	logger.Debug("%v", cert0.NotBefore)
	logger.Debug("%v", cert0.Subject.CommonName)

	expiry = cert0.NotAfter
	return
}
