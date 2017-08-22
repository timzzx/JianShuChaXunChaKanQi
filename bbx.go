// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/lxn/walk"

	. "github.com/lxn/walk/declarative"
)

func main() {

	//main
	mw := &MyMainWindow{}

	if err := (MainWindow{
		Icon:     "img/x.ico",
		AssignTo: &mw.MainWindow,
		Title:    "简书mini -- by tim_zhang",
		MenuItems: []MenuItem{
			Menu{
				Text: "&编辑",
				Items: []MenuItem{
					Separator{},
					Action{
						Text:        "退出",
						OnTriggered: func() { mw.Close() },
					},
				},
			},
			Menu{
				Text: "&帮助",
				Items: []MenuItem{
					Action{
						Text:        "关于",
						OnTriggered: mw.aboutAction_Triggered,
					},
				},
			},
		},
		MinSize: Size{1000, 600},

		Layout: VBox{MarginsZero: true},
		Children: []Widget{
			Composite{
				MaxSize: Size{0, 50},
				Layout:  HBox{},
				Children: []Widget{
					Label{Text: "关键词: "},
					LineEdit{
						AssignTo: &mw.keywords,
						Text:     "golang",
					},
					PushButton{
						AssignTo: &mw.query,
						Text:     "查询",
					},
					PushButton{
						AssignTo: &mw.nextquery,
						Text:     "下一页",
					},
					PushButton{
						AssignTo: &mw.prvquery,
						Text:     "上一页",
					},
				},
			},

			Composite{
				Layout: Grid{Columns: 2, Spacing: 10},
				Children: []Widget{
					ListBox{
						MaxSize:               Size{200, 0},
						AssignTo:              &mw.lb,
						OnCurrentIndexChanged: mw.lb_CurrentIndexChanged,
						OnItemActivated:       mw.lb_ItemActivated,
					},

					WebView{

						AssignTo: &mw.wv,
					},
				},
			},
		},
	}.Create()); err != nil {
		log.Fatal(err)
	}

	mw.query.Clicked().Attach(func() {
		go func() {
			mw.query.SetText("查询中...")
			mw.query.SetEnabled(false)
			mw.GetList()
			mw.query.SetText("查询")
			mw.query.SetEnabled(true)

		}()
	})

	mw.nextquery.Clicked().Attach(func() {
		go func() {
			mw.page = mw.page + 1
			mw.GetNextList(mw.page)
		}()
	})

	mw.prvquery.Clicked().Attach(func() {
		go func() {
			mw.page = mw.page - 1
			if mw.page == 0 {
				mw.page = 1
			}
			mw.GetNextList(mw.page)
		}()
	})

	mw.Run()
}

type MyMainWindow struct {
	*walk.MainWindow
	lb        *walk.ListBox
	te        *walk.TextEdit
	wv        *walk.WebView
	keywords  *walk.LineEdit
	model     *Model
	query     *walk.PushButton
	nextquery *walk.PushButton
	prvquery  *walk.PushButton
	page      int
}
type Item struct {
	name  string
	value string
}

type Model struct {
	walk.ListModelBase
	items []Item
}

func (m *Model) ItemCount() int {
	return len(m.items)
}

func (m *Model) Value(index int) interface{} {
	return m.items[index].name
}

func (mw *MyMainWindow) aboutAction_Triggered() {
	walk.MsgBox(mw, "关于", "mini版简书阅读器v1.0\n作者：tim_zhang", walk.MsgBoxIconQuestion)
}

