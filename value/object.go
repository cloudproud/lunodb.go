package value

import (
	"encoding/binary"

	lunopb "github.com/cloudproud/lunodb.api/proto/types"
	"github.com/gogo/protobuf/proto"
)

func EncodeObject(val map[string]any, buf []byte) (frame []byte, err error) {
	var typed *lunopb.Type
	for k, v := range val {
		buf = append(buf, []byte(k)...)
		buf = append(buf, 0)

		typed, frame, err = Encode(v, []byte{})
		if err != nil {
			return nil, err
		}

		header, err := proto.Marshal(typed)
		if err != nil {
			return nil, err
		}

		buf = make([]byte, len(header)+4)
		binary.BigEndian.PutUint32(buf[:4], uint32(len(header)))
		copy(buf[4:], header)

		buf = binary.BigEndian.AppendUint32(buf, uint32(len(frame)))
		buf = append(buf, frame...)
	}

	return buf, nil
}
