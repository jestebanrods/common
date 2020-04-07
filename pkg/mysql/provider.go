package mysql

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Provider interface {
	DBName() string
	DBSession() (*gorm.DB, error)
	DBClose()
}

type provider struct {
	env     *Env
	session *gorm.DB
}

func (p *provider) DBName() string {
	return p.env.DatabaseName
}

func (p *provider) DBSession() (*gorm.DB, error) {
	var url = fmt.Sprintf(
		"%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		p.env.User, p.env.Password, p.env.Host, p.env.DatabaseName,
	)

	sess, err := gorm.Open("mysql", url)
	if err != nil {
		return nil, err
	}

	p.session = sess

	return sess, nil
}

func (p *provider) DBClose() {
	p.session.Close()
}

func NewProvider(env *Env) *provider {
	return &provider{env: env}
}
