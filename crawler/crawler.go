package crawler

import (
	"fmt"
	"net/url"
	"net/http"
	"code.google.com/p/go.net/html"
	"github.com/ChuckHa/Squid/robots"
	"strings"
	"log"
	"io/ioutil"
)

const (
	userAgent = "Squidbot"
)

type Crawler struct {
	client *http.Client
	userAgent, rawurl string
	cURL chan *url.URL
	links, keywords []string
}

func NewCrawler(userAgent, rawurl string, cURL chan *url.URL) *Crawler {
	client := &http.Client{}
	return &Crawler{
		client: client,
		userAgent: userAgent,
		rawurl: rawurl,
		cURL: cURL,
	}
}

// Get the gatekeeper (robots.txt) from the site.
// Get the HTML and parse it.
// Save result to mongo.
func (c *Crawler) Crawl() error {
	gatekeeper := robots.NewRobotsTxtFromUrl(c.rawurl)
	if gatekeeper.NotAllowed(c.userAgent, c.rawurl) {
		return fmt.Errorf("Disallowed website from robots.txt")
	}
	content, err := c.GetHTML()
	if err != nil {
		return err
	}
	links, keywords := Parse(content)
	// TODO: Save to mongo
	//err = mongo.Save()
	log.Println(links, keywords)
	return err
}

// Attach our user agent to the headers and make the request.
func (c *Crawler) GetHTML() (string, error) {
	req, err := http.NewRequest("GET", c.rawurl, nil)
	if err != nil {
		log.Println("Error creating Request: %v", err)
	}
	req.Header.Add("User-Agent", c.userAgent)

	resp, err := c.client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Printf("%s -- %s", c.rawurl, err)
		return "", err
	}
	contents, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode > 310 {
		statusText := http.StatusText(resp.StatusCode)
		log.Printf("%v: %v on %v", resp.StatusCode, statusText, c.rawurl)
		return "", fmt.Errorf(statusText)
	}
	return string(contents), nil
}

func getLink(n *html.Node) string {
	for _, a := range n.Attr {
		if a.Key == "href" {
			return a.Val
		}
	}
	// blah
	return ""
}

func Compact(words []string) []string {
	compact := make([]string, 0)
	for _, word := range words {
		w := strings.TrimSpace(word)
		if len(w) == 0 {
			continue
		}
		compact = append(compact, w)
	}
	return compact
}

// TODO: Implement a removal of common words
// TODO: Implement stemming
// TODO: Get deeper than first child words
func getKeywords(n *html.Node) []string {
	header := n.FirstChild.Data
	words := strings.Split(header, " ")
	return Compact(words)
}

func Parse(content string) ([]string, []string) {
	links := make([]string, 0)
	keywords := make([]string, 0)
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		log.Fatalln(err)
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		// if the node is an Element and it's an attribute
		if n.Type == html.ElementNode {
			switch n.Data {
			case "a":
				links = append(links, getLink(n))
			case "h1":
				keywords = append(keywords, getKeywords(n)...)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return links, keywords
}
