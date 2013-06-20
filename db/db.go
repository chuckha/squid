package db

import (
	"labix.org/v2/mgo"
	"log"
)

const (
	mongoUrl   = "localhost"
	dbName     = "squid"
	collection = "websites"
)

var (
	session *mgo.Session
)

func GetSession() *mgo.Session {
	if session == nil {
		session, err := mgo.Dial(mongoUrl)
		if err != nil {
			log.Printf("Error connecting to mongo: %s", err)
		}
	}
	return session.Clone()
}

func GetCollection() *mgo.Collection {
	session := GetSession()
	return session.DB("").C(collection)
}
