package main

import "testing"

func TestParsing(t *testing.T) {

    messages := []string{
        "(aaa) (bbb) @sss",
        "http://google.com",
        "(toolongtobeemoticon) @mention",
        }

    results := []string{
        `{"emoticons":["aaa","bbb"],"mentions":["sss"]}`,
        `{"links":[{"url":"http://google.com","title":"Google"}]}`,
        `{"mentions":["mention"]}`,
    }
    for i, message := range messages {
        res := ParseMessage(message)
        if res != results[i] {
            t.Log(res, "!=", results[i])
            t.Fail()
        }
    }

}

