package prober

import (
	"context"
	"net"
	"net/http"
	"time"
)

func ProbeHttpHost(network, url string) (statusCode int, responseTime time.Duration, httpErr error) {
	myClient := http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _network, addr string) (net.Conn, error) {
				logger.Debug("DialContext %s %s -> %s", addr, _network, network)
				myDial := &net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 0,
					DualStack: false,
				}
				return myDial.DialContext(ctx, network, addr)
			},
			DisableKeepAlives: true,
		},
	}
	req, httpErr := http.NewRequest("GET", url, nil)
	if httpErr != nil {
		return
	}
	start := time.Now()
	res, httpErr := myClient.Do(req)
	if httpErr != nil {
		return
	}
	defer res.Body.Close()

	responseTime = time.Since(start)
	// body, _ := ioutil.ReadAll(res.Body)
	statusCode = res.StatusCode
	return
}
