package flatten

import (
	"errors"
	"fmt"
	"strconv"
)

type SeparatorStyle struct {
	Before string // Prepend to key
	Middle string // Add between keys
	After  string // Append to key
}

// Default styles
var (
	// Separate nested key components with dots, e.g. "a.b.1.c.d"
	DotStyle = SeparatorStyle{Middle: "."}

	// Separate with path-like slashes, e.g. a/b/1/c/d
	PathStyle = SeparatorStyle{Middle: "/"}

	// Separate ala Rails, e.g. "a[b][c][1][d]"
	RailsStyle = SeparatorStyle{Before: "[", After: "]"}

	// Separate with underscores, e.g. "a_b_1_c_d"
	UnderscoreStyle = SeparatorStyle{Middle: "_"}
)

// Nested input must be a map or slice
var NotValidInputError = errors.New("Not a valid input: map or slice")

func Flatten(nested map[string]interface{}, prefix string, style SeparatorStyle) (map[string]interface{}, error) {
	flatmap := make(map[string]interface{})

	err := flatten(true, flatmap, nested, prefix, style)
	if err != nil {
		return nil, err
	}

	return flatmap, nil
}

func flatten(top bool, flatMap map[string]interface{}, nested interface{}, prefix string, style SeparatorStyle) error {
	assign := func(newKey string, v interface{}) error {
		switch v.(type) {
		case map[string]interface{}, map[interface{}]interface{}, []interface{}:
			if err := flatten(false, flatMap, v, newKey, style); err != nil {
				return err
			}
		default:
			flatMap[newKey] = v
		}

		return nil
	}

	switch nested.(type) {
	case map[string]interface{}:
		for k, v := range nested.(map[string]interface{}) {
			newKey := enkey(top, prefix, k, style)
			assign(newKey, v)
		}
	case map[interface{}]interface{}:
		for k, v := range nested.(map[interface{}]interface{}) {
			newKey := enkey(top, prefix, k, style)
			assign(newKey, v)
		}
	case []interface{}:
		for i, v := range nested.([]interface{}) {
			newKey := enkey(top, prefix, strconv.Itoa(i), style)
			assign(newKey, v)
		}
	default:
		return NotValidInputError
	}

	return nil
}

func enkey(top bool, prefix string, subkey interface{}, style SeparatorStyle) string {
	key := prefix

	if top {
		key += fmt.Sprintf("%v", subkey)
	} else {
		key += style.Before + style.Middle + fmt.Sprintf("%v", subkey) + style.After
	}

	return key
}
