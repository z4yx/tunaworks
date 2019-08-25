package server

import (
	"database/sql"
	"encoding/json"
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

func (ni NullString) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.String)
}

type MonitorRec struct {
	Name         string
	StatusCode   NullInt64
	ResponseTime NullInt64
	SSLError     NullString
	SSLExpire    time.Time
	Updated      time.Time
}

type WebsiteInfo struct {
	Id    int
	Url   string
	Nodes map[int]MonitorRec
}

type LatestMonitorInfo struct {
	NodeNames map[int]string
	Websites  []WebsiteInfo
}
