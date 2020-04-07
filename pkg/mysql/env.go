package mysql

import "time"

type Env struct {
	DatabaseName string        `env:"MYSQL_DATABASE_NAME" envDefault:"test"`
	User         string        `env:"MYSQL_USER" envDefault:"root"`
	Password     string        `env:"MYSQL_PASSWORD" envDefault:"root"`
	Host         string        `env:"MYSQL_URL_HOST" envDefault:"localhost:3306"`
	ConnTimeout  time.Duration `env:"MYSQL_CONN_TIMEOUT" envDefault:"10s"`
}
