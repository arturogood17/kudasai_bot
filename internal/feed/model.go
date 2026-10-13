// estructuras XML (título, link, categorías, fecha, media:content)

package feed

import (
	"encoding/xml"
	"time"
)

type Post struct {
	Title    string
	URL      string
	PubDate  time.Time
	Category string
	PhotoURL string
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel struct {
		Item []struct {
			Title    string `xml:"title"`
			Link     string `xml:"link"`
			PubDate  string `xml:"pubDate"`
			Category string `xml:"category"`
			Content  struct {
				URL string `xml:"url,attr"`
			} `xml:"content"`
		} `xml:"item"`
	} `xml:"channel"`
}
