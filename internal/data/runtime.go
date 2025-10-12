package data

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format, correct format (2 hr|min|sec)")

type Runtime struct {
	Duration int64
	Unit     string
}

func (rt Runtime) MarshalJSON() ([]byte, error) {
	jsonVal := fmt.Sprintf("%d %s", rt.Duration, rt.Unit)

	qoutedJsonVal := strconv.Quote(jsonVal)
	return []byte(qoutedJsonVal), nil
}

func (rt *Runtime) UnmarshalJSON(jsonValue []byte) error {

	unqoutedJSONValue, err := strconv.Unquote(string(jsonValue))

	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	parts := strings.Split(unqoutedJSONValue, " ")

	if len(parts) != 2 || !slices.Contains([]string{"hr", "min", "sec"}, parts[1]) {
		return ErrInvalidRuntimeFormat
	}
	i, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	*rt = Runtime{Duration: i, Unit: parts[1]}
	return nil
}
