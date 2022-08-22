# tokens

![test](https://github.com/wmentor/tokens/workflows/test/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/wmentor/tokens/badge.svg?branch=master&v=1.0.4)](https://coveralls.io/github/wmentor/tokens?branch=master)
[![https://goreportcard.com/report/github.com/wmentor/tokens](https://goreportcard.com/badge/github.com/wmentor/tokens)](https://goreportcard.com/report/github.com/wmentor/tokens)
[![https://pkg.go.dev/github.com/wmentor/tokens](https://pkg.go.dev/badge/github.com/wmentor/tokens.svg)](https://pkg.go.dev/github.com/wmentor/tokens)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Text to tokens Go library.

# Token get insensitive mode

```go
package main

import (
	"fmt"
	"strings"

	"github.com/wmentor/tokens"
)

func main() {

	txt := "Hello, my little friend!"

	tokenizer := tokens.New(strings.NewReader(txt))

	for {
		tok, err := range tokenizer.Token()
		if err != nil { // io.EOF
			break
		}
		fmt.Println(tok)
	}
}
```

Result:

```
hello
,
my
little
friend
!
```

Case sensitive mode:

```go
package main

import (
	"fmt"
	"strings"

	"github.com/wmentor/tokens"
)

func main() {

	txt := "Hello, my little friend!"

	tokenizer := tokens.New(strings.NewReader(txt), tokens.WithCaseSensitive())

	for {
		tok, err := range tokenizer.Token()
		if err != nil { // io.EOF
			break
		}
		fmt.Println(tok)
	}
}
```

Result:

```
Hello
,
my
liTTle
friend
!
```
