package publishtask

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id primitive.ObjectID `bson:"_id"`
	Account string `bson:"account"`
	PassWord string `bson:"pass_word"`
	Email string `bson:"email"`
	Name string `bson:"name"`
	Balance int `bson:"balance"`
	CreateTime int64 `bson:"create_time"`
}


