package structutil

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/alecthomas/chroma/quick"
	"github.com/goccy/go-yaml"
)

func Sprint(v interface{}) (string, error) {
	y, err := yaml.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("unable to marshal struct to yaml: %w", err)
	}

	o, _ := os.Stdout.Stat()
	if (o.Mode() & os.ModeCharDevice) != os.ModeCharDevice {
		return string(y), nil
	}

	var b bytes.Buffer
	if err := quick.Highlight(&b, string(y), "yaml", "terminal256", "pygments"); err != nil {
		return string(y), fmt.Errorf("unable to highlight yaml: %w", err)
	}

	yh, err := io.ReadAll(&b)
	if err != nil {
		return string(y), fmt.Errorf("unable to read buffer: %w", err)
	}

	return string(yh), nil
}
