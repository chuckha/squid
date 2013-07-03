package crawler

import (
	"fmt"
	"net/http"
	"net/url"
	"code.google.com/p/go.net/html"
	"github.com/ChuckHa/Squid/db"
	"strings"
	"log"
	"io/ioutil"
)

const (
	userAgent = "Squidbot"
)

type Metadata struct {
	Site string
	Links, Keywords []string
}

func (m *Metadata) Save() error {
	c := db.GetCollection()
	defer c.Database.Session.Close()
	return c.Insert(m)
}

type Crawler struct {
	client *http.Client
	userAgent, rawurl string
}

func NewCrawler(userAgent, rawurl string) *Crawler {
	client := &http.Client{}
	return &Crawler{
		client: client,
		userAgent: userAgent,
		rawurl: rawurl,
	}
}

// Get the gatekeeper (robots.txt) from the site.
// Get the HTML and parse it.
func (c *Crawler) Crawl() (*Metadata, error) {
	resolver, _ := url.Parse(c.rawurl)
	cleanUrl := resolver.Scheme + "://" + resolver.Host + "/" + resolver.Path

	md := &Metadata{Site: cleanUrl}

	if db.Exists(cleanUrl) {
		return md, fmt.Errorf("Already visited")
	}

	content, err := c.GetHTML()
	if err != nil {
		return md, err
	}
	links, keywords := Parse(content)
	resolvedLinks := make([]string, len(links))
	for i, link := range links {
		if strings.HasPrefix(link, "http") {
			resolvedLinks[i] = link
			continue
		}
		ref, _ := url.Parse(link)
		resolvedLinks[i] = resolver.ResolveReference(ref).String()
	}
	md.Links = resolvedLinks
	md.Keywords = keywords
	return md, err
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
			// don't get things like mailto: and javascript:void(0)
			if strings.HasPrefix(a.Val, "http://") {
				return a.Val
			}
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
