package common

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type MySQLEnv struct {
	DatabaseName string        `env:"MYSQL_DATABASE_NAME" envDefault:"test"`
	User         string        `env:"MYSQL_USER" envDefault:"root"`
	Password     string        `env:"MYSQL_PASSWORD" envDefault:"root"`
	Host         string        `env:"MYSQL_URL_HOST" envDefault:"localhost:3306"`
	ConnTimeout  time.Duration `env:"MYSQL_CONN_TIMEOUT" envDefault:"10s"`
}

func NewMySQLConnection(env *MySQLEnv) (*gorm.DB, error) {
	var url = fmt.Sprintf(
		"%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		env.User, env.Password, env.Host, env.DatabaseName,
	)

	sess, err := gorm.Open("mysql", url)
	if err != nil {
		return nil, err
	}

	return sess, nil
}
