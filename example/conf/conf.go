package conf

import "github.com/arstd/log"

type conf struct {
	Name string

	DB struct {
		Dialect  string
		Host     string
		Port     int
		Username string
		Password string
		DBName   string
		Params   string
	}
}

var Conf conf

func init() {
	Conf.DB.Dialect = "mysql"
	Conf.DB.Host = "139.224.15.11"
	Conf.DB.Port = 3306
	Conf.DB.Username = "fireeyes"
	Conf.DB.Password = "fy3fsc4ptv"
	Conf.DB.DBName = "ac_fireeyes"
	Conf.DB.Params = "charset=utf8&parseTime=true&loc=Local"

	log.Json(Conf)
}
