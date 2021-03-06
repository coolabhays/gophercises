package main

import (
	"fmt"
	"gophercises/link"
	"strings"
)

var exampleHtml = `
<html>
	<body>
		<h1>Hello!</h1>
		<a href="/other-page">
			A link
			<span>to other page</span>
		</a>
		<a href="/another-page">A link to another page</a>
	</body>
</html>`

func main() {
	r := strings.NewReader(exampleHtml)  // creating a Reader from string

	links, err := link.Parse(r)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", links)
}
