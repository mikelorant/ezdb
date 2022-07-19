package progress

import (
	"fmt"
	"time"

	"github.com/schollz/progressbar/v3"
)

func New(size int64, desc string) *progressbar.ProgressBar {
	return progressbar.NewOptions64(size,
		progressbar.OptionSetDescription(desc),
		progressbar.OptionOnCompletion(func() {
			fmt.Printf("\n")
		}),
		progressbar.OptionThrottle(64*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionShowBytes(true),
	)
}
