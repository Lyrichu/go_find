package main

import (
	"flag"
	"fmt"
	"github.com/Lyrichu/go_find/internal/finder"
	"os"
)

func main() {
	// 定义命令行参数
	name := flag.String("name", "", "正则表达式匹配文件名")
	fileType := flag.String("type", "f", "搜索类型 (f: 文件, d: 目录)")
	flag.Parse()

	// 获取搜索路径
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("请指定搜索路径")
		flag.Usage()
		os.Exit(1)
	}
	path := args[0]

	// 创建 finder 实例
	f := finder.NewFinder(*name, *fileType)

	// 执行搜索
	results, err := f.Find(path)
	if err != nil {
		fmt.Printf("搜索出错: %v\n", err)
		os.Exit(1)
	}

	// 输出结果
	for _, result := range results {
		formattedResult := f.FormatResult(result)
		fmt.Printf("%s\n", formattedResult)
	}
}