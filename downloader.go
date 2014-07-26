package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {
	crawler := mockCrawler{}
	start(crawler)
}

func start(crawler crawler) {
	fmt.Printf("%#v\n", crawler)

	var initialLinks = crawler.findLinkFromInitialPage()

	fmt.Printf("%#v\n", initialLinks)

	done := make(chan bool)

	for _, link := range initialLinks {
		go findVideoLink(link, done)
	}

	for i := 0; i < len(initialLinks); i++ {
		<-done
	}

	fmt.Println("everything completed")
}

type crawler interface {
	findLinkFromInitialPage() []string
}

type realCrawler struct {
}

func (realCrawler) findLinkFromInitialPage() []string {
	var doc *goquery.Document
	var e error

	// Initialize doc, (in real use, the method would be goquery.NewDocument)
	if doc, e = goquery.NewDocument("http://somesite.net"); e != nil {
		log.Fatal(e)
	}

	links := make([]string, 0)

	// Find the review items (the type of the Selection would be *goquery.Selection)
	doc.Find("div.post-content a[title]").Each(func(i int, s *goquery.Selection) {
		var a, _ = s.Attr("href")
		links = append(links, a)
	})

	return links
}

type mockCrawler struct {
}

func (mockCrawler) findLinkFromInitialPage() []string {
	return []string{"http://www.somesite.net/?pasd=13123"}
}

func findVideoLink(link string, done chan<- bool) {

	res, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)

	res.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	str := string(body)

	r, _ := regexp.Compile("file\":\"http:.*mp4")

	result := r.FindAllString(str, -1)

	if len(result) > 0 {
		var first = result[0]
		fmt.Printf("%#v\n", first)

		first = strings.Replace(first, "file\":\"", "", -1)
		first = strings.Replace(first, "\"", "", -1)
		fmt.Printf("%#v\n", first)
		downloadVidFromLink(first, done)
	} else {
		done <- true
	}
}

func downloadVidFromLink(downloadLink string, done chan<- bool) {
	res, err := http.Get(downloadLink)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var filename = getFileName(downloadLink)

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		out, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()

		io.Copy(out, res.Body)

		fmt.Println("saved file ", filename)
	} else {
		fmt.Println(filename, " already exist")
	}

	done <- true
}

func getFileName(downloadLink string) string {
	hasher := md5.New()
	io.WriteString(hasher, downloadLink)

	var hashStr = hex.EncodeToString(hasher.Sum(nil)) + ".mp4"

	fmt.Printf("%#v", hashStr)
	return hashStr
}
