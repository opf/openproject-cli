package printer

import (
	"time"

	"github.com/briandowns/spinner"
)

type Func[T any] func() (T, error)

var loadingSpinner *spinner.Spinner

func init() {
	loadingSpinner = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	loadingSpinner.Suffix = " fetching data"
	_ = loadingSpinner.Color("yellow")
}

func WithSpinner[T any](f Func[T]) (T, error) {
	loadingSpinner.Start()

	t, err := f()

	loadingSpinner.Stop()

	return t, err
}
