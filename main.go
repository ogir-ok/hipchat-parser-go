package main

import (
    "fmt"
    "log"
    "time"
    "regexp"
    "net/url"
    "encoding/json"

    "github.com/mvdan/xurls"

    "packages/url_title"
)

const mentionRegex string = "\\B@(?P<mention>\\w+)"
const emoticonRegex string = "\\((?P<emoticon>\\w{1,15})\\)"

// Timeout to fetch title of link in message
const fetchTimeout time.Duration = time.Duration(2) * time.Second


type Link struct {
    url string
    title string
}

func (s *Link) MarshalJSON() ([]byte, error) {
    return []byte(fmt.Sprintf(`{"url":"%s","title":"%s"}`, s.url, s.title)), nil
}

// Finds links in text and gives slice of Link struct objects
func parseLinks(message string) []Link {
    linksChain := make(chan Link)

    matchedURLs := xurls.Relaxed.FindAllString(message, -1)

    for _, matchedURL := range matchedURLs {
        parsedURL, _ := url.Parse(matchedURL)
        if parsedURL.Scheme == "" {
            parsedURL.Scheme = "http"
        }
        formatedURL := parsedURL.String()

        // Async fetch of titles in urls (all at the same time)
        go func() {
            link := Link{formatedURL, ""}

            titleChan := make(chan string, 1)
            go func() {
                title, _ := url_title.GetURLTitle(formatedURL)
                titleChan <- title
            }()

            select {
            case title := <- titleChan:
                link.title = title
            case <-time.After(fetchTimeout):
                log.Print("Timeout fetch for url ", formatedURL)
            }

            linksChain <- link
        }()
    }

    links := make([]Link, 0)

    for i:=0; i < len(matchedURLs); i++ {
        link := <-linksChain
        links = append(links, link)
    }

    return links
}

// Parses message for given regex. Regex should have submatch
func parseRegexp(message string,regex string) []string {
    r := regexp.MustCompile(regex)
    res := make([]string, 0)
    matches := r.FindAllStringSubmatch(message, -1)
    for _, match := range matches  {
        res = append(res, match[1])
    }
    return res
}


// Parses message. Returns json with links, emoticons, mentions
func ParseMessage(message string) string {
    res := make(map[string]interface{})
    links := parseLinks(message)
    mentions := parseRegexp(message, mentionRegex)
    emoticons := parseRegexp(message, emoticonRegex)

    if len(links) > 0 {
        res["links"] = links
    }

    if len(emoticons) > 0 {
        res["emoticons"] = emoticons
    }

    if len(mentions) > 0 {
        res["mentions"] = mentions
    }

    resJson, _ := json.Marshal(res)
    return string(resJson)
}


func main() {
    messages := []string{
        "(aaa) (bbb) @sss",
        "http://vk.com",
        "(toolongtobeemoticon) @mention",
        `Hello world,
        @sarah @jane @parker
        yahoo.com
        http://google.com
        http://i.beleive.does.not.exist.co
        http://eloquentjavascript.net/Eloquent_JavaScript.pdf
        http://www.axmag.com/download/pdfurl-guide.pdf
        google.com vk.com
        (test) (tes) (ttt)
        `,}

    for _, message := range messages {
        fmt.Println("Parsing: ", message)
        res := ParseMessage(message)
        fmt.Println(res)
    }
}
