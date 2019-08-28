package internal

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type NullInt64 struct {
	sql.NullInt64
}

type NullString struct {
	sql.NullString
}

func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

func (ni *NullInt64) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		ni.Valid = false
		return nil
	}
	fmt.Printf("data %s\n", data)
	ni.Valid = true
	return json.Unmarshal(data, &ni.Int64)
}

func (ni NullString) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.String)
}

func (ni *NullString) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		ni.Valid = false
		return nil
	}
	ni.Valid = true
	return json.Unmarshal(data, &ni.String)
}

type MonitorRec struct {
	Name         string
	Protocol     int
	StatusCode   NullInt64
	ResponseTime NullInt64
	SSLError     NullString
	SSLExpire    time.Time
	Updated      time.Time
}

type WebsiteInfo struct {
	Id     int
	Url    string
	Nodes4 map[int]MonitorRec
	Nodes6 map[int]MonitorRec
}

type NodeInfo struct {
	Name      string
	Heartbeat time.Time
	Proto     int
}

type LatestMonitorInfo struct {
	NodeNames map[int]string
	NodeInfo  map[int]NodeInfo
	Websites  []WebsiteInfo
}

type ProbeResult struct {
	WebsiteId    int
	Protocol     int
	StatusCode   NullInt64
	ResponseTime NullInt64
	SSLError     NullString
	SSLExpire    time.Time
}

type Website struct {
	Id  int
	Url string
}

type AllWebsites struct {
	Websites []Website
}
