# hipchat-parser-go
hipchat-like parser for text, rewritten using golang

# usage example
<pre>
$: go run main.go "aaas dd @aa (aa) google.com"
</pre>
<pre>
Parsing:  aaas dd @aa (aa) google.com
{"emoticons":["aa"],"links":[{"url":"google.com","title":"Google"}],"mentions":["aa"]}
</pre>

# tests
run `go test`
