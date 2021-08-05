package prober

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	internal "github.com/z4yx/tunaworks/internal"
)

type ProberCtx struct {
	allWebsites internal.AllWebsites
	cfg         *ProberConfig
	client      *http.Client
	baseUrl     string
}

func (ctx *ProberCtx) getWebsites() error {
	proto := 0
	if ctx.cfg.IPv4 {
		proto |= 1
	}
	if ctx.cfg.IPv6 {
		proto |= 2
	}
	request, err := http.NewRequest("GET", ctx.baseUrl+"/prober/websites?Proto="+strconv.Itoa(proto), nil)
	if err != nil {
		return err
	}
	request.Header.Add("X-Token", ctx.cfg.Token)
	resp, err := ctx.client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &ctx.allWebsites)

	return err
}

func (ctx *ProberCtx) reportResult(result *internal.ProbeResult) error {
	jbytes, err := json.Marshal(result)
	if err != nil {
		return err
	}
	logger.Debug("result: %s", string(jbytes))
	request, err := http.NewRequest("POST", ctx.baseUrl+"/prober/result", bytes.NewBuffer(jbytes))
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-Token", ctx.cfg.Token)
	resp, err := ctx.client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (ctx *ProberCtx) probeWebsites() {
	networks := make([]string, 0, 2)
	if ctx.cfg.IPv4 {
		networks = append(networks, "tcp4")
	}
	if ctx.cfg.IPv6 {
		networks = append(networks, "tcp6")
	}
	for _, site := range ctx.allWebsites.Websites {
		u, err := url.Parse(site.Url)
		if err != nil {
			logger.Error("Invalid url %s: %s", site.Url, err.Error())
			continue
		}
		for _, network := range networks {
			logger.Debug("Probing %d %s with %s", site.Id, site.Url, network)
			result := internal.ProbeResult{
				WebsiteId: site.Id,
			}
			if network == "tcp4" {
				result.Protocol = 4
			} else if network == "tcp6" {
				result.Protocol = 6
			}
			if u.Port() == "" {
				if u.Scheme == "https" {
					u.Host += ":443"
				} else {
					u.Host += ":80"
				}
			}
			if u.Scheme == "https" {
				sslInfo, sslErr := ProbeSSLHost(network, u.Host)
				logger.Debug("ProbeSSLHost (%v) %v", sslInfo.NotAfter, sslErr)
				if sslErr != nil {
					result.SSLError = internal.NullString{
						sql.NullString{
							String: sslErr.Error(),
							Valid:  true,
						},
					}
					ctx.reportResult(&result)
					continue
				}
				result.SSLExpire = sslInfo.NotAfter
				result.SSLInfo = sslInfo
			}
			UA := "Harmless TUNAworks Prober / " + ctx.cfg.Version
			statusCode, responseTime, httpErr := ProbeHttpHost(network, u.String(), UA)
			logger.Debug("ProbeHttpHost %v %v %v", statusCode, responseTime, httpErr)

			if httpErr != nil {
				result.SSLError = internal.NullString{
					sql.NullString{
						String: httpErr.Error(),
						Valid:  true,
					},
				}
			} else {
				result.StatusCode = internal.NullInt64{
					sql.NullInt64{
						Int64: int64(statusCode),
						Valid: true,
					},
				}
				result.ResponseTime = internal.NullInt64{
					sql.NullInt64{
						Int64: int64(responseTime / time.Millisecond),
						Valid: true,
					},
				}
			}
			ctx.reportResult(&result)
		}
	}
}

func (ctx *ProberCtx) Run() {
	tickerProbe := time.NewTicker(time.Duration(ctx.cfg.Interval) * time.Second)
	tickerUpdateList := time.NewTicker(10 * time.Minute)

	err := ctx.getWebsites()
	if err == nil {
		ctx.probeWebsites()
	} else {
		logger.Error("getWebsites %s", err.Error())
	}

	for {
		select {
		case <-tickerUpdateList.C:
			err = ctx.getWebsites()
			if err != nil {
				logger.Error("getWebsites %s", err.Error())
			}

		case <-tickerProbe.C:
			ctx.probeWebsites()
		}
	}
}

func MakeProber(cfg *ProberConfig) *ProberCtx {
	proto := "http://"
	if cfg.Https {
		proto = "https://"
	}
	ret := &ProberCtx{
		cfg:     cfg,
		baseUrl: proto + cfg.Server,
		client:  &http.Client{Timeout: time.Second * 10},
	}
	return ret
}
