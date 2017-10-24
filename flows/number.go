package flows

import (
	"fmt"

	ipfix "pm.cn.tuwien.ac.at/ipfix/go-ipfix"
)

// DateTimeSeconds represents time in units of seconds from 00:00 UTC, Januray 1, 1970 according to RFC5102.
type DateTimeSeconds = ipfix.DateTimeSeconds

// DateTimeMilliseconds represents time in units of milliseconds from 00:00 UTC, Januray 1, 1970 according to RFC5102.
type DateTimeMilliseconds = ipfix.DateTimeMilliseconds

// DateTimeMicroseconds represents time in units of microseconds from 00:00 UTC, Januray 1, 1970 according to RFC5102.
type DateTimeMicroseconds = ipfix.DateTimeMicroseconds

// DateTimeNanoseconds represents time in units of nanoseconds from 00:00 UTC, Januray 1, 1970 according to RFC5102.
type DateTimeNanoseconds = ipfix.DateTimeNanoseconds

func ToFloat(a interface{}) float64 {
	switch i := a.(type) {
	case float64:
		return i
	case float32:
		return float64(i)
	case int64:
		return float64(i)
	case int32:
		return float64(i)
	case int16:
		return float64(i)
	case int8:
		return float64(i)
	case int:
		return float64(i)
	case uint64:
		return float64(i)
	case uint32:
		return float64(i)
	case uint16:
		return float64(i)
	case uint8:
		return float64(i)
	case uint:
		return float64(i)
	case DateTimeSeconds:
		return float64(i)
	case DateTimeMilliseconds:
		return float64(i)
	case DateTimeMicroseconds:
		return float64(i)
	case DateTimeNanoseconds:
		return float64(i)
	case nil:
		return 0
	case bool:
		if i {
			return 1
		}
		return 0
	}
	panic(fmt.Sprintf("Can't convert %v to float", a))
}

func ToInt(a interface{}) int64 {
	switch i := a.(type) {
	case float64:
		return int64(i)
	case float32:
		return int64(i)
	case int64:
		return i
	case int32:
		return int64(i)
	case int16:
		return int64(i)
	case int8:
		return int64(i)
	case int:
		return int64(i)
	case uint64:
		return int64(i)
	case uint32:
		return int64(i)
	case uint16:
		return int64(i)
	case uint8:
		return int64(i)
	case uint:
		return int64(i)
	case DateTimeSeconds:
		return int64(i)
	case DateTimeMilliseconds:
		return int64(i)
	case DateTimeMicroseconds:
		return int64(i)
	case DateTimeNanoseconds:
		return int64(i)
	case nil:
		return 0
	case bool:
		if i {
			return 1
		}
		return 0
	}
	panic(fmt.Sprintf("Can't convert %v to int", a))
}

func ToUInt(a interface{}) uint64 {
	switch i := a.(type) {
	case float64:
		return uint64(i)
	case float32:
		return uint64(i)
	case int64:
		return uint64(i)
	case int32:
		return uint64(i)
	case int16:
		return uint64(i)
	case int8:
		return uint64(i)
	case int:
		return uint64(i)
	case uint64:
		return i
	case uint32:
		return uint64(i)
	case uint16:
		return uint64(i)
	case uint8:
		return uint64(i)
	case uint:
		return uint64(i)
	case DateTimeSeconds:
		return uint64(i)
	case DateTimeMilliseconds:
		return uint64(i)
	case DateTimeMicroseconds:
		return uint64(i)
	case DateTimeNanoseconds:
		return uint64(i)
	case nil:
		return 0
	case bool:
		if i {
			return 1
		}
		return 0
	}
	panic(fmt.Sprintf("Can't convert %v to int", a))
}

type NumberType int

const (
	IntType NumberType = iota
	UIntType
	FloatType
	SecondsType
	MillisecondsType
	MicrosecondsType
	NanosecondsType
)

func cleanUp(a interface{}) (NumberType, NumberType, interface{}) {
	switch i := a.(type) {
	case float64:
		return FloatType, FloatType, a
	case float32:
		return FloatType, FloatType, float64(i)
	case int64:
		return IntType, IntType, a
	case int32:
		return IntType, IntType, int64(i)
	case int16:
		return IntType, IntType, int64(i)
	case int8:
		return IntType, IntType, int64(i)
	case int:
		return IntType, IntType, int64(i)
	case uint64:
		return UIntType, UIntType, a
	case uint32:
		return UIntType, UIntType, uint64(i)
	case uint16:
		return UIntType, UIntType, uint64(i)
	case uint8:
		return UIntType, UIntType, uint64(i)
	case uint:
		return UIntType, UIntType, uint64(i)
	case DateTimeSeconds:
		return SecondsType, UIntType, a
	case DateTimeMilliseconds:
		return MillisecondsType, UIntType, a
	case DateTimeMicroseconds:
		return MicrosecondsType, UIntType, a
	case DateTimeNanoseconds:
		return NanosecondsType, UIntType, a
	case bool:
		if i {
			return UIntType, UIntType, uint64(1)
		}
		return UIntType, UIntType, uint64(0)
	}
	panic(fmt.Sprintf("Can't upconvert %v", a))
}

func uintToFloat(val interface{}) float64 {
	return float64(val.(uint64))
}

func intToFloat(val interface{}) float64 {
	return float64(val.(int64))
}

func uintToInt(val interface{}) int64 {
	return int64(val.(uint64))
}

func floatToInt(val interface{}) int64 {
	return int64(val.(float64))
}

func intToTime(val interface{}, kind NumberType) interface{} {
	switch kind {
	case SecondsType:
		return DateTimeSeconds(val.(int64))
	case MillisecondsType:
		return DateTimeMilliseconds(val.(int64))
	case MicrosecondsType:
		return DateTimeMicroseconds(val.(int64))
	case NanosecondsType:
		return DateTimeNanoseconds(val.(int64))
	}
	panic("This should never happen")
}

