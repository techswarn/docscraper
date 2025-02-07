package main

import (
	//"fmt"
	"encoding/csv"
	"github.com/gocolly/colly/v2"
	"strings"
	"log"
	"os"
	"net/http"
	"fmt"

)
// links stores information about a digitalocean App platform docs

type Link struct {
	Title       string
	Description string
	URL         string
}

func main() {

	fName := "doclinks.csv"
	file , err := os.Create(fName)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	
	writer.Write([]string{"Instruction", "URL"})


	c := colly.NewCollector(
	    // Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("docs.digitalocean.com"),
	)
	// Links := make([]Link, 0, 500)
	// _ = Links

	//On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if !strings.HasPrefix(link, "/products/app-platform") {
			return
		}
		// start scraping the page under the link found
		e.Request.Visit(link)
	})

	c.OnHTML(`div[id=header-subheader]`, func(e *colly.HTMLElement) {
		log.Println("Doc found", e.Request.URL)
		resp, err := http.Get(fmt.Sprintf("%v",e.Request.URL))
		if err != nil {
			log.Fatal("Cannot get the page", err)
		}

		log.Printf("Response is %d", resp.StatusCode)

		if resp.StatusCode == 200 {
			title := strings.Split(e.ChildText("h1"), "\n")[0]
			log.Println(title)
			writer.Write([]string{title, e.Request.URL.String()})
		}
	})

	c.Visit("https://docs.digitalocean.com/products/app-platform")
}