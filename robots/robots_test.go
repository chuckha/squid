package robots

import (
	"testing"
)

const (
	userAgent   = "Squidbot"
	robotsText1 = "User-agent: *\nDisallow: /no_bots.txt\nDisallow: /filez\nDisallow: /a_really/really/really/really/deep/world.html\n"
	robotsText2 = "User-agent: *\nDisallow:\n"
	robotsText3 = "User-agent: GoogleBot\nDisallow: /\n"
	robotsText4 = "User-agent: Squidbot\nDisallow: /\n"
	robotsText5 = "User-agent: *\nDisallow: /forums/post.txt\n"
	robotsText6 = "User-agent: *\nDisallow: /help\n"
	robotsText7 = "User-agent: *\nDisallow: /help/\n"
)

type RobotTxtTest struct {
	UserAgent string
	Robotstxt string
	Url       string
	Expected  bool
}

var tests = []RobotTxtTest{
	{
		userAgent,
		robotsText7,
		"/help/index.html",
		false,
	},
	{
		userAgent,
		robotsText7,
		"/help.txt",
		true,
	},
	{
		userAgent,
		robotsText6,
		"/help/index.html",
		false,
	},
	{
		userAgent,
		robotsText6,
		"/help.txt",
		false,
	},
	{
		userAgent,
		robotsText5,
		"/forums",
		true,
	},
	{
		userAgent,
		robotsText4,
		"/yunolikeme",
		false,
	},
	{
		userAgent,
		robotsText3,
		"/not_googlebot_so_i/enjoy/crawling/this/site",
		true,
	},
	{
		userAgent,
		robotsText2,
		"/everything/will/be/true",
		true,
	},
	{
		userAgent,
		robotsText2,
		"/",
		true,
	},
	{
		userAgent,
		robotsText1,
		"http://google.com/no_bots.txt",
		false,
	},
	{
		userAgent,
		robotsText1,
		"/no_bots.txt",
		false,
	},
	{
		userAgent,
		robotsText1,
		"/filez",
		false,
	},
	{
		userAgent,
		robotsText1,
		"/a_really/really/really/really/deep/world.html",
		false,
	},
	{
		userAgent,
		robotsText1,
		"/something_else",
		true,
	},
	{
		userAgent,
		robotsText1,
		"/a_really/really/really/really/deep/file.txt",
		true,
	},
}

func TestSimpleRobotsTxt(t *testing.T) {
	url := "http://google.com/robots.txt"
	for _, test := range tests {
		robotsTxt := NewRobotsTxtFromText(url, test.Robotstxt)
		allowed := robotsTxt.Allowed(test.UserAgent, test.Url)
		val := allowed == test.Expected
		if !val {
			t.Errorf("Expecting: %v, got: %v on test: %v", test.Expected, allowed, test)
		}
	}
}

type RobotsTxtUrlTests struct {
	Given, Expected string
}

var robotsTxtUrlTests = []RobotsTxtUrlTests {
	{
		"http://google.com",
		"http://google.com/robots.txt",
	},
	{
		"https://duckduckgo.com/",
		"https://duckduckgo.com/robots.txt",
	},
	{
		"http://www.chuckha.com/blog/Introduction-to-python/",
		"http://www.chuckha.com/robots.txt",
	},
}

func TestGetRobotsTxtUrl(t *testing.T) {
	for _, test := range robotsTxtUrlTests {
		url := GetRobotsTxtUrl(test.Given)
		if test.Expected != url {
			t.Errorf("Expected: %v Got: %v", test.Expected, url)
		}
	}
}


