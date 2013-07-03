package db

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
)

const (
	mongoUrl   = "localhost"
	dbName     = "squid"
	collection = "websites"
)

var (
	session *mgo.Session
	err error
)

func GetSession() *mgo.Session {
	if session == nil {
		session, err = mgo.Dial(mongoUrl)
		if err != nil {
			log.Printf("Error connecting to mongo: %s", err)
		}
	}
	return session.Clone()
}

func GetCollection() *mgo.Collection {
	session := GetSession()
	return session.DB(dbName).C(collection)
}

func Exists(url string) bool {
	C := GetCollection()
	defer C.Database.Session.Close()
	count, err := C.Find(bson.M{"site": url}).Count()
	if err != nil {
		log.Printf("Error talking to mongo: %v", err)
	}
	return count >= 1
}
