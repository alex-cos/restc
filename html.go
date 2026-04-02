package restc

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func getBodyConcatainedText(data []byte) (string, error) {
	doc, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrParseHTML, err)
	}
	body := findNodeByTagName(doc, "body")
	if body == nil {
		return "", fmt.Errorf("%w: %w",
			ErrParseHTML,
			errors.New("can't find body node in HTML response"),
		)
	}
	return concatAllText(body, ": "), nil
}

func findNodeByTagName(node *html.Node, name string) *html.Node {
	if node == nil {
		return nil
	}
	if node.Type == html.ElementNode && strings.EqualFold(node.Data, name) {
		return node
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		found := findNodeByTagName(child, name)
		if found != nil {
			return found
		}
	}

	return nil
}

func concatAllText(node *html.Node, sep string) string {
	if node == nil {
		return ""
	}
	if node.Type == html.ElementNode && strings.EqualFold(node.Data, "script") {
		return ""
	}
	if node.FirstChild == nil && node.Type == html.TextNode {
		return node.Data
	}
	texts := []string{}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		txt := concatAllText(child, sep)
		if strings.TrimSpace(txt) != "" {
			texts = append(texts, txt)
		}
	}

	return strings.Join(texts, sep)
}
