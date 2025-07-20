package value

import (
	"fmt"
	"net/netip"
)

func EncodeInet[T netip.Prefix | *netip.Prefix](val T, buf []byte) ([]byte, error) {
	switch v := any(val).(type) {
	case netip.Prefix:
		return v.AppendBinary(buf)
	case *netip.Prefix:
		if v == nil {
			return buf, nil
		}

		return v.AppendBinary(buf)
	default:
		return buf, fmt.Errorf("unsupported inet type: %T", val)
	}
}
