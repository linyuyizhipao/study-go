package web

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/url"
	"time"
)


func Run(){


	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello word")
	})
	r.GET("/welcome", func(c *gin.Context) {

		res,err :=Get("https://www.google.com/")
		if err!=nil{
			fmt.Println(res,222)
		}
		fmt.Println(111111,res)

	})
	//监听端口默认为8080
	r.Run(":8000")


}


var client *http.Client
// 发送GET请求
// url:请求地址
// response:请求返回的内容
func Get(urlPath string) (response string,err error) {

	ipAddress := "127.0.0.1:1080"
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://"+ipAddress)
	}

	tr := &http.Transport{
		Proxy: proxy,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client = &http.Client{
		Timeout: 5 * time.Second,
		Transport:tr,
	}

	resp, err := client.Get(urlPath)
	if err != nil {
		fmt.Println(12)
		return
	}
	defer func(){
		if err := resp.Body.Close(); err != nil{
			fmt.Printf("Get关闭 resp.Body.Close()发生了错误%s",err.Error())
		}
	}()

	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	n:=1
	for {
		n, err = resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			return
		}
	}
	response = result.String()
	return
}