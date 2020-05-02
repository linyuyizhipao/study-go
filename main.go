package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"
)

type Goods struct {
	Id primitive.ObjectID `bson:"_id"`
	Name string `bson:"name"`
	Price int `bson:"price"`
	Inventory int `bson:"inventory"`
	CreateTime int64 `bson:"create_time"`
}

type Order struct {
	Id primitive.ObjectID `bson:"_id"`
	Uid primitive.ObjectID `bson:"uid"`
	GoodsId primitive.ObjectID `bson:"goods_id"`
	Price int `bson:"price"`
	CreateTime int64 `bson:"create_time"`
}
type User struct {
	Id primitive.ObjectID `bson:"_id"`
	UserName string `bson:"user_name"`
	Balance int `bson:"balance"`
	CreateTime int64 `bson:"create_time"`
}

type UserGoodsByOrderRel struct {
	GoodsId primitive.ObjectID
	Uid primitive.ObjectID
}
var userGoodsByOrderRelChan = make(chan UserGoodsByOrderRel,5000)


func main() {
	//1.创建5000个用户并发发生并发抢一个商品
	userBuyGoods()
	//2.消费队列处理生成有效订单
	generateOrder()

	//3.websocket通知消息
	inform()

	fmt.Println("结束完毕")
}

func userBuyGoods(){

	//生成5000个用户
	user := new (User)
	goods := new (Goods)
	userGoodsByOrderRel := new(UserGoodsByOrderRel)
	userModel :=Mon.Collection("tbl_user")
	goodsModel :=Mon.Collection("tbl_goods")
	ordersModel :=Mon.Collection("tbl_order")

	if _,err := userModel.DeleteMany(context.TODO(),bson.M{});err!=nil{
		fmt.Println(err)
	}
	if _,err := goodsModel.DeleteMany(context.TODO(),bson.M{});err!=nil{
		fmt.Println(err)
	}
	if _,err := ordersModel.DeleteMany(context.TODO(),bson.M{});err!=nil{
		fmt.Println(err)
	}

	for i:=0;i<5000;i++{
		user.Balance = 1000
		user.CreateTime = time.Now().Unix()
		user.Id = primitive.NewObjectID()
		user.UserName = "hugo"
		insertId,err :=userModel.InsertOne(context.TODO(),user)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(insertId,"用户生成成功")
	}

	//生成一个商品
	goods.Id = primitive.NewObjectID()
	goods.CreateTime = time.Now().Unix()
	goods.Name = "短袖"
	goods.Price = 128
	goods.Inventory = 500
	goodsId,err :=goodsModel.InsertOne(context.TODO(),goods)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(goodsId,"商品生成成功")

	filter := bson.M{"user_name":"hugo"}
	findOptios :=options.Find()
	findOptios.SetLimit(5000)
	cur,err :=userModel.Find(context.TODO(),filter,findOptios)
	if err != nil {
		fmt.Println(err)
	}

	for cur.Next(context.TODO()){
		err :=cur.Decode(user)
		if err != nil {
			fmt.Println(err)
		}
		userGoodsByOrderRel.Uid = user.Id
		userGoodsByOrderRel.GoodsId = goods.Id
		userGoodsQueue := *userGoodsByOrderRel
		userGoodsByOrderRelChan <- userGoodsQueue
	}
	if err :=cur.Close(context.TODO());err!=nil{
		fmt.Println(err)
	}

}
func generateOrder(){
	order := new(Order)
	goodsInfo := new(Goods)
	user := new(User)
	filter := bson.M{"name":"短袖"}
	if err :=Mon.Collection("tbl_goods").FindOne(context.TODO(),filter).Decode(goodsInfo);err!=nil{
		fmt.Println(err)
	}
	inventory := goodsInfo.Inventory

	defer func() {
		close(userGoodsByOrderRelChan)
	}()

	for userGoodsByOrder := range userGoodsByOrderRelChan{
		fmt.Println("已经有人在抢了")
		if inventory <= 0 {
			fmt.Println("商品抢完了")
			return
		}
		findFilter := bson.M{"_id":userGoodsByOrder.Uid}
		//检查用户是否合法
		if err :=Mon.Collection("tbl_user").FindOne(context.TODO(),findFilter).Decode(user);err!=nil || user.Balance < 0 || user.Balance < goodsInfo.Price {
			fmt.Println(err,"余额：",user.Balance)
			continue
		}

		order.Id = primitive.NewObjectID()
		order.GoodsId = userGoodsByOrder.GoodsId
		order.Uid = userGoodsByOrder.Uid
		order.Price = goodsInfo.Price
		order.CreateTime = time.Now().Unix()
		_,err :=Mon.Collection("tbl_order").InsertOne(context.TODO(),order)
		if err != nil {
			fmt.Println(err)
		}
		filter :=bson.M{"_id":userGoodsByOrder.Uid}
		update := bson.D{
			{"$inc", bson.D{
				{"age", -goodsInfo.Price},
			}},
		}
		if _,err :=Mon.Collection("tbl_user").UpdateOne(context.TODO(),filter,update);err!=nil{
			fmt.Println(err)
		}

		inventory--

	}
}
func inform(){

}




func init(){
	Dbconnect()
}

var Mon *mongo.Database
func GetContext() (ctx context.Context) {
	ctx, _ = context.WithTimeout(context.TODO(), 10*time.Second)
	return
}
func Dbconnect() {
	want, err := readpref.New(readpref.SecondaryMode) //表示只使用辅助节点
	if err != nil {
		fmt.Println(err)
	}
	wc := writeconcern.New(writeconcern.WMajority())
	readconcern.Majority()
	//链接mongo服务
	opt := options.Client().ApplyURI("mongodb://localhost:27017")
	opt.SetLocalThreshold(3 * time.Second)     //只使用与mongo操作耗时小于3秒的
	opt.SetMaxConnIdleTime(5 * time.Second)    //指定连接可以保持空闲的最大毫秒数
	opt.SetMaxPoolSize(10)                    //使用最大的连接数
	opt.SetReadPreference(want)                //表示只使用辅助节点
	opt.SetReadConcern(readconcern.Majority()) //指定查询应返回实例的最新数据确认为，已写入副本集中的大多数成员
	opt.SetWriteConcern(wc)                    //请求确认写操作传播到大多数mongod实例
	ctx := GetContext()

	if database, err := mongo.Connect(ctx, opt); err != nil {
		fmt.Println(err)
	}else{
		//UseSession(client)
		//判断服务是否可用
		if err = database.Ping(ctx, readpref.Primary()); err != nil {
			fmt.Println(err)
		}

		Mon = database.Database("testing_base")
	}
}
