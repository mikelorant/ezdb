package app

import (
	"fmt"

	"github.com/mikelorant/ezdb2/internal/selector"
)

func Select(value string, list []string, msg string) (string, error) {
	if value == "" {
		value, err := selector.Select(list,
			selector.WithMessage(msg),
		)
		if err != nil {
			return value, fmt.Errorf("unable to select: %w", err)
		}
		return value, nil
	}

	return value, nil
}

func SelectWithExclude(value string, list []string, msg string, exclude []string) (string, error) {
	if value == "" {
		value, err := selector.Select(list,
			selector.WithMessage(msg),
			selector.WithExclude(exclude),
		)
		if err != nil {
			return value, fmt.Errorf("unable to select: %w", err)
		}
		return value, nil
	}

	return value, nil
}
