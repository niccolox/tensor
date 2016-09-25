package users

import (
	"gopkg.in/mgo.v2/bson"
	"bitbucket.pearson.com/apseng/tensor/models"
	"time"
	"bitbucket.pearson.com/apseng/tensor/db"
	"log"
)

func addActivity(crdID bson.ObjectId, userID bson.ObjectId, desc string) {

	a := models.Activity{
		ID: bson.NewObjectId(),
		ActorID: userID,
		Type: "user",
		ObjectID: crdID,
		Description: desc,
		Created: time.Now(),
	}

	c := db.C(db.ACTIVITY_STREAM)

	if err := c.Insert(a); err != nil {
		log.Println("Failed to add new Activity", err)
	}
}

