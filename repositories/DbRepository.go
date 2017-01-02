package repositories

import (
	"time"

	"github.com/iftekhersunny/shortx/configs"

	mgo "gopkg.in/mgo.v2"
)

// Database repository struct
type DbRepository struct {
	//
}

// Init method
func (dbRepository DbRepository) init() (*mgo.Database, *mgo.Session) {

	dialInfo := &mgo.DialInfo{
		Addrs:    []string{configs.DB_HOST + ":" + configs.DB_PORT},
		Timeout:  60 * time.Second,
		Database: configs.DB_NAME,
		Username: configs.DB_USER,
		Password: configs.DB_PASSWORD,
	}

	session, err := mgo.DialWithInfo(dialInfo)

	if err != nil {
		panic("failed to connect database")
	}

	return session.DB(configs.DB_NAME), session
}
