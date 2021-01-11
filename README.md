# tokens

There are two ways fetch tokens from text. First is token chan and second is token callback.

# Token chan

```go
package main

import (
	"fmt"
	"strings"

	"github.com/wmentor/tokens"
)

func main() {

	txt := "Hello, my little friend!"

	for tok := range tokens.Stream(strings.NewReader(txt)) {
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



# Token callback

```go
package main

import (
	"fmt"
	"strings"

	"github.com/wmentor/tokens"
)

func main() {

	txt := "Hello, my little friend!"

	tokens.Process(strings.NewReader(txt), func(w string) {
		fmt.Println(w)
	})
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

	txt := "Hello, my liTtLe fRiEnd!"

	tokens.Process(strings.NewReader(txt), func(w string) {
		fmt.Println(w)
	}, tokens.OptCaseSensitive)
}
```

Result:

```
Hello
,
my
liTtLe
fRiEnd
!
```

```go
package main

import (
	"fmt"
	"strings"

	"github.com/wmentor/tokens"
)

func main() {

	txt := "Hello, my little frIend!"

	tokens.Process(strings.NewReader(txt), func(w string) {
		fmt.Println(w)
	}, tokens.OptCaseSensitive)
}
```

Result:

```
Hello
,
my
little
frIend
!
```

Unlike the first case, we don't create new chan and goroutine. This method is more efficient especially when we process a large number of short lines.
