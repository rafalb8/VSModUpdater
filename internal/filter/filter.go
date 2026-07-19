package filter

import (
	"iter"
	"strings"
)

type ExclusionFilter[T any] []Rule

func NewExclusion[T any](expr string) (ExclusionFilter[T], error) {
	expr = strings.ReplaceAll(expr, ",", " ")

	filter := ExclusionFilter[T]{}
	for field := range strings.FieldsSeq(expr) {
		if field == "" {
			continue
		}

		field, newF, neg := strings.Cut(field, "^")
		if neg {
			field = newF
		}

		var rule Rule
		var err error

		before, after, isRange := strings.Cut(field, "-")
		if isRange {
			rule, err = RangeRule(before, after, neg)
		} else {
			rule, err = MatchRule(field, neg)
		}

		if err != nil {
			return nil, err
		}
		filter = append(filter, rule)
	}
	return filter, nil
}

func (f ExclusionFilter[T]) Check(idx int) bool {
	decision := Neutral
	for _, rule := range f {
		d := rule(idx)
		if d != Neutral {
			decision = d
		}
	}
	return decision == Exclude
}

func (f ExclusionFilter[T]) Filter(items iter.Seq2[int, T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i, val := range items {
			if !f.Check(i) {
				if !yield(val) {
					return
				}
			}
		}
	}
}
