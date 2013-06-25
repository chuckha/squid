// A package for parsing robots.txt
package robots

import (
	"bufio"
	"net/url"
	"net/http"
	"strings"
	"log"
	"io"
)

type RobotsTxt struct {
	DisallowAll, AllowAll bool
	// User-agents to disallowed URLs
	Rules map[string][]string
	Url *url.URL
	contents io.Reader
}

func NewRobotsTxtFromUrl(rawurl string) *RobotsTxt {
	parsedUrl, _ := url.Parse(rawurl)
	r := &RobotsTxt{
		Url: parsedUrl,
	}
	robotsUrl := GetRobotsTxtUrl(rawurl)
	r.GetRobotsTxtFromUrl(robotsUrl)
	r.ParseRobots()
	return r
}

// A helper method for testing.
// Supply a base URL and the contents of robots.txt.
func NewRobotsTxtFromText(rawurl, contents string) *RobotsTxt {
	parsedUrl, _ := url.Parse(rawurl)
	r := &RobotsTxt{
		Url: parsedUrl,
	}
	r.contents = strings.NewReader(contents)
	r.ParseRobots()
	return r
}

// Actually get the contents from some robots.txt url.
func (r *RobotsTxt) GetRobotsTxtFromUrl(robotsUrl string) {
	resp, err := http.Get(robotsUrl)
	if err != nil {
		log.Fatal(err)
	}
	r.contents = resp.Body
	resp.Body.Close()
	if resp.StatusCode > 310 {
		log.Println("Robots.txt not found")
		r.AllowAll = true
	}
}

// Build a map of UserAgents => Rules.
func (r *RobotsTxt) ParseRobots() {
	rules := make(map[string][]string)
	robots := bufio.NewScanner(r.contents)
	robots.Split(bufio.ScanLines)
	var currentUserAgent string
	for ; robots.Scan() ; {
		text := robots.Text()
		if strings.HasPrefix(text, "User-agent") {
			currentUserAgent = strings.TrimSpace(strings.Split(text, ":")[1])
			rules[currentUserAgent] = make([]string, 0)
		}
		if strings.HasPrefix(text, "Disallow") {
			if text == strings.TrimSpace("Disallow: /") && currentUserAgent == "*" {
				r.DisallowAll = true
			}
			path := strings.TrimSpace(strings.Split(text, ":")[1])
			if len(path) > 0 {
				rules[currentUserAgent] = append(rules[currentUserAgent], path)
			}
		}
	}
	r.Rules = rules
}

// Ask if a specific UserAgent and URL that it wants to crawl is an allowed action.
func (r *RobotsTxt) Allowed(ua, rawurl string) bool {
	parsedUrl, _ := url.Parse(rawurl)

	if r.DisallowAll { return false }
	if r.AllowAll { return true }

	userAgents := []string{ua, "*"}

	// TODO: Implement Allowed rules
	for _, ua := range userAgents {
		if _, ok := r.Rules[ua]; ok {
			for _, rule := range r.Rules[ua] {
				if rule == "/" {
					return false
				}
				if strings.HasPrefix(parsedUrl.Path, rule) {
					return false
				}
			}
		}
	}
	// No matched rule:
	return true
}

// GetRobotsTxtUrl returns the location of robots.txt given a URL
// that points to somewhere on the server.
func GetRobotsTxtUrl(rawurl string) string {
	u, _ := url.Parse(rawurl)
	return u.Scheme + "://" + u.Host + "/robots.txt"
}

