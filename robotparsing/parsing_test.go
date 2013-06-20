package robotparsing

import (
	"testing"
)

type RobotsTxtTest struct {
	in, out string
}

var tests = []RobotsTxtTest{
	{
		"http://google.com",
		"http://google.com/robots.txt",
	},
	{
		"https://google.com",
		"https://google.com/robots.txt",
	},
	{
		"https://google.com/hello/world",
		"https://google.com/robots.txt",
	},
	{
		"https://google.com/file%20one%20two",
		"https://google.com/robots.txt",
	},
	{
		"http://www.google.com/",
		"http://www.google.com/robots.txt",
	},
}

func TestGetRobotsTxtUrl(t *testing.T) {
	// return a robots.txt link given any valid url
	for _, test := range tests {
		robourl := GetRobotsTxtUrl(test.in)
		if robourl != test.out {
			t.Errorf("We did not put robots in the correct spot: %s\n", robourl)
		}
	}
}
