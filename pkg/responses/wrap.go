package responses

import "fmt"

func Wrap(ctx string, err error) error {
	return fmt.Errorf("%s: %w", ctx, err)
}
