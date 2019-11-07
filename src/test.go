package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/xml"
)


type News struct {
	Titles [] string `xml:"url>news>title"`
	Keywords [] string `xml:"url>news>keywords"`
	Locations [] string `xml:"url>loc"`
}


func main() {
	var n News
	p := "https://www.washingtonpost.com/news-sitemaps/politics.xml"
	resp1, _ := http.Get(p)
	fmt.Println(resp1.Body)
	bytes1, _ := ioutil.ReadAll(resp1.Body)
	xml.Unmarshal(bytes1, &n) 
	for _, v := range n.Keywords {
		fmt.Println(v)
	}
	resp1.Body.Close()
}