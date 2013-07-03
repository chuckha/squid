// A package for parsing robots.txt
package robots

import (
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Rules map[string][]string

func (r Rules) Add(key, value string) {
	r[key] = append(r[key], value)
}

type RobotsTxt struct {
	DisallowAll, AllowAll bool
	// User-agents to disallowed URLs
	Rules    Rules
	Url      *url.URL
	contents string
}

func NewRobotsTxtFromUrl(rawurl string) *RobotsTxt {
	parsedUrl, _ := url.Parse(rawurl)
	r := &RobotsTxt{
		Url: parsedUrl,
	}
	robotsUrl := GetRobotsTxtUrl(rawurl)
	r.GetRobotsTxtFromUrl(robotsUrl)
	r.Rules = GetRules(r.contents)
	return r
}

// Actually get the contents from some robots.txt url.
func (r *RobotsTxt) GetRobotsTxtFromUrl(robotsUrl string) {
	resp, err := http.Get(robotsUrl)
	if err != nil {
		log.Fatal(err)
	}
	// What have we learened? We must read it and close it!
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	r.contents = string(body)
	if resp.StatusCode > 310 {
		log.Printf("Robots.txt not found for url %v", robotsUrl)
		r.AllowAll = true
	}
}

func GetRules(contents string) Rules {
	rules := make(Rules)
	robotsText := bufio.NewScanner(strings.NewReader(contents))
	var currentUserAgent string
	for robotsText.Scan() {
		text := CleanInput(robotsText.Text())
		// Ignore comments
		if strings.HasPrefix(text, "#") {
			continue
		} else if strings.HasPrefix(text, "user-agent") {
			currentUserAgent = strings.TrimSpace(strings.Split(text, ":")[1])
		} else if strings.HasPrefix(text, "disallow") {
			path := strings.TrimSpace(strings.Split(text, ":")[1])
			rules.Add(currentUserAgent, path)
		}
	}
	return rules
}

func CleanInput(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// Ask if a specific UserAgent and URL that it wants to crawl is an allowed action.
// BUG(ChuckHa): Will fail when UserAgent: * and Disallow: / followed by UserAgent: Squidbot and Disallow:
func (r *RobotsTxt) Allowed(ua, rawurl string) bool {
	ua = CleanInput(ua)
	parsedUrl, _ := url.Parse(rawurl)

	// Check specific user agents first
	userAgents := []string{"*", ua}

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

func (r *RobotsTxt) NotAllowed(ua, rawurl string) bool {
	return !r.Allowed(ua, rawurl)
}

// GetRobotsTxtUrl returns the location of robots.txt given a URL
// that points to somewhere on the server.
func GetRobotsTxtUrl(rawurl string) string {
	u, _ := url.Parse(rawurl)
	return u.Scheme + "://" + u.Host + "/robots.txt"
}
