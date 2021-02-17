package main

import (
	"bytes"
	"context"
	"errors"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

var ErrLicenseNotFound = errors.New("no license found")

type License string

var predefinedLicenses = map[string]License{
	"flutter":             License("BSD"),
	"flutter_web_plugins": License("BSD"),
	"flutter_test":        License("BSD"),
	"sky_engine":          License("BSD"),
}

func getLicense(ctx context.Context, name string) (License, error) {
	if l, ok := predefinedLicenses[name]; ok {
		return l, nil
	}
	req, err := http.NewRequestWithContext(ctx, "GET", "https://pub.dev/packages/"+name, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "github.com/plan3t-one/flutter_licenses")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	node, err := html.Parse(resp.Body)
	if err != nil {
		return "", err
	}

	infoBox := getElementByClass(node, "detail-info-box")

	if infoBox == nil {
		return "", ErrLicenseNotFound
	}

	for c := infoBox.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.FirstChild.Data == "License" {
			// the next sibling of the infoBox is our license tag
			text := &bytes.Buffer{}
			collectText(c.NextSibling.NextSibling, text)
			str := text.String()

			str = strings.TrimSuffix(str, " (LICENSE)")

			if str != "" {
				return License(str), nil
			} else {
				return "", ErrLicenseNotFound
			}
		}
	}

	return "", ErrLicenseNotFound
}

func getAttribute(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}

func checkClass(n *html.Node, id string) bool {
	if n.Type == html.ElementNode {
		s, ok := getAttribute(n, "class")
		if ok && s == id {
			return true
		}
	}
	return false
}

func traverse(n *html.Node, class string) *html.Node {
	if checkClass(n, class) {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result := traverse(c, class)
		if result != nil {
			return result
		}
	}

	return nil
}

func getElementByClass(n *html.Node, class string) *html.Node {
	return traverse(n, class)
}

func collectText(n *html.Node, buf *bytes.Buffer) {
	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectText(c, buf)
	}
}
