package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"hugo/mon"
	"strconv"
	"sync"
	"time"
)

type Abs struct {
	UserName string `bson:"username"`
	Email string `bson:"email"`
}

func main(){
	switch1()
	return

	abs := new(Abs)
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(time.Now().Second())

	singResult :=mon.Client.Collection("tbl_data").FindOne(context.TODO(),bson.M{"username":"胡工827"})
	singResult.Decode(abs)
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(time.Now().Second())
	return

	wg := new(sync.WaitGroup)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go insertMon(wg)
	}
	wg.Wait()

	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))

}

func insertMon(wg *sync.WaitGroup){

	defer func() {
		wg.Done()
	}()

	data := Abs{"hugo","136586551@163.com"}

	for i := 1; i < 100000; i++ {
		iStr := strconv.Itoa(i)
		data.UserName = fmt.Sprintf("胡工%s",iStr)
		insertRes,err:=mon.Client.Collection("tbl_data").InsertOne(context.TODO(),data)
		if err != nil {
			fmt.Println(err,insertRes)
		}
	}

}

func doSth(){
	//ctx,cancel :=context.WithTimeout(context.Background(),time.Duration(150)* time.Microsecond)
	timer := time.NewTimer(2*time.Second)

	//go sth(ctx)

	for range timer.C{
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	}

	return
	for {
		select{
		case <-timer.C:
			fmt.Println("定时器时间到了，开始要取消cancel()\n")
			//cancel()
			return
		default:
			//fmt.Println("还没有开始123")

		}
	}

}

func sth(ctx context.Context){
	fmt.Println(1111)
	for{
		time.Sleep(time.Microsecond * 100)
		select {
		case <-ctx.Done():
			fmt.Println("取消结束了")
			return
		default:
			fmt.Println("还没有被结束567")
		}
	}
}

//定时器每隔2秒执行一次代码块
func timer1(){
	timer := time.NewTimer(2*time.Second)
	for range timer.C{
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	}
}

//定时器隔2秒执行一次代码块,仅仅执行一次
func timer2(){
	timer := time.NewTimer(2*time.Second)
	i :=0
	for true {
		select{
		case <-timer.C:
			fmt.Println("定时器执行了")
		default:
			if i == 0 {
				fmt.Println("定时器还没执行了333333")
			}
		}
		i++
	}
}

//context 执行
//for 里面的break  会直接终止for循坏
//select 里面的case里面的break是终止case里面的代码块执行，不会终止到select外层的for
//select 里面不用break也只会执行一个case
func context1(){
	ctx,cancel :=context.WithTimeout(context.Background(),2 * time.Second)
	i:=0

	go func(ctx context.Context) {
		for{

			select{
				case <-time.After(10 * time.Nanosecond):
					fmt.Println("时间1")
					fmt.Println("78787878")

				case <-ctx.Done():
					fmt.Println("取消的命令来了")
					return
				default:
					if i == 0 {
						fmt.Println("11122333")
					}
					break
			}
			i++
		}
	}(ctx)

	time.Sleep(1 * time.Second)
	fmt.Println("cancel来取消的")
	cancel()

}

//执行到的switch的case如果有 fallthrough  就会在不判断的情况下直接执行下一个case
func switch1(){
	df :=1
	switch df{
	case 1:
		fmt.Println(11)
		fallthrough
	case 2,3,4:
		fmt.Println(222)
		fallthrough
	case 22,32,42:
		fmt.Println(8989)
	case 222,322,422:
		fmt.Println(84378434)

	}
}

//后面sync.atomic




