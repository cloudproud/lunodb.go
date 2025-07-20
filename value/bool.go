package value

import (
	"fmt"
)

func EncodeBool[T bool | *bool](val T, buf []byte) ([]byte, error) {
	switch v := any(val).(type) {
	case bool:
		if v {
			buf = append(buf, 1)
		} else {
			buf = append(buf, 0)
		}
	case *bool:
		if v == nil {
			return buf, nil
		}

		if *v {
			buf = append(buf, 1)
		} else {
			buf = append(buf, 0)
		}
	default:
		return buf, fmt.Errorf("unsupported boolean type: %T", val)
	}

	return buf, nil
}
