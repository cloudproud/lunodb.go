package value

import (
	"encoding/binary"
	"fmt"
	"math"
)

func EncodeFloat32[T float32 | *float32](val T, buf []byte) ([]byte, error) {
	switch v := any(val).(type) {
	case float32:
		return binary.BigEndian.AppendUint32(buf, math.Float32bits(v)), nil
	case *float32:
		if v == nil {
			return buf, nil
		}

		return binary.BigEndian.AppendUint32(buf, math.Float32bits(*v)), nil
	default:
		return buf, fmt.Errorf("unsupported float32 type: %T", val)
	}
}

func EncodeFloat64[T float64 | *float64](val T, buf []byte) ([]byte, error) {
	switch v := any(val).(type) {
	case float64:
		return binary.BigEndian.AppendUint64(buf, math.Float64bits(v)), nil
	case *float64:
		if v == nil {
			return buf, nil
		}

		return binary.BigEndian.AppendUint64(buf, math.Float64bits(*v)), nil
	default:
		return buf, fmt.Errorf("unsupported float64 type: %T", val)
	}
}
