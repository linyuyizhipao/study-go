package file

import (
	"testing"
)

//将整个文件一次性读取到内存  使用  ioutil 包实现
func TestReadFilecontent(t *testing.T){
	content,_ :=ReadFileContent("./config.txt")
	t.Log(content,"读取到了")
}
//将整个文件一次性读取到内存  使用  ioutil 包实现
func TestChunkReadFile(t *testing.T){
	ChunkReadFile("./config.txt")
}

//将整个文件一次性读取到内存  使用  ioutil 包实现
func TestReadFileByLine(t *testing.T){
	readFileByLine("./config.txt")
}
//写入测试
func TestWriteFile(t *testing.T){
	writeFile()
}

func TestAppendFile(t *testing.T){
	appendWriteFile()
}

func TestConcurrentWriteFile(t *testing.T){
	concurrentWriteFile()
}