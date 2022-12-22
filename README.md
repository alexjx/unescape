# unescape

This is an utility library for remove ANSI sequence from text. It is created for helping processing ANSI text produced in CI jobs.

It is inspired by [muesli/ansi](https://github.com/muesli/ansi)

## Unescape Writer

```go
import "github.com/alexjx/unescape"

w := unescape.Unescaper{Forward: os.Stdout}
w.Write([]byte("\x1b[31mHello, world!\x1b[0m"))
w.Close()
```
