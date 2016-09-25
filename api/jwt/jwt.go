package jwt

import (
	"time"
	"gopkg.in/dgrijalva/jwt-go.v3"
	"bitbucket.pearson.com/apseng/tensor/db"
	"bitbucket.pearson.com/apseng/tensor/models"
	"log"
	"gopkg.in/mgo.v2/bson"
	"errors"
)

type LocalToken struct {
	Token  string
	Expire string
}

func NewAuthToken() (*LocalToken, error) {
	// Initial middleware default setting.
	HeaderAuthMiddleware.MiddlewareInit()

	// Create the token
	token := jwt.New(jwt.GetSigningMethod(HeaderAuthMiddleware.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	collection := db.C(db.USERS)

	var admin models.User

	if err := collection.Find(bson.M{"username": "admin"}).One(&admin); err != nil {
		log.Println("User not found, Create JWT Token faild")
		return nil, errors.New("User not found, Create JWT Token faild")
	}

	expire := time.Now().Add(HeaderAuthMiddleware.Timeout)
	claims["id"] = admin.ID
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = time.Now().Unix()

	tokenString, err := token.SignedString(HeaderAuthMiddleware.Key)

	if err != nil {
		log.Println("Create JWT Token faild")
		return nil, errors.New("Create JWT Token faild")
	}

	return &LocalToken{Token:  tokenString, Expire: expire.Format(time.RFC3339), }, nil
}
