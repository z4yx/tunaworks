package server

import internal "github.com/z4yx/tunaworks/internal"

func (s *Server) UpdateNodeProtocol(token string, proto int) (err error) {

	r, err := s.db.Exec("UPDATE nodes SET proto=?,heartbeat=datetime('now') WHERE token=?",
		proto, token)
	if err != nil {
		return
	}
	affected, _ := r.RowsAffected()
	logger.Debug("RowsAffected %d", affected)
	return
}

func (s *Server) QueryNodes(active bool) (nodes map[int]internal.NodeInfo, err error) {
	where := ""
	if active {
		where = "WHERE active=1"
	}

	rows, err := s.db.Query("SELECT node, name, proto FROM nodes " + where)
	if err != nil {
		return
	}
	defer rows.Close()
	nodes = make(map[int]internal.NodeInfo)
	for rows.Next() {
		var info internal.NodeInfo
		var id int
		err = rows.Scan(&id, &info.Name, &info.Proto)
		if err != nil {
			return
		}
		nodes[id] = info
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
	//logger.Debug("rec %p %v",s.db,rec)
	r, err := s.db.Exec(`INSERT INTO records(http_code, response_time, site, node, protocol, ssl_err, ssl_expire)
VALUES (?, ?, ?, ?, ?, ?, ?)`, rec.StatusCode, rec.ResponseTime, rec.WebsiteId, node, rec.Protocol, rec.SSLError, rec.SSLExpire)
	if err == nil {
		return err
	}
	affected, _ := r.RowsAffected()
	logger.Debug("RowsAffected %d", affected)
	return err
}

func (s *Server) InsertRecordWithSSLInfo(node int, rec *internal.ProbeResult) error {
	//logger.Debug("rec %p %v",s.db,rec)
	r, err := s.db.Exec(`INSERT INTO records(http_code, response_time, site, node, protocol, ssl_err, ssl_expire,
		have_ocsp,ocsp_err,ocsp_this_update,ocsp_next_update)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		rec.StatusCode, rec.ResponseTime, rec.WebsiteId, node, rec.Protocol, rec.SSLError, rec.SSLExpire,
		rec.SSLInfo.HaveOCSPStapling, rec.SSLInfo.OCSPStaplingErr, rec.SSLInfo.OCSPThisUpdate, rec.SSLInfo.OCSPNextUpdate)
	if err == nil {
		return err
	}
	affected, _ := r.RowsAffected()
	logger.Debug("RowsAffected %d", affected)
	return err
}

func (s *Server) QueryLatestMonitorInfo() (ret *internal.LatestMonitorInfo, err error) {
	nodeInfo, err := s.QueryNodes(true)
	if err != nil {
		return
	}
	node2name := make(map[int]string, len(nodeInfo))
	for i, val := range nodeInfo {
		node2name[i] = val.Name
	}
	rows, err := s.db.Query(`SELECT records.updated,tmp.site,url,tmp.node,tmp.protocol,http_code,response_time,ssl_err,ssl_expire 
	FROM (SELECT site,node,protocol,MAX(updated) AS u FROM records GROUP BY site,node,protocol) tmp
	INNER JOIN records ON tmp.site = records.site
		AND tmp.site = records.site 
		AND tmp.node = records.node 
		AND tmp.protocol = records.protocol 
		AND tmp.u = records.updated
	INNER JOIN sites ON tmp.site = sites.site AND sites.active = 1`)
	if err != nil {
		return
	}
	defer rows.Close()
	ret = &internal.LatestMonitorInfo{
		Websites:  make([]internal.WebsiteInfo, 0, 10),
		NodeNames: node2name,
		NodeInfo:  nodeInfo,
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
