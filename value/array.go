package value

func EncodeArray[T any](val []T, buf []byte) ([]byte, error) {
	return buf, nil
}

// func (a Array) Encode(buf []byte) ([]byte, error) {
// 	if a.Items == nil {
// 		return buf, nil
// 	}

// 	buf = binary.BigEndian.AppendUint32(buf, uint32(len(a.Items)))

// 	for index := range a.Items {
// 		frame, err := Encode(a.underlying, a.Items[index])
// 		if err != nil {
// 			return nil, err
// 		}

// 		buf = binary.BigEndian.AppendUint32(buf, uint32(len(frame)))
// 		buf = append(buf, frame...)
// 	}

// 	return buf, nil
// }
