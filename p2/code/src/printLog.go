package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func openTxt(txt string) string {
	filePath := txt
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("file Open failed = ", err)
		return ""
	}
	defer file.Close()              // 关闭文本流
	reader := bufio.NewReader(file) // 读取文本数据
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		fmt.Print(str)
	}

	return ""
}

func main() {

	//output the result
	fmt.Println("===================For test1======================")
	openTxt("./test1.txt")
	fmt.Printf("\n \n \n")

	fmt.Println("===================For test2======================")
	openTxt("./test2.txt")
	fmt.Printf("\n \n \n")

}
