package crawler

import (
	"testing"
)

type ParseTest struct {
	Input string
	ExpectedLinks, ExpectedKeywords []string
}

var parseTests = []ParseTest{
	{
		`<html>
			<head></head>
			<body>
				<h1>These are all keywords</h1>
				<a href="/this/is/an/internal/link">link</a>
				<p>this is a <a href="http://google.com">link</a></p>
			</body>
		</html>`,
		[]string{"/this/is/an/internal/link", "http://google.com"},
		[]string{"These", "are", "all", "keywords"},
	},
	{
		`<html><body><h1>Hello <a href="TEST">world</a></h1></body></html>`,
		[]string{"TEST"},
		[]string{"Hello"}, // This is dubious. See todo in crawler.go
	},
	{
		`<html><body><a href="/not/a/test"><h1>hi there</h1></a></body></html>`,
		[]string{"/not/a/test"},
		[]string{"hi", "there"},
	},
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		links, keywords := Parse(test.Input)
		for i := range links {
			if links[i] != test.ExpectedLinks[i] {
				t.Errorf("Expected to find:\n %v\nFound: %v\n", test.ExpectedLinks[i], links[i])
			}
		}
		for i := range keywords {
			if keywords[i] != test.ExpectedKeywords[i] {
				t.Errorf("Expected to find:\n %v\nFound: %v\n", test.ExpectedKeywords[i], keywords[i])
			}
		}
	}
}
