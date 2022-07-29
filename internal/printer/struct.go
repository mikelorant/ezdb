package printer

import (
	"bytes"
	"io"
	"os"

	"github.com/alecthomas/chroma/quick"
	"github.com/goccy/go-yaml"
)

func Struct(v interface{}) string {
	y, err := yaml.Marshal(v)
	if err != nil {
		return ""
	}

	o, _ := os.Stdout.Stat()
	if (o.Mode() & os.ModeCharDevice) != os.ModeCharDevice {
		return string(y)
	}

	var b bytes.Buffer
	if err := quick.Highlight(&b, string(y), "yaml", "terminal256", "pygments"); err != nil {
		return string(y)
	}

	yh, err := io.ReadAll(&b)
	if err != nil {
		return string(y)
	}

	return string(yh)
}
