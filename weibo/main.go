//package main
//
//import (
//	"fmt"
//	"github.com/anaskhan96/soup"
//	"io/ioutil"
//	"log"
//	"net/http"
//)
//
//func main() {
//	requestUrl := "https://www.liaoxuefeng.com/"
//	// 发送Get请求
//	rsp, err := http.Get(requestUrl)
//	if err != nil {
//		log.Println(err.Error())
//		return
//	}
//	log.Println("1")
//	body, err := ioutil.ReadAll(rsp.Body)
//	if err != nil {
//		log.Println(err.Error())
//		return
//	}
//	log.Println("2")
//	content := string(body)
//	defer rsp.Body.Close()
//
//	// 下面主要是解析标签
//	doc := soup.HTMLParse(content)
//	links := doc.Find("div", "class", "uk-margin").FindAll("a")
//	log.Println(doc)
//	for _, link := range links {
//		fmt.Println(link.Text(), "| Link :", link.Attrs()["href"])
//	}
//	fmt.Println("\n\n")
//	for _, link := range links {
//		fmt.Println(link.Attrs()["href"])
//	}
//}

package main

import (
	"github.com/anaskhan96/soup"
	"io/ioutil"
	"log"
	"net/http"
)

type Link struct {
	Text string
	Href string
}

func main() {
	requestUrl := "https://tophub.today/n/KqndgxeLl9"
	// 发送Get请求
	rsp, err := http.Get(requestUrl)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("1")
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("2")
	content := string(body)
	defer rsp.Body.Close()

	slice := make([]Link, 0)

	// 下面主要是解析标签
	doc := soup.HTMLParse(content)
	links := doc.Find("div", "class", "cc-dc-c").FindAll("a")
	//log.Println(doc)
	for _, link := range links {
		//fmt.Println(link.Text(), "| Link :", link.Attrs()["href"])
		t := Link{Text: link.Text(), Href: "https://tophub.today" + link.Attrs()["href"]}
		slice = append(slice, t)
	}

	log.Println(slice[0])

}
