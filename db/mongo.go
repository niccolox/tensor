package db

import (
	"gopkg.in/mgo.v2"
	"pearson.com/hilbert-space/util"
	"time"
)

var MongoDb *mgo.Database


// Mongodb database
func Connect() error {

	cfg := util.Config.MongoDB

	// We need this object to establish a session to our MongoDB.
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    cfg.Hosts,
		Timeout:  60 * time.Second,
		Database: cfg.DbName,
		Username: cfg.Username,
		Password: cfg.Password,
		ReplicaSetName: cfg.ReplicaSet,
		Mechanism: "SCRAM-SHA-1",
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		return err
	}

	// Switch the session to a monotonic behavior.
	//session.SetMode(mgo.Monotonic, true)

	if err := session.Ping(); err != nil {
		return err
	}

	MongoDb = session.DB(cfg.DbName)
	return nil
}