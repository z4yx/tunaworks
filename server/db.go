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

func (s *Server) QuerySites(active bool) (sites map[int]string, err error) {
	where := ""
	if active {
		where = "WHERE active=1"
	}

	rows, err := s.db.Query("SELECT site, url FROM sites " + where)
	if err != nil {
		return
	}
	defer rows.Close()
	sites = make(map[int]string)
	for rows.Next() {
		var url string
		var id int
		err = rows.Scan(&id, &url)
		if err != nil {
			return
		}
		sites[id] = url
	}
	return
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
	rows, err := s.db.Query(`SELECT updated,tmp.site,url,node,http_code,response_time,ssl_err,ssl_expire
	FROM (SELECT * FROM records ORDER BY site,node,updated DESC) tmp
	INNER JOIN sites
	ON tmp.site = sites.site AND sites.active = 1
	GROUP BY tmp.site, tmp.node`)
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
				Id:    site,
				Url:   url,
				Nodes: make(map[int]internal.MonitorRec),
			}
			lastSite = site
		}
		siteInfo.Nodes[node] = record
	}
	if siteInfo != nil {
		ret.Websites = append(ret.Websites, *siteInfo)
	}
	return
}
