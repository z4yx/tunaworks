package prober

import (
	"crypto/tls"
	"net"
	"time"
	internal "github.com/z4yx/tunaworks/internal"
	"golang.org/x/crypto/ocsp"
	"database/sql"
)

func ProbeSSLHost(network, hostWithPort string) (sslInfo internal.SSLInfo, sslErr error) {
	conf := &tls.Config{
		InsecureSkipVerify: false,
	}
	myDial := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 0,
		DualStack: false,
	}
	conn, sslErr := tls.DialWithDialer(myDial, network, hostWithPort, conf)
	if sslErr != nil {
		return
	}
	defer conn.Close()
	cert0 := conn.ConnectionState().PeerCertificates[0]
	logger.Debug("\tNotAfter: %v", cert0.NotAfter)
	logger.Debug("\tNotBefore: %v", cert0.NotBefore)
	logger.Debug("\tCommonName: %v", cert0.Subject.CommonName)

	sslInfo.NotAfter, sslInfo.NotBefore, sslInfo.CommonName = 
	cert0.NotAfter, cert0.NotBefore, cert0.Subject.CommonName

	ocspErrStr := ""
	ocspResp := conn.OCSPResponse()
	sslInfo.HaveOCSPStapling = len(ocspResp) > 0
	if(sslInfo.HaveOCSPStapling){
		resp, ocspErr := ocsp.ParseResponse(ocspResp, nil)
		if ocspErr != nil {
			ocspErrStr = ocspErr.Error()
		} else if resp.Status == ocsp.Revoked {
			ocspErrStr = "Revoked at " + resp.RevokedAt.String()
		} else if resp.Status == ocsp.Unknown {
			ocspErrStr = "Cert Status Unknown"
		} else if resp.Status == ocsp.Good {
			sslInfo.OCSPThisUpdate, sslInfo.OCSPNextUpdate =
			resp.ThisUpdate, resp.NextUpdate
		} else {
			ocspErrStr = "Unknown Status from OCSP Response: " + string(resp.Status)
		}
	}
	if ocspErrStr != "" {
		sslInfo.OCSPStaplingErr = internal.NullString{
			sql.NullString{
				String: ocspErrStr,
				Valid:  true,
			},
		}
	}
	return
}
