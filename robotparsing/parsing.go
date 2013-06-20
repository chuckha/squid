package robotparsing

import (
	"net/url"
)

func GetRobotsTxtUrl(link string) string {
	u, _ := url.Parse(link)
	return u.Scheme + "://" + u.Host + "/robots.txt"
}
