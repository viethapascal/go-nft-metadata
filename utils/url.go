package utils

import (
	url2 "net/url"
	"path/filepath"
)

func BuildUrl(prefix string, path ...string) string {
	u1, _ := url2.Parse(prefix)
	url := url2.URL{
		Scheme: "https",
		Host:   u1.Host,
		Path:   filepath.Join(path...),
	}
	return url.String()
}
