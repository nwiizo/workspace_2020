package main

import (
	"fmt"
	"github.com/lucperkins/rek"
	"time"
)

func main() {
	type Comment struct {
		Body string `json:"body"`
	}
	comment := Comment{Body: "This movie sucked"}
	headers := map[string]string{"My-Custom-Header", "foo,bar,baz"}
	res, err := rek.Post("https://httpbin.org/post",
		rek.Json(comment),
		rek.Headers(headers),
		rek.BasicAuth("user", "pass"),
		rek.Timeout(5*time.Second),
	)

	fmt.Println(res.StatusCode())
	fmt.Println(res.Text())
}
