package flows

// Unsigned8 represents a uint8 according to RFC5102
type Unsigned8 uint8

// Unsigned16 represents a uint16 according to RFC5102
type Unsigned16 uint16

// Unsigned32 represents a uint32 according to RFC5102
type Unsigned32 uint32

// Unsigned64 represents a uint64 according to RFC5102
type Unsigned64 uint64

// Signed8 represents a int8 according to RFC5102
type Signed8 int8

// Signed16 represents a int16 according to RFC5102
type Signed16 int16

// Signed32 represents a int32 according to RFC5102
type Signed32 int32

// Signed64 represents a int64 according to RFC5102
type Signed64 int64

// Float32 represents a float32 according to RFC5102
type Float32 float32

// Float64 represents a float64 according to RFC5102
type Float64 float64

// Number represents a number according to RFC5102
type Number interface {
	// Add numbers and return result
	Add(Number) Number //oper:a+b
	// Log returns log(number)
	Log() Number //oper:Float64(math.Log(float64(a)))

	// Less returns true if the number is smaller than the argument
	Less(Number) bool //oper:a<b
	// Greater returns true if the number is greater than the argument
	Greater(Number) bool //oper:a>b

	// ToFloat returns the number converted to float64
	ToFloat() float64 //oper:float64(a)
	// ToInt returns the number converted to int64
	ToInt() int64 //oper:int64(a)
	// ToUint returns the number converted to uint64
	ToUint() uint64 //oper:uint64(a)
	To64() Number
}

//go:generate go run tool/number_generate.go

// To64 returns the number converted to a 64 bit wide Number
func (a Unsigned8) To64() Number {
	return Signed64(a)
}

// To64 returns the number converted to a 64 bit wide Number
func (a Unsigned16) To64() Number {
	return Signed64(a)
}

// To64 returns the number converted to a 64 bit wide Number
func (a Unsigned32) To64() Number {
	return Signed64(a)
}

// To64 returns the number converted to a 64 bit wide Number
func (a Unsigned64) To64() Number {
	return Signed64(a)
}

// To64 returns the number converted to a 64 bit wide Number
func (a Signed8) To64() Number {
	return Signed64(a)
}

// To64 returns the number converted to a 64 bit wide Number
func (a Signed16) To64() Number {
	return Signed64(a)
}

// To64 returns the number converted to a 64 bit wide Number
func (a Signed32) To64() Number {
	return Signed64(a)
}

// To64 returns the number converted to a 64 bit wide Number
func (a Signed64) To64() Number {
	return Signed64(a)
}

// To64 returns the number converted to a 64 bit wide Number
func (a Float32) To64() Number {
	return Float64(a)
}

// To64 returns the number converted to a 64 bit wide Number
func (a Float64) To64() Number {
	return Float64(a)
}