package data

import (
	"fmt"
	"strconv"
)

type Runtime struct {
	Duration int64
	Unit     string
}

func (rt Runtime) MarshalJSON() ([]byte, error) {
	jsonVal := fmt.Sprintf("%d %s", rt.Duration, rt.Unit)

	qoutedJsonVal := strconv.Quote(jsonVal)
	return []byte(qoutedJsonVal), nil
}
