package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	cfg    *Config
	engine *gin.Engine
	db     *sql.DB
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
		c.JSON(http.StatusOK, inf)
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

	s.engine.GET("/latest", s.getLatestMonitorInfo)

	return s
}
