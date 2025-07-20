package value

import (
	"encoding/binary"
	"fmt"
)

func EncodeInt8[T int8 | *int8](val T, buf []byte) ([]byte, error) {
	switch v := any(val).(type) {
	case int8:
		return append(buf, byte(v)), nil
	case *int8:
		if v == nil {
			return buf, nil
		}
		return append(buf, byte(*v)), nil
	default:
		return buf, fmt.Errorf("unsupported int8 type: %T", val)
	}
}

func EncodeInt16[T int16 | *int16](val T, buf []byte) ([]byte, error) {
	switch v := any(val).(type) {
	case int16:
		return binary.BigEndian.AppendUint16(buf, uint16(v)), nil
	case *int16:
		if v == nil {
			return buf, nil
		}
		return binary.BigEndian.AppendUint16(buf, uint16(*v)), nil
	default:
		return buf, fmt.Errorf("unsupported int16 type: %T", val)
	}
}

func EncodeInt32[T int32 | *int32](val T, buf []byte) ([]byte, error) {
	switch v := any(val).(type) {
	case int32:
		return binary.BigEndian.AppendUint32(buf, uint32(v)), nil
	case *int32:
		if v == nil {
			return buf, nil
		}
		return binary.BigEndian.AppendUint32(buf, uint32(*v)), nil
	default:
		return buf, fmt.Errorf("unsupported int32 type: %T", val)
	}
}

func EncodeInt64[T int64 | *int64 | int | *int](val T, buf []byte) ([]byte, error) {
	switch v := any(val).(type) {
	case int64:
		return binary.BigEndian.AppendUint64(buf, uint64(v)), nil
	case *int64:
		if v == nil {
			return buf, nil
		}
		return binary.BigEndian.AppendUint64(buf, uint64(*v)), nil
	case int:
		return binary.BigEndian.AppendUint64(buf, uint64(v)), nil
	case *int:
		if v == nil {
			return buf, nil
		}
		return binary.BigEndian.AppendUint64(buf, uint64(*v)), nil
	default:
		return buf, fmt.Errorf("unsupported int64 type: %T", val)
	}
}

func EncodeUint8[T uint8 | *uint8](val T, buf []byte) ([]byte, error) {
	switch v := any(val).(type) {
	case uint8:
		return append(buf, byte(v)), nil
	case *uint8:
		if v == nil {
			return buf, nil
		}
		return append(buf, byte(*v)), nil
	default:
		return buf, fmt.Errorf("unsupported uint8 type: %T", val)
	}
}

func EncodeUint16[T uint16 | *uint16](val T, buf []byte) ([]byte, error) {
	switch v := any(val).(type) {
	case uint16:
		return binary.BigEndian.AppendUint16(buf, v), nil
	case *uint16:
		if v == nil {
			return buf, nil
		}
		return binary.BigEndian.AppendUint16(buf, *v), nil
	default:
		return buf, fmt.Errorf("unsupported uint16 type: %T", val)
	}
}

func EncodeUint32[T uint32 | *uint32](val T, buf []byte) ([]byte, error) {
	switch v := any(val).(type) {
	case uint32:
		return binary.BigEndian.AppendUint32(buf, v), nil
	case *uint32:
		if v == nil {
			return buf, nil
		}
		return binary.BigEndian.AppendUint32(buf, *v), nil
	default:
		return buf, fmt.Errorf("unsupported uint32 type: %T", val)
	}
}

func EncodeUint64[T uint64 | *uint64](val T, buf []byte) ([]byte, error) {
	switch v := any(val).(type) {
	case uint64:
		return binary.BigEndian.AppendUint64(buf, v), nil
	case *uint64:
		if v == nil {
			return buf, nil
		}
		return binary.BigEndian.AppendUint64(buf, *v), nil
	default:
		return buf, fmt.Errorf("unsupported uint64 type: %T", val)
	}
}
