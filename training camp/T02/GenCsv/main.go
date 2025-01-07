package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

type Link struct {
	Id      int
	Title   string
	WebLink string
}

func main() {
	// 创建一个 tutorials.csv 文件
	csvFile, err := os.Create("tutorials.csv")
	if err != nil {
		panic(err)
	}
	fmt.Println(os.Getwd())
	defer csvFile.Close()
	// 设置随机数种子
	// 设置随机种子

	tutorials := []Link{}
	// 初始化字典数据
	for i := 0; i < 3000; i++ {
		//rand.Seed(time.Now().UnixNano())

		// 生成 [9999, 10000] 范围的随机整数
		min_, max_ := 100000, 999999
		randomNum := rand.Intn(max_-min_+1) + min_
		T := strconv.Itoa(randomNum)
		if i%7 == 0 {
			T = "baidu.com"
		}
		tutorials = append(tutorials, Link{
			Id:      i + 1,
			Title:   "tt" + strconv.Itoa(i),
			WebLink: "http://" + T,
		})

	}

	// 初始化一个 csv writer，并通过这个 writer 写入数据到 csv 文件
	csvFile.WriteString("\xEF\xBB\xBF")

	writer := csv.NewWriter(csvFile)
	for _, tutorial := range tutorials {
		line := []string{
			strconv.Itoa(tutorial.Id), // 将 int 类型数据转化为字符串
			tutorial.Title,
			tutorial.WebLink,
		}
		// 将切片类型行数据写入 csv 文件
		err := writer.Write(line)
		if err != nil {
			panic(err)
		}
	}
	// 将 writer 缓冲中的数据都推送到 csv 文件，至此就完成了数据写入到 csv 文件
	writer.Flush()

	// 打开这个 csv 文件
	file, err := os.Open("tutorials.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

}
