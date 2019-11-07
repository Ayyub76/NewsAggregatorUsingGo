package main

import (
	"fmt"
	"net/http"
	"html/template"
	"io/ioutil"
	"encoding/xml"
	"strings"
)


type SiteMapIndex struct {
	Locations [] string `xml:"sitemap>loc"`
}

type News struct {
	Titles [] string `xml:"url>news>title"`
	Keywords [] string `xml:"url>news>keywords"`
	Locations [] string `xml:"url>loc"`
}

type NewsMap struct {
	Keyword string
	Location string
}

type NewsAggPage struct {
	Title string
	News map[string]NewsMap
}


func main() {
	http.HandleFunc("/", index_handler)
	http.HandleFunc("/agg/", news_agg_handler)
	http.ListenAndServe(":8000", nil)
}

func index_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Test Web!!</h1>\n")
	fmt.Fprintf(w, "<p1>This is a test webpage made using golang!!!</p1>\n")
}

func news_agg_handler(w http.ResponseWriter, r *http.Request) {
	var s SiteMapIndex
	var n News
	var p string
	News_map := make(map[string]NewsMap)
	resp, _ := http.Get("https://www.washingtonpost.com/news-sitemaps/index.xml")
	bytes, _ := ioutil.ReadAll(resp.Body)
	
	xml.Unmarshal(bytes, &s)

	//fmt.Println(s)

	for _, val := range s.Locations {
		p = strings.Trim(val, "\n")
		resp1, _ := http.Get(p)
		bytes1, _ := ioutil.ReadAll(resp1.Body)
		xml.Unmarshal(bytes1, &n)

		for ind, _ := range n.Keywords {
			News_map[n.Titles[ind]] = NewsMap{n.Keywords[ind], n.Locations[ind]}
			//fmt.Println(ind1, n.Titles[ind])
		}
		resp1.Body.Close()
	}
	resp.Body.Close()

	page := NewsAggPage{Title: "Sample title", News: News_map}
	t, err := template.ParseFiles("sampletemplate.html")
	if err == nil {
		fmt.Println(t.Execute(w, page))
    } else {
    	fmt.Println(err)
    }
}