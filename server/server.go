package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
	internal "github.com/z4yx/tunaworks/internal"
)

type Server struct {
	cfg    *Config
	engine *gin.Engine
	db     *sql.DB
	recent sync.Map
}

type empty struct {
}

func contextErrorLogger(c *gin.Context) {
	errs := c.Errors.ByType(gin.ErrorTypeAny)
	if len(errs) > 0 {
		for _, err := range errs {
			logger.Errorf(`"in request "%s %s: %s"`,
				c.Request.Method, c.Request.URL.Path,
				err.Error())
		}
	}
	// pass on to the next middleware in chain
	c.Next()
}

func (s *Server) getLatestMonitorInfo(c *gin.Context) {
	inf, err := s.QueryLatestMonitorInfo()
	if err != nil {
		logger.Errorf("getLatestMonitorInfo: %s", err.Error())
		c.JSON(http.StatusInternalServerError, empty{})
	} else {
		for node, info := range inf.NodeInfo {
			t, ok := s.recent.Load(node)
			// logger.Debug("%d %v %v", node ,ok, t)
			if ok {
				info.Heartbeat = t.(time.Time)
			} else {
				info.Heartbeat = time.Unix(0, 0)
			}
			inf.NodeInfo[node] = info
		}
		c.JSON(http.StatusOK, inf)
	}
}

func (s *Server) getAllWebsites(c *gin.Context) {
	token := c.GetHeader("X-Token")
	if proto, valid := c.GetQuery("Proto"); valid && token != "" {
		if proto_i, err := strconv.Atoi(proto); err == nil {
			s.UpdateNodeProtocol(token, proto_i)
		}
	}
	inf, err := s.QuerySites(true)
	if err != nil {
		logger.Errorf("getAllWebsites: %s", err.Error())
		c.JSON(http.StatusInternalServerError, empty{})
	} else {
		c.JSON(http.StatusOK, inf)
	}
}

func (s *Server) insertProberRecord(c *gin.Context) {
	auth, node := s.AuthNode(c.GetHeader("X-Token"))
	if !auth {
		c.JSON(http.StatusForbidden, empty{})
	} else {
		var result internal.ProbeResult
		// raw, _ := c.GetRawData()
		// logger.Debug("%v", raw)
		if err := c.ShouldBindJSON(&result); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		s.recent.Store(node, time.Now())
		err := s.InsertRecord(node, &result)
		if err != nil {
			c.JSON(http.StatusInternalServerError, empty{})
		} else {
			c.JSON(http.StatusOK, empty{})
		}
	}
}

func (s *Server) Run() {
	db, err := sql.Open(s.cfg.Server.DBProvider, s.cfg.Server.DBName)
	if err != nil {
		panic(err)
	}
	s.db = db

	addr := fmt.Sprintf("%s:%d", s.cfg.Server.Addr, s.cfg.Server.Port)

	httpServer := &http.Server{
		Addr:         addr,
		Handler:      s.engine,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}

}

func MakeServer(cfg *Config) *Server {

	s := &Server{
		engine: gin.New(),
		cfg:    cfg,
	}
	s.engine.Use(gin.Recovery())
	if cfg.Debug {
		s.engine.Use(gin.Logger())
	}
	s.engine.Use(contextErrorLogger)

	s.engine.Static("/assets", "./assets")
	s.engine.StaticFile("/", "./assets/html/index.html")
	s.engine.StaticFile("/ssl", "./assets/html/ssl.html")
	s.engine.GET("/monitor/latest", s.getLatestMonitorInfo)
	s.engine.GET("/prober/websites", s.getAllWebsites)
	s.engine.POST("/prober/result", s.insertProberRecord)

	return s
}
