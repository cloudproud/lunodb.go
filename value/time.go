package value

// import (
// 	"encoding/binary"
// 	"fmt"
// 	"time"

// 	"github.com/cloudproud/lunodb.go/value/duration"
// )

// func NewDate(value time.Time, null bool) *Date {
// 	return &Date{
// 		value: value,
// 		null:  null,
// 	}
// }

// type Date struct {
// 	value time.Time
// 	null  bool
// }

// func (d Date) String() string {
// 	if d.null {
// 		return "date(nil)"
// 	}

// 	return fmt.Sprintf("date(%s)", d.value.Format(time.DateOnly))
// }

// func (d Date) Encode(buf []byte) ([]byte, error) {
// 	if d.null {
// 		return buf, nil
// 	}

// 	enc, err := d.value.MarshalBinary()
// 	return append(buf, enc...), err
// }

// type Time struct {
// 	value time.Time
// 	null  bool
// }

// func NewTime(value time.Time, null bool) *Time {
// 	return &Time{
// 		value: value,
// 		null:  null,
// 	}
// }

// func (i Time) String() string {
// 	if i.null {
// 		return "time(nil)"
// 	}

// 	return fmt.Sprintf("time(%s)", i.value.Format(time.TimeOnly))
// }

// func (t Time) Encode(buf []byte) ([]byte, error) {
// 	if t.null {
// 		return buf, nil
// 	}

// 	enc, err := t.value.MarshalBinary()
// 	return append(buf, enc...), err
// }

// type Timestamp struct {
// 	value time.Time
// 	null  bool
// }

// func NewTimestamp(value time.Time, null bool) *Timestamp {
// 	return &Timestamp{
// 		value: value,
// 		null:  null,
// 	}
// }

// func (i Timestamp) String() string {
// 	if i.null {
// 		return "timestamp(nil)"
// 	}

// 	return fmt.Sprintf("timestamp(%s)", i.value.Format(time.RFC3339Nano))
// }

// func (t Timestamp) Encode(buf []byte) ([]byte, error) {
// 	if t.null {
// 		return buf, nil
// 	}

// 	enc, err := t.value.MarshalBinary()
// 	return append(buf, enc...), err
// }

// type Duration struct {
// 	value duration.Duration
// 	null  bool
// }

// func NewDuration(value duration.Duration, null bool) *Duration {
// 	return &Duration{
// 		value: value,
// 		null:  null,
// 	}
// }

// func (i Duration) String() string {
// 	if i.null {
// 		return "duration(nil)"
// 	}

// 	return fmt.Sprintf("duration(%s)", i.value.String())
// }

// func (i Duration) Encode(buf []byte) ([]byte, error) {
// 	if i.null {
// 		return buf, nil
// 	}

// 	nano, month, days := i.value.EncodeBigInt()
// 	buf = append(buf, make([]byte, 24)...)
// 	binary.BigEndian.PutUint64(buf[len(buf)-24:], uint64(days))
// 	binary.BigEndian.PutUint64(buf[len(buf)-16:], uint64(month))
// 	binary.BigEndian.PutUint64(buf[len(buf)-8:], uint64(nano.Int64()))

// 	return buf, nil
// }
