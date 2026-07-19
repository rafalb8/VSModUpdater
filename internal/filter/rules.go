package filter

import (
	"cmp"
	"fmt"
	"strconv"
)

type Decision uint8

const (
	Neutral Decision = iota
	Exclude
	Include
)

type Rule func(idx int) Decision

func RangeRule(before, after string, neg bool) (Rule, error) {
	min, err1 := strconv.Atoi(before)
	max, err2 := strconv.Atoi(after)
	if err := cmp.Or(err1, err2); err != nil {
		return nil, fmt.Errorf("invalid int in range: %s-%s", before, after)
	}

	return func(n int) Decision {
		if n >= min && n <= max {
			if neg {
				return Include
			}
			return Exclude
		}
		return Neutral
	}, nil
}

func MatchRule(field string, neg bool) (Rule, error) {
	val, err := strconv.Atoi(field)
	if err != nil {
		return nil, fmt.Errorf("invalid int: %s", field)
	}

	return func(n int) Decision {
		if n == val {
			if neg {
				return Include
			}
			return Exclude
		}
		return Neutral
	}, nil
}
