// Package server TODO
package server

import (
	_ "github.com/go-sql-driver/mysql" // mysql TODO
	"github.com/qianlnk/codidoc/config"
	"github.com/xormplus/xorm"
)

// Server TODO
type Server struct {
	cfg *config.Config

	mysqlCli *xorm.Engine
}

// New TODO
func New(cfg *config.Config) (*Server, error) {
	engine, err := xorm.NewDB("mysql", cfg.MysqlDSN)
	if err != nil {
		return nil, err
	}

	return &Server{
		cfg:      cfg,
		mysqlCli: engine,
	}, nil
}

// Run TODO
func (s *Server) Run() {
	go s.Download()
	go s.PushLoop()
}
