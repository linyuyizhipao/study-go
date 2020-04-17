package file

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"time"
)

//将整个文件一次性读取到内存  使用  ioutil 包实现
func ReadFileContent(filePath string) (content string, err error) {
	data, rErr := ioutil.ReadFile(filePath)
	if rErr != nil {
		fmt.Println("file read err", rErr)
		err = rErr
		return
	}
	content = string(data)
	return
}

//分块读文件内容，利用到了缓存 bufio 包
func ChunkReadFile(filePath string) {
	chunkSize := 3 //每次读取文件内容的大小
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err = f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	r := bufio.NewReader(f) //通过文件资源获取一个缓存读取器，后面迭代这个读取器就能套出文件内容
	b := make([]byte, chunkSize)

	for {
		_, err := r.Read(b)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Printf("本次读到的块内容:%s,\n", string(b))
	}
}

//逐行读取文件内容
func readFileByLine(filePath string) {
	r, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		if err := r.Close(); err != nil {
			fmt.Println(err)
			return
		}
	}()

	s :=bufio.NewScanner(r)

	for s.Scan(){
		fmt.Printf("每行的内容:%s,\n",s.Text())
	}

	if err := s.Err(); err != nil {
		fmt.Printf("s逐行读取发生了错误:%s",err.Error())
	}
}


//将字符串一行一行的追加写入一个新的文件，若文件原本存在则覆盖原有的后，一行行的追加我们自己的内容
func writeFile(){
	f ,err := os.Create("a.txt")
	defer func(){
		if err := f.Close();err!=nil{
			fmt.Println(err)
		}
	}()

	if err != nil{
		fmt.Println(err)
		return
	}
	strArr := make([]string,2)
	strArr = append(strArr,"2dfdfdfdf")
	strArr = append(strArr,"3gggg")
	for _,str := range strArr {
		_,err :=fmt.Fprintln(f,str)
		if err != nil {
			fmt.Println(err.Error(),2323233232)
		}
	}
	fmt.Println("写入成功")
}

func appendWriteFile(){
	f, err := os.OpenFile("a.txt", os.O_APPEND|os.O_WRONLY, 0644)
	defer func(){
		if err:=f.Close();err != nil {
			fmt.Println(err)
		}
	}()


	if err != nil {
		fmt.Println(err)
	}
	strs :=[]string{"a撒大声地","fdfdfdf"}

	for _,v := range strs {
		fmt.Fprintln(f,v)
	}
	fmt.Println("追加成功")
}

//并发写入文件
//其实对接文件的写入还是只有一个
func concurrentWriteFile(){
	//0:初始化变量
	productionChan := make(chan string,100)
	wg :=sync.WaitGroup{}
	//1.先起10个生产者
	for i:=1;i<=10;i++{
		wg.Add(1)
		go func(){
			//每个work 就不停的往目标channel里面写东西，间隔200毫秒，写1000个位置
			//这样就很容易产生了，10台机器每台机器1000的写入量，并发的写入channel，
			for j :=1; j < 1000;j++{
				time.Sleep(time.Microsecond * 20)
				jstr := strconv.Itoa(j)
				productionChan <- "hugo" + jstr
			}
			wg.Done()
		}()
	}

	//wg.Wait() 检测生产者生产 是否完毕，完毕则关闭通道 productionChan,后面的range productionChan 就会收到信号停止range
	go func(){
		wg.Wait()
		close(productionChan)
	}()

	f, err := os.OpenFile("a.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err,222)
	}

	defer func(){
		if err:=f.Close();err != nil {
			fmt.Println(err)
		}
	}()

	//写入文件的消费者只能有一个
	for v := range productionChan {
		if _,err :=fmt.Fprintln(f,v);err!=nil{
			fmt.Println(err,111)
		}
	}



}
