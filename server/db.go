package server

import internal "github.com/z4yx/tunaworks/internal"

func (s *Server) QueryNodes(active bool) (nodes map[int]string, err error) {
	where := ""
	if active {
		where = "WHERE active=1"
	}

	rows, err := s.db.Query("SELECT node, name FROM nodes " + where)
	if err != nil {
		return
	}
	defer rows.Close()
	nodes = make(map[int]string)
	for rows.Next() {
		var name string
		var id int
		err = rows.Scan(&id, &name)
		if err != nil {
			return
		}
		nodes[id] = name
	}
	return
}

func (s *Server) QuerySites(active bool) (sites internal.AllWebsites, err error) {
	where := ""
	if active {
		where = "WHERE active=1"
	}

	rows, err := s.db.Query("SELECT site, url FROM sites " + where)
	if err != nil {
		return
	}
	defer rows.Close()

	sites.Websites = make([]internal.Website, 0, 10)
	for rows.Next() {
		var url string
		var id int
		err = rows.Scan(&id, &url)
		if err != nil {
			return
		}
		sites.Websites = append(sites.Websites, internal.Website{
			Id:  id,
			Url: url,
		})
	}
	return
}

func (s *Server) AuthNode(token string) (succ bool, node int) {
	succ = false
	node = -1
	rows, err := s.db.Query("SELECT node FROM nodes WHERE token=?", token)
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&node)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		return true, node
	}
	return
}

func (s *Server) InsertRecord(node int, rec *internal.ProbeResult) error {
	r, err := s.db.Exec(`INSERT INTO records(http_code, response_time, site, node, protocol, ssl_err, ssl_expire)
VALUES (?, ?, ?, ?, ?, ?, ?)`, rec.StatusCode, rec.ResponseTime, rec.WebsiteId, node, rec.Protocol, rec.SSLError, rec.SSLExpire)
	affected, _ := r.RowsAffected()
	logger.Debug("RowsAffected %d", affected)
	return err
}

func (s *Server) QueryLatestMonitorInfo() (ret *internal.LatestMonitorInfo, err error) {
	node2name, err := s.QueryNodes(true)
	if err != nil {
		return
	}
	// nodes := make([]string, len(node2name))
	// for i, val := range node2name {
	// 	nodes[i] = val
	// }
	rows, err := s.db.Query(`SELECT updated,tmp.site,url,node,protocol,http_code,response_time,ssl_err,ssl_expire
	FROM (SELECT * FROM records ORDER BY site,node,protocol,updated DESC) tmp
	INNER JOIN sites
	ON tmp.site = sites.site AND sites.active = 1
	GROUP BY tmp.site, tmp.node, tmp.protocol`)
	if err != nil {
		return
	}
	defer rows.Close()
	ret = &internal.LatestMonitorInfo{
		Websites:  make([]internal.WebsiteInfo, 0, 10),
		NodeNames: node2name,
	}
	lastSite := -1
	var siteInfo *internal.WebsiteInfo
	for rows.Next() {
		var record internal.MonitorRec
		var site, node int
		var url string
		var exist bool
		err = rows.Scan(
			&record.Updated,
			&site,
			&url,
			&node,
			&record.Protocol,
			&record.StatusCode,
			&record.ResponseTime,
			&record.SSLError,
			&record.SSLExpire)
		if err != nil {
			return
		}
		if record.Name, exist = node2name[node]; !exist {
			// Skip inactive nodes
			continue
		}
		if site != lastSite {
			if siteInfo != nil {
				ret.Websites = append(ret.Websites, *siteInfo)
			}
			siteInfo = &internal.WebsiteInfo{
				Id:     site,
				Url:    url,
				Nodes4: make(map[int]internal.MonitorRec),
				Nodes6: make(map[int]internal.MonitorRec),
			}
			lastSite = site
		}
		if record.Protocol == 4 {
			siteInfo.Nodes4[node] = record
		} else if record.Protocol == 6 {
			siteInfo.Nodes6[node] = record
		}
	}
	if siteInfo != nil {
		ret.Websites = append(ret.Websites, *siteInfo)
	}
	return
}
