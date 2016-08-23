package main


import (
    "fmt"
    "sync"
    "regexp"
    "net/url"
    "encoding/json"

    "github.com/mvdan/xurls"

    "packages/url_title"
)

const mention_regex string = "\\B@(?P<value>\\w+)"
const emoticon_regex string = "\\((?P<value>\\w{1,15})\\)"

const fetch_timeout int = 1

type link struct {
    url string
    title string
}

func (s *link) MarshalJSON() ([]byte, error) {
    return []byte(fmt.Sprintf(`{"url":"%s","title":"%s"}`, s.url, s.title)), nil
}



func parse_links(message string) interface{} {
    links_chain := make(chan link)
    var wg sync.WaitGroup

    matched_urls := xurls.Relaxed.FindAllString(message, -1)

    for _, matched_url := range matched_urls {
        parsed_url, _ := url.Parse(matched_url)
        if parsed_url.Scheme == "" {
            parsed_url.Scheme = "http"
        }
        formated_url := parsed_url.String()

        wg.Add(1)
        go func() {
            defer wg.Done()
            lnk := link{formated_url, ""}
            title, ok := url_title.GetURLTitle(formated_url, fetch_timeout)
            if ok {
                lnk.title = title
            }
            links_chain <- lnk
        }()
    }

    links := make([]link, 0)

    for i:=0; i < len(matched_urls); i++ {
        link := <-links_chain
        links = append(links, link)
    }

    return links
}

func parse_regexp(message string,regex string) []string {
    r := regexp.MustCompile(regex)
    res := make([]string, 0)
    matches := r.FindAllStringSubmatch(message, -1)
    for _, match := range matches  {
        res = append(res, match[1])
    }
    return res
}

func parse_message(message string) {
    res := make(map[string]interface{})
    res["links"] = parse_links(message)
    res["mentions"] = parse_regexp(message, mention_regex)
    res["emoticons"] = parse_regexp(message, emoticon_regex)
    res_json, _ := json.Marshal(res)
    fmt.Println(string(res_json))
}


func main() {
    parse_message(`Hello world,
    @sarah @jane @parker
    yahoo.com
    http://abbac.co
    http://google.com
    http://i.beleive.does.not.exist.co
    google.com vk.com
    (test) (tes) (ttt)
    `)
}