func uintToTime(val interface{}, kind NumberType) interface{} {
	switch kind {
	case SecondsType:
		return DateTimeSeconds(val.(uint64))
	case MillisecondsType:
		return DateTimeMilliseconds(val.(uint64))
	case MicrosecondsType:
		return DateTimeMicroseconds(val.(uint64))
	case NanosecondsType:
		return DateTimeNanoseconds(val.(int64))
	}
	panic("This should never happen")
}

func floatToTime(val interface{}, kind NumberType) interface{} {
	switch kind {
	case SecondsType:
		return DateTimeSeconds(val.(float64))
	case MillisecondsType:
		return DateTimeMilliseconds(val.(float64))
	case MicrosecondsType:
		return DateTimeMicroseconds(val.(float64))
	case NanosecondsType:
		return DateTimeNanoseconds(val.(float64))
	}
	panic("This should never happen")
}

func scaleTimetoNano(from NumberType, val interface{}) interface{} {
	switch from {
	case SecondsType:
		return DateTimeSeconds(val.(DateTimeSeconds) * 1e9)
	case MillisecondsType:
		return DateTimeMilliseconds(val.(DateTimeMilliseconds) * 1e6)
	case MicrosecondsType:
		return DateTimeMicroseconds(val.(DateTimeMicroseconds) * 1e3)
	case NanosecondsType:
		return val
	}
	panic("This should never happen")
}

// UpConvert returns either two Signed64 or two Float64 depending on the numbers
func UpConvert(a, b interface{}) (dst NumberType, family NumberType, ai, bi interface{}) {
	var tA, tB, fA, fB NumberType
	tA, fA, ai = cleanUp(a)
	tB, fB, bi = cleanUp(b)
	if tA == tB {
		dst = tA
		family = fA
		return
	}
	if tA == IntType {
		if tB == UIntType {
			return tA, fA, ai, uintToInt(bi)
		}
		if tB == FloatType {
			return tB, fB, intToFloat(ai), bi
		}
		return tB, fB, intToTime(ai, tB), tB
	}
	if tA == UIntType {
		if tB == IntType {
			return tB, fB, uintToInt(ai), bi
		}
		if tB == FloatType {
			return tB, fB, uintToFloat(ai), bi
		}
		return tB, fB, uintToTime(ai, tB), tB
	}
	if tA == FloatType {
		if tB == UIntType {
			return tA, fA, ai, uintToFloat(bi)
		}
		if tB == IntType {
			return tA, fA, ai, intToFloat(bi)
		}
		return tB, fB, floatToTime(ai, tB), tB
	}
	// both types are time - but differen timebases
	return NanosecondsType, UIntType, scaleTimetoNano(tA, a), scaleTimetoNano(tB, b)
}

func FixType(val interface{}, t NumberType) interface{} {
	switch t {
	case SecondsType:
		return DateTimeSeconds(val.(uint64))
	case MillisecondsType:
		return DateTimeMilliseconds(val.(uint64))
	case MicrosecondsType:
		return DateTimeMicroseconds(val.(uint64))
	case NanosecondsType:
		return DateTimeNanoseconds(val.(uint64))
	}
	return val
}

func cleanUpType(a ipfix.Type) ipfix.Type {
	switch a {
	case ipfix.OctetArrayType:
		return ipfix.IllegalType
	case ipfix.Unsigned8Type:
		return ipfix.Unsigned64Type
	case ipfix.Unsigned16Type:
		return ipfix.Unsigned64Type
	case ipfix.Unsigned32Type:
		return ipfix.Unsigned64Type
	case ipfix.Unsigned64Type:
		return ipfix.Unsigned64Type
	case ipfix.Signed8Type:
		return ipfix.Signed64Type
	case ipfix.Signed16Type:
		return ipfix.Signed64Type
	case ipfix.Signed32Type:
		return ipfix.Signed64Type
	case ipfix.Signed64Type:
		return ipfix.Signed64Type
	case ipfix.Float32Type:
		return ipfix.Float64Type
	case ipfix.Float64Type:
		return ipfix.Float64Type
	case ipfix.BooleanType:
		return ipfix.Unsigned64Type
	case ipfix.MacAddressType:
		return ipfix.IllegalType
	case ipfix.StringType:
		return ipfix.IllegalType
	case ipfix.DateTimeSecondsType:
		return a
	case ipfix.DateTimeMillisecondsType:
		return a
	case ipfix.DateTimeMicrosecondsType:
		return a
	case ipfix.DateTimeNanosecondsType:
		return a
	case ipfix.Ipv4AddressType:
		return ipfix.IllegalType
	case ipfix.Ipv6AddressType:
		return ipfix.IllegalType
	}
	panic(fmt.Sprintf("Can't upconvert %s", a))
}

func UpConvertTypes(a, b ipfix.Type) ipfix.Type {
	tA := cleanUpType(a)
	tB := cleanUpType(b)
	if tA == tB {
		return tA
	}
	if tA == ipfix.Signed64Type {
		if tB == ipfix.Unsigned64Type {
			return tA
		}
		if tB == ipfix.Float64Type {
			return tB
		}
		return tB
	}
	if tA == ipfix.Unsigned64Type {
		if tB == ipfix.Signed64Type {
			return tB
		}
		if tB == ipfix.Float64Type {
			return tB
		}
		return tB
	}
	if tA == ipfix.Float64Type {
		if tB == ipfix.Unsigned64Type {
			return tA
		}
		if tB == ipfix.Signed64Type {
			return tA
		}
		return tB
	}
	// both types are time - but differen timebases
	return ipfix.DateTimeNanosecondsType
}
