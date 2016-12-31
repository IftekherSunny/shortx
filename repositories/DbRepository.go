package repositories

import (
	"github.com/iftekhersunny/shortx/configs"

	mgo "gopkg.in/mgo.v2"
)

// Database repository struct
type DbRepository struct {
	//
}

// Init method
func (dbRepository DbRepository) init() (*mgo.Database, *mgo.Session) {

	session, err := mgo.Dial(configs.DB_HOST)

	if err != nil {
		panic("failed to connect database")
	}

	return session.DB(configs.DB_NAME), session
}
