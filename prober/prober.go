package prober

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	internal "github.com/z4yx/tunaworks/internal"
)

type ProberCtx struct {
	allWebsites internal.AllWebsites
	cfg         *ProberConfig
}

func (ctx *ProberCtx) getWebsites() error {
	resp, err := http.Get("https://" + ctx.cfg.Server + "/prober/websites")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(body), &ctx.allWebsites)

	return err
}

func (ctx *ProberCtx) reportResult(result *internal.ProbeResult) error {
	jbytes, err := json.Marshal(result)
	if err != nil {
		return err
	}
	resp, err := http.Post("https://"+ctx.cfg.Server+"/prober/result", "application/json", bytes.NewBuffer(jbytes))
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
		logger.Debug("Probing %d %s", site.Id, site.Url)
		u, err := url.Parse(site.Url)
		if err != nil {
			logger.Error("Invalid url %s", err.Error())
			continue
		}
		result := internal.ProbeResult{
			WebsiteId: site.Id,
		}
		for _, network := range networks {
			logger.Debug(network)
			if network == "tcp4" {
				result.Protocol = 4
			} else if network == "tcp6" {
				result.Protocol = 6
			}
			if u.Scheme == "https" {
				expiry, sslErr := ProbeSSLHost(network, u.Host)
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
				result.SSLExpire = expiry
			}
			statusCode, responseTime, httpErr := ProbeHttpHost(network, site.Url)
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
						Int64: int64(responseTime),
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

	ctx.getWebsites()

	for {
		select {
		case <-tickerUpdateList.C:
			ctx.getWebsites()
		case <-tickerProbe.C:
			ctx.probeWebsites()
		}
	}
}

func MakeProber(cfg *ProberConfig) *ProberCtx {
	ret := &ProberCtx{
		cfg: cfg,
	}
	return ret
}
