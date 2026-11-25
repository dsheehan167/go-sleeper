package sleeper

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// FlexibleString can unmarshal from either a string or int
type FlexibleString string

// UnmarshalJSON implements json.Unmarshaler, accepting a JSON string or integer.
func (f *FlexibleString) UnmarshalJSON(data []byte) error {
	// Try string first
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*f = FlexibleString(s)
		return nil
	}

	// Try int
	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		*f = FlexibleString(strconv.Itoa(i))
		return nil
	}

	return errors.New("cannot unmarshal into FlexibleString")
}

// String returns the string value
func (f FlexibleString) String() string {
	return string(f)
}

// Int converts to int, returns error if conversion fails
func (f FlexibleString) Int() (int, error) {
	i, err := strconv.Atoi(string(f))
	if err != nil {
		return 0, fmt.Errorf("cannot convert %s to int: %w", f, err)
	}
	return i, nil
}
