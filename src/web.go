package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

var wg sync.WaitGroup

type SiteMapIndex struct {
	Locations []string `xml:"sitemap>loc"`
}

type News struct {
	Titles    []string `xml:"url>news>title"`
	Keywords  []string `xml:"url>news>keywords"`
	Locations []string `xml:"url>loc"`
}

type NewsMap struct {
	Keyword  string
	Location string
}

type NewsAggPage struct {
	Title string
	News  map[string]NewsMap
}

func main() {
	http.HandleFunc("/", index_handler)
	http.HandleFunc("/agg/", news_agg_handler)
	http.ListenAndServe(":8000", nil)
}

func index_handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>News Aggregator!!</h1>\n")
	fmt.Fprintf(w, "<p1>News aggregator for washingtonpost</p1>\n")
}

func news_routine(c chan News, loc string){
	defer wg.Done()
	var n News
	p := strings.Trim(loc, "\n")
	resp1, _ := http.Get(p)
	bytes1, _ := ioutil.ReadAll(resp1.Body)
	xml.Unmarshal(bytes1, &n)
	resp1.Body.Close()
	fmt.Println("Added:", p)
	c <- n
}

func news_agg_handler(w http.ResponseWriter, r *http.Request) {
	var s SiteMapIndex
	News_map := make(map[string]NewsMap)
	queue := make(chan News, 100)
	resp, _ := http.Get("https://www.washingtonpost.com/news-sitemaps/index.xml")
	bytes, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(bytes, &s)
	resp.Body.Close()

	for _, val := range s.Locations {
		wg.Add(1)
		go news_routine(queue, val)
	}
	wg.Wait()
	close(queue)
	for ele := range queue {
		for ind, _ := range ele.Keywords {
			News_map[ele.Titles[ind]] = NewsMap{ele.Keywords[ind], ele.Locations[ind]}
		}
	}

	page := NewsAggPage{Title: "Sample title", News: News_map}
	t, err := template.ParseFiles("sampletemplate.html")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t.Execute(w, page))
}
