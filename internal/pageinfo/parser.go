package pageinfo

import (
	"bytes"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type Result struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

func ParseHTML(body []byte, finalURL string) Result {
	res := Result{URL: finalURL}
	root, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return res
	}

	var (
		titleText string
		descText  string
		iconHref  string
	)

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch strings.ToLower(n.Data) {
			case "title":
				if titleText == "" && n.FirstChild != nil {
					titleText = strings.TrimSpace(n.FirstChild.Data)
				}
			case "meta":
				name := attr(n, "name")
				prop := attr(n, "property")
				content := strings.TrimSpace(attr(n, "content"))
				if descText == "" && content != "" {
					if strings.EqualFold(name, "description") || strings.EqualFold(prop, "og:description") {
						descText = content
					}
				}
				if titleText == "" && content != "" {
					if strings.EqualFold(prop, "og:title") {
						titleText = content
					}
				}
			case "link":
				rel := strings.ToLower(attr(n, "rel"))
				href := strings.TrimSpace(attr(n, "href"))
				if iconHref == "" && href != "" {
					if strings.Contains(rel, "icon") {
						iconHref = href
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}

	walk(root)

	res.Title = titleText
	res.Description = descText
	res.Icon = resolveIcon(iconHref, finalURL)

	return res
}

func resolveIcon(href, pageURL string) string {
	if pageURL == "" {
		return href
	}
	parsed, err := url.Parse(pageURL)
	if err != nil {
		return href
	}
	if href == "" {
		return parsed.Scheme + "://" + parsed.Host + "/favicon.ico"
	}
	iconURL, err := url.Parse(href)
	if err != nil {
		return href
	}
	return parsed.ResolveReference(iconURL).String()
}

func attr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if strings.EqualFold(a.Key, key) {
			return a.Val
		}
	}
	return ""
}