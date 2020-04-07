package database

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type MySQLEnv struct {
	MySQLDatabaseName string        `env:"MYSQL_DATABASE_NAME" envDefault:"courses"`
	MySQLUser         string        `env:"MYSQL_USER" envDefault:"jossie"`
	MySQLPassword     string        `env:"MYSQL_PASSWORD" envDefault:"jossie"`
	MySQLHost         string        `env:"MYSQL_URL_HOST" envDefault:"localhost:3306"`
	MySQLConnTimeout  time.Duration `env:"MYSQL_CONN_TIMEOUT" envDefault:"10s"`
}

type MySQLProvider interface {
	DBName() string
	DBSession() (*gorm.DB, error)
	DBClose()
}

type mySQLProvider struct {
	env     *MySQLEnv
	session *gorm.DB
}

func (p *mySQLProvider) DBName() string {
	return p.env.MySQLDatabaseName
}

func (p *mySQLProvider) DBSession() (*gorm.DB, error) {
	var url = fmt.Sprintf(
		"%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		p.env.MySQLUser, p.env.MySQLPassword, p.env.MySQLHost, p.env.MySQLDatabaseName,
	)

	var err error

	p.session, err = gorm.Open("mysql", url)
	if err != nil {
		return nil, err
	}

	return p.session, nil
}

func (p *mySQLProvider) DBClose() {
	p.session.Close()
}

func NewMySQLProvider(env *MySQLEnv) *mySQLProvider {
	return &mySQLProvider{env: env}
}