func (mw *MyMainWindow) GetList() {
	keywords := mw.keywords.Text()
	if len(keywords) <= 0 {
		walk.MsgBox(mw, "查询地址", "请填写关键词", walk.MsgBoxIconWarning)
		return
	}
	enkeywords := url.QueryEscape(keywords)
	body := httpDo("GET", "http://www.jianshu.com/search/do?q="+enkeywords+"&type=note&page=1&order_by=default")
	mw.page = 1
	var dat map[string]interface{}
	msg := "查询成功"
	json.Unmarshal([]byte(body), &dat)
	if v, ok := dat["entries"]; ok {
		entries := v.([]interface{})
		m := &Model{items: make([]Item, len(entries))}
		re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
		for i, all := range entries {
			entries1 := all.(map[string]interface{})
			title, _ := entries1["title"]
			slug, _ := entries1["slug"]
			m.items[i] = Item{re.ReplaceAllString(title.(string), ""), re.ReplaceAllString(slug.(string), "")}
		}

		mw.lb.SetModel(m)
		mw.model = m
		return
	} else {
		msg = "查询请10秒一次"
		walk.MsgBox(mw, "查询", msg, walk.MsgBoxIconInformation)
	}

	return
}

func (mw *MyMainWindow) GetNextList(page int) {
	keywords := mw.keywords.Text()

	if len(keywords) <= 0 {
		walk.MsgBox(mw, "查询地址", "请填写关键词", walk.MsgBoxIconWarning)
		return
	}
	enkeywords := url.QueryEscape(keywords)
	body := httpDo("GET", "http://www.jianshu.com/search/do?q="+enkeywords+"&type=note&page="+strconv.Itoa(page)+"&order_by=default")
	var dat map[string]interface{}
	msg := "查询成功"
	json.Unmarshal([]byte(body), &dat)
	if v, ok := dat["entries"]; ok {
		entries := v.([]interface{})
		m := &Model{items: make([]Item, len(entries))}
		re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
		for i, all := range entries {
			entries1 := all.(map[string]interface{})
			title, _ := entries1["title"]
			slug, _ := entries1["slug"]
			m.items[i] = Item{re.ReplaceAllString(title.(string), ""), re.ReplaceAllString(slug.(string), "")}
		}

		mw.lb.SetModel(m)
		mw.model = m
		return
	} else {
		msg = "查询请10秒一次"
		walk.MsgBox(mw, "查询", msg, walk.MsgBoxIconInformation)
	}

	return
}

func (mw *MyMainWindow) lb_CurrentIndexChanged() {

	go func() {
		i := mw.lb.CurrentIndex()
		defer mw.lb.SetCurrentIndex(-1)
		item := &mw.model.items[i]
		doc, err := goquery.NewDocument("http://www.jianshu.com/p/" + item.value)
		if err != nil {
			//log.Fatal(err)
		}
		doc.Find(".show-content").Each(func(i int, s *goquery.Selection) {
			re, _ := regexp.Compile("<img[\\S\\s]+?/>")
			userFile := "data/xxx.html"
			fout, err := os.Create(userFile)
			defer fout.Close()
			if err != nil {
				fmt.Println(userFile, err)
				return
			}
			ht, _ := s.Html()
			fout.WriteString("<meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\" /><h3>" + item.name + "</h3><a href='http://www.jianshu.com/p/" + item.value + "' target='_blank'>文章地址</a>" + re.ReplaceAllString(ht, ""))

			mw.wv.SetURL("file:///" + getCurrentDirectory() + "/data/xxx.html")

		})
		doc.Find(".normal-comment-list").Each(func(i int, s *goquery.Selection) {
			ht, _ := s.Html()
			appendToFile("data/xxx.html", "<h3>评论</h3>"+ht)
		})
	}()
}

func (mw *MyMainWindow) lb_ItemActivated() {
	value := mw.model.items[mw.lb.CurrentIndex()].value
	fmt.Println(value)
}

func appendToFile(fileName string, content string) error {
	// 以只写的模式，打开文件
	f, err := os.OpenFile(fileName, os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("cacheFileList.yml file create failed. err: " + err.Error())
	} else {
		// 查找文件末尾的偏移量
		n, _ := f.Seek(0, os.SEEK_END)
		// 从末尾的偏移量开始写入内容
		_, err = f.WriteAt([]byte(content), n)
	}
	defer f.Close()
	return err
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func httpDo(method string, url string) (bodys string) {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {

	}

	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req.Header.Set("Accept", "application/json")
	}
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")
	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	bodys = string(body)

	if err != nil {
		// handle error
		return ""
	}

	// fmt.Println(string(body))
	return bodys
}
