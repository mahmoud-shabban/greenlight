package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format, correct format: duration mins")

type Runtime int64

func (rt Runtime) MarshalJSON() ([]byte, error) {
	jsonVal := fmt.Sprintf("%d mins", rt)

	qoutedJsonVal := strconv.Quote(jsonVal)
	return []byte(qoutedJsonVal), nil
}

func (rt *Runtime) UnmarshalJSON(jsonValue []byte) error {

	unqoutedJSONValue, err := strconv.Unquote(string(jsonValue))

	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	parts := strings.Split(unqoutedJSONValue, " ")

	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	i, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	*rt = Runtime(i)
	return nil
}
