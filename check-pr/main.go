package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gocarina/gocsv"
)

func main() {
	bispr, e := getListPR()
	if e != nil {
		fmt.Println("csv 解析失败，请对照文档检查格式和转义有无问题。")
		os.Exit(1)
	}

	bis, e := getList()
	if e != nil {
		fmt.Println("自动检测失败，等待手动处理")
		os.Exit(1)
	}

	// 先检查是否有删除
	deleted := []string{}
LABEL1:
	for i := range bis {
		for j := range bispr {
			if bis[i].BizID == bispr[j].BizID {
				continue LABEL1
			}
		}
		deleted = append(deleted, bis[i].BizID)
	}

	if len(deleted) != 0 {
		fmt.Printf("删除了以下 bizid，如果是误删建议关闭本 pull request 重新提一个，如果确定要删除，等待手动处理\n\n%s\n", strings.Join(deleted, "\n"))
		os.Exit(1)
	}

	// 检查有无重复
	duplicated := []string{}
	m := map[string]int8{}
	for i := range bispr {
		_, ok := m[bispr[i].BizID]
		if ok {
			duplicated = append(duplicated, bispr[i].BizID)
		} else {
			m[bispr[i].BizID] = 0
		}
	}

	if len(duplicated) != 0 {
		fmt.Printf("以下 bizid 重复，建议修改或者关闭本 pull request 重新提一个\n\n%s\n", strings.Join(duplicated, "\n"))
		os.Exit(1)
	}

	fmt.Println("自动检测通过，等待手动处理")
	os.Exit(0)
}

type bizInfo struct {
	Name        string `csv:"name"`
	BizID       string `csv:"bizid"`
	Description string `csv:"description"`
}

func getList() ([]*bizInfo, error) {

	r, e := http.Get("https://github.com/hellodword/wechat-feeds/raw/main/list.csv")
	if e != nil {
		return nil, e
	}

	defer r.Body.Close()

	bis := []*bizInfo{}
	e = gocsv.Unmarshal(r.Body, &bis)
	if e != nil {
		return nil, e
	}

	return bis, nil
}

func getListPR() ([]*bizInfo, error) {

	b, e := ioutil.ReadFile("../list.csv")
	if e != nil {
		return nil, e
	}

	bis := []*bizInfo{}
	e = gocsv.Unmarshal(bytes.NewReader(b), &bis)
	if e != nil {
		return nil, e
	}

	return bis, nil
}
