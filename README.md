# go-textextract
A Golang library for extracting title and content from HTML
## Installation
The recommended way to install go-textextract

	go get github.com/tarylei/go-textextract
	
## Examples

	package main

	import (
		"fmt"

		"github.com/go-textextract/textextract"

		"github.com/ddliu/go-httpclient"
	)

	func main() {
		response, _ := httpclient.Get("http://wengengmiao.baijia.baidu.com/article/655838", nil)
		body, err := response.ToString()
		if err != nil {
			fmt.Println(err)
		}
		t := textextract.NewExtract(body)
		//当待抽取的网页正文中遇到成块的新闻标题未剔除时，只要增大此阈值即可
		//阈值增大，准确率提升，召回率下降；值变小，噪声会大，但可以保证抽到只有一句话的正文
		t.SetThreshold(86)
		fmt.Println(t.ExtractTitle())
		fmt.Println(t.ExtractText())
	}

