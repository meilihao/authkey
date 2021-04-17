package db

import (
	"fmt"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/meilihao/layer"
)

type DBConfig struct {
	Host         string
	Port         int
	Name         string
	Username     string
	Password     string
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	ShowSQL      bool   `yaml:"show_sql"`
	Loc          string `yaml:"loc"`
}

func InitMySQL2Layer(conf *DBConfig) (*layer.Layer, error) {
	if conf.Loc == "" {
		conf.Loc = url.QueryEscape("Asia/Shanghai")
	}

	l, err := layer.New(
		layer.WithDB("mysql", fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=%s`,
			conf.Username,
			conf.Password,
			conf.Host,
			conf.Port,
			conf.Name,
			conf.Loc)),
		layer.WithConnMaxLifetime(7*60*60*time.Second),
		layer.WithMaxIdleConns(20),
		layer.WithMaxOpenConns(100),
		//layer.WithDryRun(true),
		layer.WithDebug(true),
	)
	if err != nil {
		return nil, err
	}

	return l, nil
}
