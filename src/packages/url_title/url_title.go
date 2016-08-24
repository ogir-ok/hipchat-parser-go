// contains utility to get url's title
package url_title

import (
    "io"
    "net/http"
    "golang.org/x/net/html"
)

// https://siongui.github.io/2016/05/10/go-get-html-title-via-net-html/. 
// probably there is a better way to parse
func isTitleElement(n *html.Node) bool {
    return n.Type == html.ElementNode && n.Data == "title"
}

func traverse(n *html.Node) (string, bool) {
    if isTitleElement(n) {
        return n.FirstChild.Data, true
    }

    for c := n.FirstChild; c != nil; c = c.NextSibling {
        result, ok := traverse(c)
        if ok {
            return result, ok
        }
    }

    return "", false
}

func getHtmlTitle(r io.Reader) (string, bool) {
    doc, err := html.Parse(r)
    if err != nil {
        panic("Fail to parse html")
    }

    return traverse(doc)
}

func GetURLTitle(url string) (string, bool) {
    http.Get(url)
    res, err := http.Get(url)
    if err == nil {
        defer res.Body.Close()
        title, ok := getHtmlTitle(res.Body)
        if ok {
            return title, ok
        }
    }
    return "", false
}
