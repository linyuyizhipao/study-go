package publishtask

import (
	"bytes"
	"crypto/rand"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/big"
	"time"
)

type UserService struct {
	User
}
func (us *UserService) createRandAccount()(name string){
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0;i < 9  ;i++  {
		randomInt,_ := rand.Int(rand.Reader,bigInt)
		container += string(str[randomInt.Int64()])
	}
	name = container
	return
}
func (us *UserService) createRandName()(name string){
	var container string
	var str = "afghijklmnopqrstuvwxyzABCDEFGHIJNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0;i < 10  ;i++  {
		randomInt,_ := rand.Int(rand.Reader,bigInt)
		container += string(str[randomInt.Int64()])
	}
	name = container
	return
}

func (us *UserService) createRandBalance()(balance int,err error){
	bigBalance,err := rand.Int(rand.Reader,big.NewInt(int64(100000)))
	if err != nil {
		return
	}
	balance = int(bigBalance.Int64())
	return
}
func (us *UserService) createRandEmail()(email string,err error){
	randEmailStr := "123456789"
	len := len(randEmailStr)
	for i:=0;i<8;i++{
		n,err :=rand.Int(rand.Reader,big.NewInt(int64(len)))
		if err != nil {
			return email,err
		}
		index := n.Int64()
		randStr :=randEmailStr[index:index+1]
		email += randStr
	}
	email += "163.com"
	return
}

func (us *UserService) getRandUser()(u User){
	user := new(User)
	user.Id = primitive.NewObjectID()
	user.Name = us.createRandName()
	user.Balance,_ = us.createRandBalance()
	user.CreateTime = time.Now().Unix()
	user.Email,_ = us.createRandEmail()
	user.Account = us.createRandAccount()
	user.PassWord = "123"
	u = *user
	return
}