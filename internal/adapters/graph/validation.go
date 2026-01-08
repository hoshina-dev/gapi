package graph

import "errors"

// validateTolerance ensures tolerance is not negative and returns nil if it's 0 or less
func validateTolerance(tolerance *float64) (*float64, error) {
	if tolerance == nil {
		return nil, nil
	}
	if *tolerance < 0 {
		return nil, errors.New("tolerance must be non-negative")
	}
	if *tolerance == 0 {
		return nil, nil
	}
	return tolerance, nil
}
