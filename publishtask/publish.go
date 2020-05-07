package publishtask

import (
	"context"
	"fmt"
	"hugo/mon"
)

//创建100万用户
//每个用户随机1000到100000余额
func CreateUser() {
	userModel := mon.Client.Collection("tbl_user")
	userService := new(UserService)
	for i := 1; i <= 1000000; i++ {
		user := userService.getRandUser()
		if _, err := userModel.InsertOne(context.TODO(), &user); err != nil {
			fmt.Println(err)
		}
	}
}

//用户创建一个任务
func createTask() {}

//根据任务信息产生具体任务行为
func publishTask() {}

//用户收到推送任务的通知
func receiveTask() {}

//用户领取任务
func getTask() {}

//用户完成任务
func finishTask() {}
