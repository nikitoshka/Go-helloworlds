package links

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

func forEachNode(node *html.Node, pre, post func(node *html.Node)) {
	if node == nil {
		return
	}

	if pre != nil {
		pre(node)
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		forEachNode(child, pre, post)
	}

	if post != nil {
		post(node)
	}
}

func Exract(url string) ([]string, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("Error: failed to get %s; status code %d",
			url, resp.StatusCode)
	}

	node, err := html.Parse(resp.Body)
	resp.Body.Close()

	if err != nil {
		return nil, err
	}

	var links []string

	var findLinks = func(node *html.Node) {
		if node == nil {
			return
		}

		if node.Type != html.ElementNode ||
			node.Data != "a" {
			return
		}

		for _, attr := range node.Attr {
			if attr.Key != "href" {
				continue
			}

			link, err := resp.Request.URL.Parse(attr.Val)
			if err != nil {
				continue
			}

			links = append(links, link.String())
		}
	}

	forEachNode(node, findLinks, nil)

	return links, nil
}
