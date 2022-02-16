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
	fmt.Println("===============AliceLog===============")
	openTxt("./testEntry/test1/log/Alice.txt")
	fmt.Println("===============BobLog===============")
	openTxt("./testEntry/test1/log/Bob.txt")
	fmt.Println("===============ChadLog===============")
	openTxt("./testEntry/test1/log/Chad.txt")
	fmt.Printf("\n \n \n")

	fmt.Println("===================For test2======================")
	//output the result
	fmt.Println("===============AliceLog===============")
	openTxt("./testEntry/test2/log/Alice.txt")
	fmt.Println("===============BobLog===============")
	openTxt("./testEntry/test2/log/Bob.txt")
	fmt.Println("===============ChadLog===============")
	openTxt("./testEntry/test2/log/Chad.txt")
	fmt.Println("===============DoughLog===============")
	openTxt("./testEntry/test2/log/Dough.txt")
	fmt.Printf("\n \n \n")

	fmt.Println("===================For test3======================")
	//output the result
	fmt.Println("===============AliceLog===============")
	openTxt("./testEntry/test3/log/Alice.txt")
	fmt.Println("===============BobLog===============")
	openTxt("./testEntry/test3/log/Bob.txt")
	fmt.Println("===============ChadLog===============")
	openTxt("./testEntry/test3/log/Chad.txt")
	fmt.Printf("\n \n \n")

}
