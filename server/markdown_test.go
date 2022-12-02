package server

import (
	"fmt"
	"testing"
)

func TestToMarkdown(t *testing.T) {
	doc := `!_navbar.md

* [Github地址](https://github.com/qianlnk)

	 ![](/uploads/upload_bf0640d65fd27713e617dea2e03679f4.png)

test
`

	md, _ := toMarkdown(doc)

	fmt.Println(*md)

	md.localizaiton("/Users/xiezhenjia/go/src/github.com/qianlnk/codidoc/docs")
}
