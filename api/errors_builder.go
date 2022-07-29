package api

import "fmt"

func buildArgumentRequiredError(argument string) error {
	return fmt.Errorf("%q is required", argument)
}
