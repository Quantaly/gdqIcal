package lib

import "golang.org/x/net/html"

func findTag(n *html.Node, tagName string) *html.Node {
	if n.Type == html.ElementNode && n.Data == tagName {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		maybeRet := findTag(c, tagName)
		if maybeRet != nil {
			return maybeRet
		}
	}
	// couldn't find the tag
	return nil
}

func nextElement(n *html.Node) *html.Node {
	var ret *html.Node
	for ret = n.NextSibling; ret != nil && ret.Type != html.ElementNode; ret = ret.NextSibling {
	}
	return ret
}

func prevElement(n *html.Node) *html.Node {
	var ret *html.Node
	for ret = n.PrevSibling; ret != nil && ret.Type != html.ElementNode; ret = ret.PrevSibling {
	}
	return ret
}
