package goposit

import (
	"math/big"
	"math/bits"
)

// Posit8 is an 8 bit posit with 0 exponent bits
type Posit8 struct{ impl *SlowPosit }

// NewPosit8 makes a new posit with 8 bits and 0 es bits, the initial value is zero
func NewPosit8() Posit8 { return Posit8{impl: NewSlowPosit(8, 0)} }

// Add takes the sum of two posits
// p.Add(x) creates a new posit z which is p+x, neither p nor x are altered
func (p Posit8) Add(x Posit8) Posit8 { return Posit8{impl: p.impl.Add(x.impl)} }

// AddExact returns exactly the sum of two posits, represented as two
// more posits, the first result is a posit which is the sum truncated
// (round to negative infinity) and the second posit is the difference between
// the first posit and the true sum. That is to say true_sum = out2 + out1.
func (p Posit8) AddExact(x Posit8) (Posit8, Posit8) {
	res, diff := p.impl.AddExact(x.impl)
	return Posit8{impl: res}, Posit8{impl: diff}
}

// Sub takes the difference of two posits
// p.Sub(x) creates a new posit z which is p-x, neither p nor x are altered
func (p Posit8) Sub(x Posit8) Posit8 { return Posit8{impl: p.impl.Sub(x.impl)} }

// SubExact returns exactly the difference of two posits, represented as two
// more posits, the first result is a posit which is the difference truncated
// (round to negative infinity) and the second posit is the difference between
// the first posit and the true sum. That is to say true_sum = out2 + out1.
func (p Posit8) SubExact(x Posit8) (Posit8, Posit8) {
	res, diff := p.impl.SubExact(x.impl)
	return Posit8{impl: res}, Posit8{impl: diff}
}

// Mul takes the product of two posits
// p.Mul(x) creates a new posit z which is x*y, neither p nor x are altered
func (p Posit8) Mul(x Posit8) Posit8 { return Posit8{impl: p.impl.Mul(x.impl)} }

// MulPromote takes the product of two posits
// p.Mul(x) creates a new posit z which is p*x represented as the next larger posit size
// next larger means the bit width is doubled and the exponent size is increased by 1
// this function is guaranteed not to round
func (p Posit8) MulPromote(x Posit8) Posit16 { return Posit16{impl: p.impl.MulPromote(x.impl)} }

// Div takes the quotent of two posits
// p.Div(x) creates a new posit z which is p/x, neither p nor x are altered
func (p Posit8) Div(x Posit8) Posit8 { return Posit8{impl: p.impl.Div(x.impl)} }

// DivPromote takes the quotent of two posits
// p.DivPromote(x) creates a new posit z which is p/x, as the next larger posit size
// next larger means the bit width is doubled and the exponent size is increased by 1
// this function is guaranteed not to round
func (p Posit8) DivPromote(x Posit8) Posit16 { return Posit16{impl: p.impl.DivPromote(x.impl)} }

// Sqrt finds the square root of a posit
// p.Sqrt() returns a new posit, p is not altered
func (p Posit8) Sqrt() Posit8 { return Posit8{impl: p.impl.Sqrt()} }

// Size specific

// FromInt creates a new posit which is set to the value of a integer
// p.FromInt() outputs a new posit z, p is not altered
func (p Posit8) FromInt(i int8) Posit8 { return Posit8{impl: p.impl.FromInt(int64(i))} }

// FromUint creates a new posit which is set to the value of an unsigned integer
// p.FromUint() outputs a new posit z, p is not altered
func (p Posit8) FromUint(i uint8) Posit8 { return Posit8{impl: p.impl.FromUint(uint64(i))} }

// Int outputs an int8 representation of the value of the posit
// lost bits are rounded to nearest even and if the number is greater than or equal
// to two to the 8 -1 power, the maximum 0x7f for positive input or
// 0xff for negative input will return.
// Likewise if the number is less than or equal to 1/2 it will round down to zero.
func (p Posit8) Int() int8 {
	x := p.impl.Int()
	if x > 0x7f {
		return 0x7f
	}
	if x < -0x7f {
		return -0x7f
	}
	return int8(x)
}

// Uint outputs a uint8 representation of the value of the posit, if the posit
// is less than or equal to 1/2, 0 will be returned and if the value is greater
// than two to the maxBits power, 0xffffffffffffffff will be returned.
func (p Posit8) Uint() uint8 {
	x := p.impl.Int()

	if x > 0xff {
		return 0xff
	}

	return uint8(x)
}

// Exp outputs the exponent of the posit, specifically the number z for which
// 0.5*2**z <= positValue < 1*2**z
func (p Posit8) Exp() int8 {
	x := p.impl.Exp()

	if x > 0x7f || x < -0x7f {
		panic("posit with bigger exponent than should be possible")
	}

	return int8(x)
}

// Mant outputs a posit with the same mantissa but an exponent of 0, that is a
// number which is greater than or equal to 0.5 and less than 1
func (p Posit8) Mant() Posit8 { return Posit8{impl: p.impl.Mant()} }

// ExpAdd returns a new posit with x added to the exponent, effectively multiplying the
// posit value by 2**x. It is equivilent to a bit-shift in integer math. If bits is
// negative then the exponent will be decreased, if bits causes the exponent to increase
// in magnitude, it might cause the posit to round.
func (p Posit8) ExpAdd(x int8) Posit8 { return Posit8{impl: p.impl.ExpAdd(int32(x))} }

// Up returns the same number up-casted to the next larger posit, in this case larger means
// twice the nbits and an exponent size which is one greater.
func (p Posit8) Up() Posit16 {
	return Posit16{impl: p.impl.Up()}
}

// Bits outputs a uint8 containing the raw binary format of the posit
func (p Posit8) Bits() uint8 {
	words := p.impl.Bits.Bits()
	if len(words) == 1 {
		return uint8(words[0])
	} else if len(words) == 2 {
		if bits.UintSize != 32 {
			panic("unexpected number of bits")
		}
		return (uint8(words[1]) << (8 / 2)) | uint8(words[0])
	}
	panic("unexpected word size")
}

// SetBits outputs a new posit with the bits set to those which you specify
func (p Posit8) SetBits(bits uint8) Posit8 {
	bi := new(big.Int).SetUint64(uint64(bits))
	return Posit8{impl: &SlowPosit{es: p.impl.es, nbits: p.impl.nbits, Bits: bi}}
}

// Clone makes a copy of a posit
func (p Posit8) Clone() Posit8 {
	return Posit8{impl: p.impl.Clone()}
}

// Posit16 is an 16 bit posit with 1 exponent bits
type Posit16 struct{ impl *SlowPosit }

// NewPosit16 makes a new posit with 16 bits and 1 es bits, the initial value is zero
func NewPosit16() Posit16 { return Posit16{impl: NewSlowPosit(16, 1)} }

// Add takes the sum of two posits
// p.Add(x) creates a new posit z which is p+x, neither p nor x are altered
func (p Posit16) Add(x Posit16) Posit16 { return Posit16{impl: p.impl.Add(x.impl)} }

// AddExact returns exactly the sum of two posits, represented as two
// more posits, the first result is a posit which is the sum truncated
// (round to negative infinity) and the second posit is the difference between
// the first posit and the true sum. That is to say true_sum = out2 + out1.
func (p Posit16) AddExact(x Posit16) (Posit16, Posit16) {
	res, diff := p.impl.AddExact(x.impl)
	return Posit16{impl: res}, Posit16{impl: diff}
}

// Sub takes the difference of two posits
// p.Sub(x) creates a new posit z which is p-x, neither p nor x are altered
func (p Posit16) Sub(x Posit16) Posit16 { return Posit16{impl: p.impl.Sub(x.impl)} }

// SubExact returns exactly the difference of two posits, represented as two
// more posits, the first result is a posit which is the difference truncated
// (round to negative infinity) and the second posit is the difference between
// the first posit and the true sum. That is to say true_sum = out2 + out1.
func (p Posit16) SubExact(x Posit16) (Posit16, Posit16) {
	res, diff := p.impl.SubExact(x.impl)
	return Posit16{impl: res}, Posit16{impl: diff}
}

// Mul takes the product of two posits
// p.Mul(x) creates a new posit z which is x*y, neither p nor x are altered
func (p Posit16) Mul(x Posit16) Posit16 { return Posit16{impl: p.impl.Mul(x.impl)} }

// MulPromote takes the product of two posits
// p.Mul(x) creates a new posit z which is p*x represented as the next larger posit size
// next larger means the bit width is doubled and the exponent size is increased by 1
// this function is guaranteed not to round
func (p Posit16) MulPromote(x Posit16) Posit32 { return Posit32{impl: p.impl.MulPromote(x.impl)} }

// Div takes the quotent of two posits
// p.Div(x) creates a new posit z which is p/x, neither p nor x are altered
func (p Posit16) Div(x Posit16) Posit16 { return Posit16{impl: p.impl.Div(x.impl)} }

// DivPromote takes the quotent of two posits
// p.DivPromote(x) creates a new posit z which is p/x, as the next larger posit size
// next larger means the bit width is doubled and the exponent size is increased by 1
// this function is guaranteed not to round
func (p Posit16) DivPromote(x Posit16) Posit32 { return Posit32{impl: p.impl.DivPromote(x.impl)} }

// Sqrt finds the square root of a posit
// p.Sqrt() returns a new posit, p is not altered
func (p Posit16) Sqrt() Posit16 { return Posit16{impl: p.impl.Sqrt()} }

// Size specific

// FromInt creates a new posit which is set to the value of a integer
// p.FromInt() outputs a new posit z, p is not altered
func (p Posit16) FromInt(i int16) Posit16 { return Posit16{impl: p.impl.FromInt(int64(i))} }

// FromUint creates a new posit which is set to the value of an unsigned integer
// p.FromUint() outputs a new posit z, p is not altered
func (p Posit16) FromUint(i uint16) Posit16 { return Posit16{impl: p.impl.FromUint(uint64(i))} }

// Int outputs an int16 representation of the value of the posit
// lost bits are rounded to nearest even and if the number is greater than or equal
// to two to the 16 -1 power, the maximum 0x7fff for positive input or
// 0xffff for negative input will return.
// Likewise if the number is less than or equal to 1/2 it will round down to zero.
func (p Posit16) Int() int16 {
	x := p.impl.Int()
	if x > 0x7fff {
		return 0x7fff
	}
	if x < -0x7fff {
		return -0x7fff
	}
	return int16(x)
}

// Uint outputs a uint16 representation of the value of the posit, if the posit
// is less than or equal to 1/2, 0 will be returned and if the value is greater
// than two to the maxBits power, 0xffffffffffffffff will be returned.
func (p Posit16) Uint() uint16 {
	x := p.impl.Int()

	if x > 0xffff {
		return 0xffff
	}

	return uint16(x)
}

// Exp outputs the exponent of the posit, specifically the number z for which
// 0.5*2**z <= positValue < 1*2**z
func (p Posit16) Exp() int16 {
	x := p.impl.Exp()

	if x > 0x7fff || x < -0x7fff {
		panic("posit with bigger exponent than should be possible")
	}

	return int16(x)
}

// Mant outputs a posit with the same mantissa but an exponent of 0, that is a
// number which is greater than or equal to 0.5 and less than 1
func (p Posit16) Mant() Posit16 { return Posit16{impl: p.impl.Mant()} }

// ExpAdd returns a new posit with x added to the exponent, effectively multiplying the
// posit value by 2**x. It is equivilent to a bit-shift in integer math. If bits is
// negative then the exponent will be decreased, if bits causes the exponent to increase
// in magnitude, it might cause the posit to round.
func (p Posit16) ExpAdd(x int16) Posit16 { return Posit16{impl: p.impl.ExpAdd(int32(x))} }

// Up returns the same number up-casted to the next larger posit, in this case larger means
// twice the nbits and an exponent size which is one greater.
func (p Posit16) Up() Posit32 {
	return Posit32{impl: p.impl.Up()}
}

// Down returns the same number down-casted to the next smaller posit, see Up() for the
// definition of larger and smaller.
func (p Posit16) Down() Posit8 {
	return Posit8{impl: p.impl.Up()}
}

// Bits outputs a uint16 containing the raw binary format of the posit
func (p Posit16) Bits() uint16 {
	words := p.impl.Bits.Bits()
	if len(words) == 1 {
		return uint16(words[0])
	} else if len(words) == 2 {
		if bits.UintSize != 32 {
			panic("unexpected number of bits")
		}
		return (uint16(words[1]) << (16 / 2)) | uint16(words[0])
	}
	panic("unexpected word size")
}

// SetBits outputs a new posit with the bits set to those which you specify
func (p Posit16) SetBits(bits uint16) Posit16 {
	bi := new(big.Int).SetUint64(uint64(bits))
	return Posit16{impl: &SlowPosit{es: p.impl.es, nbits: p.impl.nbits, Bits: bi}}
}

// Clone makes a copy of a posit
func (p Posit16) Clone() Posit16 {
	return Posit16{impl: p.impl.Clone()}
}

// Posit32 is an 32 bit posit with 2 exponent bits
type Posit32 struct{ impl *SlowPosit }

// NewPosit32 makes a new posit with 32 bits and 2 es bits, the initial value is zero
func NewPosit32() Posit32 { return Posit32{impl: NewSlowPosit(32, 2)} }

// Add takes the sum of two posits
// p.Add(x) creates a new posit z which is p+x, neither p nor x are altered
func (p Posit32) Add(x Posit32) Posit32 { return Posit32{impl: p.impl.Add(x.impl)} }

// AddExact returns exactly the sum of two posits, represented as two
// more posits, the first result is a posit which is the sum truncated
// (round to negative infinity) and the second posit is the difference between
// the first posit and the true sum. That is to say true_sum = out2 + out1.
func (p Posit32) AddExact(x Posit32) (Posit32, Posit32) {
	res, diff := p.impl.AddExact(x.impl)
	return Posit32{impl: res}, Posit32{impl: diff}
}

// Sub takes the difference of two posits
// p.Sub(x) creates a new posit z which is p-x, neither p nor x are altered
func (p Posit32) Sub(x Posit32) Posit32 { return Posit32{impl: p.impl.Sub(x.impl)} }

// SubExact returns exactly the difference of two posits, represented as two
// more posits, the first result is a posit which is the difference truncated
// (round to negative infinity) and the second posit is the difference between
// the first posit and the true sum. That is to say true_sum = out2 + out1.
func (p Posit32) SubExact(x Posit32) (Posit32, Posit32) {
	res, diff := p.impl.SubExact(x.impl)
	return Posit32{impl: res}, Posit32{impl: diff}
}

// Mul takes the product of two posits
// p.Mul(x) creates a new posit z which is x*y, neither p nor x are altered
func (p Posit32) Mul(x Posit32) Posit32 { return Posit32{impl: p.impl.Mul(x.impl)} }

// MulPromote takes the product of two posits
// p.Mul(x) creates a new posit z which is p*x represented as the next larger posit size
// next larger means the bit width is doubled and the exponent size is increased by 1
// this function is guaranteed not to round
func (p Posit32) MulPromote(x Posit32) Posit64 { return Posit64{impl: p.impl.MulPromote(x.impl)} }

// Div takes the quotent of two posits
// p.Div(x) creates a new posit z which is p/x, neither p nor x are altered
func (p Posit32) Div(x Posit32) Posit32 { return Posit32{impl: p.impl.Div(x.impl)} }

// DivPromote takes the quotent of two posits
// p.DivPromote(x) creates a new posit z which is p/x, as the next larger posit size
// next larger means the bit width is doubled and the exponent size is increased by 1
// this function is guaranteed not to round
func (p Posit32) DivPromote(x Posit32) Posit64 { return Posit64{impl: p.impl.DivPromote(x.impl)} }

// Sqrt finds the square root of a posit
// p.Sqrt() returns a new posit, p is not altered
func (p Posit32) Sqrt() Posit32 { return Posit32{impl: p.impl.Sqrt()} }

// Size specific

// FromInt creates a new posit which is set to the value of a integer
// p.FromInt() outputs a new posit z, p is not altered
func (p Posit32) FromInt(i int32) Posit32 { return Posit32{impl: p.impl.FromInt(int64(i))} }

// FromUint creates a new posit which is set to the value of an unsigned integer
// p.FromUint() outputs a new posit z, p is not altered
func (p Posit32) FromUint(i uint32) Posit32 { return Posit32{impl: p.impl.FromUint(uint64(i))} }

// Int outputs an int32 representation of the value of the posit
// lost bits are rounded to nearest even and if the number is greater than or equal
// to two to the 32 -1 power, the maximum 0x7fffffff for positive input or
// 0xffffffff for negative input will return.
// Likewise if the number is less than or equal to 1/2 it will round down to zero.
func (p Posit32) Int() int32 {
	x := p.impl.Int()
	if x > 0x7fffffff {
		return 0x7fffffff
	}
	if x < -0x7fffffff {
		return -0x7fffffff
	}
	return int32(x)
}

// Uint outputs a uint32 representation of the value of the posit, if the posit
// is less than or equal to 1/2, 0 will be returned and if the value is greater
// than two to the maxBits power, 0xffffffffffffffff will be returned.
func (p Posit32) Uint() uint32 {
	x := p.impl.Int()

	if x > 0xffffffff {
		return 0xffffffff
	}

	return uint32(x)
}

// Exp outputs the exponent of the posit, specifically the number z for which
// 0.5*2**z <= positValue < 1*2**z
func (p Posit32) Exp() int32 {
	x := p.impl.Exp()

	if x > 0x7fffffff || x < -0x7fffffff {
		panic("posit with bigger exponent than should be possible")
	}

	return int32(x)
}

// Mant outputs a posit with the same mantissa but an exponent of 0, that is a
// number which is greater than or equal to 0.5 and less than 1
func (p Posit32) Mant() Posit32 { return Posit32{impl: p.impl.Mant()} }

// ExpAdd returns a new posit with x added to the exponent, effectively multiplying the
// posit value by 2**x. It is equivilent to a bit-shift in integer math. If bits is
// negative then the exponent will be decreased, if bits causes the exponent to increase
// in magnitude, it might cause the posit to round.
func (p Posit32) ExpAdd(x int32) Posit32 { return Posit32{impl: p.impl.ExpAdd(int32(x))} }

// Up returns the same number up-casted to the next larger posit, in this case larger means
// twice the nbits and an exponent size which is one greater.
func (p Posit32) Up() Posit64 {
	return Posit64{impl: p.impl.Up()}
}

// Down returns the same number down-casted to the next smaller posit, see Up() for the
// definition of larger and smaller.
func (p Posit32) Down() Posit16 {
	return Posit16{impl: p.impl.Up()}
}

// Bits outputs a uint32 containing the raw binary format of the posit
func (p Posit32) Bits() uint32 {
	words := p.impl.Bits.Bits()
	if len(words) == 1 {
		return uint32(words[0])
	} else if len(words) == 2 {
		if bits.UintSize != 32 {
			panic("unexpected number of bits")
		}
		return (uint32(words[1]) << (32 / 2)) | uint32(words[0])
	}
	panic("unexpected word size")
}

// SetBits outputs a new posit with the bits set to those which you specify
func (p Posit32) SetBits(bits uint32) Posit32 {
	bi := new(big.Int).SetUint64(uint64(bits))
	return Posit32{impl: &SlowPosit{es: p.impl.es, nbits: p.impl.nbits, Bits: bi}}
}

// Clone makes a copy of a posit
func (p Posit32) Clone() Posit32 {
	return Posit32{impl: p.impl.Clone()}
}

// Posit64 is an 64 bit posit with 3 exponent bits
type Posit64 struct{ impl *SlowPosit }

// NewPosit64 makes a new posit with 64 bits and 3 es bits, the initial value is zero
func NewPosit64() Posit64 { return Posit64{impl: NewSlowPosit(64, 3)} }

// Add takes the sum of two posits
// p.Add(x) creates a new posit z which is p+x, neither p nor x are altered
func (p Posit64) Add(x Posit64) Posit64 { return Posit64{impl: p.impl.Add(x.impl)} }

// AddExact returns exactly the sum of two posits, represented as two
// more posits, the first result is a posit which is the sum truncated
// (round to negative infinity) and the second posit is the difference between
// the first posit and the true sum. That is to say true_sum = out2 + out1.
func (p Posit64) AddExact(x Posit64) (Posit64, Posit64) {
	res, diff := p.impl.AddExact(x.impl)
	return Posit64{impl: res}, Posit64{impl: diff}
}

// Sub takes the difference of two posits
// p.Sub(x) creates a new posit z which is p-x, neither p nor x are altered
func (p Posit64) Sub(x Posit64) Posit64 { return Posit64{impl: p.impl.Sub(x.impl)} }

// SubExact returns exactly the difference of two posits, represented as two
// more posits, the first result is a posit which is the difference truncated
// (round to negative infinity) and the second posit is the difference between
// the first posit and the true sum. That is to say true_sum = out2 + out1.
func (p Posit64) SubExact(x Posit64) (Posit64, Posit64) {
	res, diff := p.impl.SubExact(x.impl)
	return Posit64{impl: res}, Posit64{impl: diff}
}

// Mul takes the product of two posits
// p.Mul(x) creates a new posit z which is x*y, neither p nor x are altered
func (p Posit64) Mul(x Posit64) Posit64 { return Posit64{impl: p.impl.Mul(x.impl)} }

// MulPromote takes the product of two posits
// p.Mul(x) creates a new posit z which is p*x represented as the next larger posit size
// next larger means the bit width is doubled and the exponent size is increased by 1
// this function is guaranteed not to round
func (p Posit64) MulPromote(x Posit64) Posit128 { return Posit128{impl: p.impl.MulPromote(x.impl)} }

// Div takes the quotent of two posits
// p.Div(x) creates a new posit z which is p/x, neither p nor x are altered
func (p Posit64) Div(x Posit64) Posit64 { return Posit64{impl: p.impl.Div(x.impl)} }

// DivPromote takes the quotent of two posits
// p.DivPromote(x) creates a new posit z which is p/x, as the next larger posit size
// next larger means the bit width is doubled and the exponent size is increased by 1
// this function is guaranteed not to round
func (p Posit64) DivPromote(x Posit64) Posit128 { return Posit128{impl: p.impl.DivPromote(x.impl)} }

// Sqrt finds the square root of a posit
// p.Sqrt() returns a new posit, p is not altered
func (p Posit64) Sqrt() Posit64 { return Posit64{impl: p.impl.Sqrt()} }

// Size specific

// FromInt creates a new posit which is set to the value of a integer
// p.FromInt() outputs a new posit z, p is not altered
func (p Posit64) FromInt(i int64) Posit64 { return Posit64{impl: p.impl.FromInt(int64(i))} }

// FromUint creates a new posit which is set to the value of an unsigned integer
// p.FromUint() outputs a new posit z, p is not altered
func (p Posit64) FromUint(i uint64) Posit64 { return Posit64{impl: p.impl.FromUint(uint64(i))} }

// Int outputs an int64 representation of the value of the posit
// lost bits are rounded to nearest even and if the number is greater than or equal
// to two to the 64 -1 power, the maximum 0x7fffffffffffffff for positive input or
// 0xffffffffffffffff for negative input will return.
// Likewise if the number is less than or equal to 1/2 it will round down to zero.
func (p Posit64) Int() int64 {
	x := p.impl.Int()
	if x > 0x7fffffffffffffff {
		return 0x7fffffffffffffff
	}
	if x < -0x7fffffffffffffff {
		return -0x7fffffffffffffff
	}
	return int64(x)
}

// Uint outputs a uint64 representation of the value of the posit, if the posit
// is less than or equal to 1/2, 0 will be returned and if the value is greater
// than two to the maxBits power, 0xffffffffffffffff will be returned.
func (p Posit64) Uint() uint64 {
	x := p.impl.Int()

	return uint64(x)
}

// Exp outputs the exponent of the posit, specifically the number z for which
// 0.5*2**z <= positValue < 1*2**z
func (p Posit64) Exp() int64 {
	x := p.impl.Exp()

	return int64(x)
}

// Mant outputs a posit with the same mantissa but an exponent of 0, that is a
// number which is greater than or equal to 0.5 and less than 1
func (p Posit64) Mant() Posit64 { return Posit64{impl: p.impl.Mant()} }

// ExpAdd returns a new posit with x added to the exponent, effectively multiplying the
// posit value by 2**x. It is equivilent to a bit-shift in integer math. If bits is
// negative then the exponent will be decreased, if bits causes the exponent to increase
// in magnitude, it might cause the posit to round.
func (p Posit64) ExpAdd(x int64) Posit64 { return Posit64{impl: p.impl.ExpAdd(int32(x))} }

// Up returns the same number up-casted to the next larger posit, in this case larger means
// twice the nbits and an exponent size which is one greater.
func (p Posit64) Up() Posit128 {
	return Posit128{impl: p.impl.Up()}
}

// Down returns the same number down-casted to the next smaller posit, see Up() for the
// definition of larger and smaller.
func (p Posit64) Down() Posit32 {
	return Posit32{impl: p.impl.Up()}
}

// Bits outputs a uint64 containing the raw binary format of the posit
func (p Posit64) Bits() uint64 {
	words := p.impl.Bits.Bits()
	if len(words) == 1 {
		return uint64(words[0])
	} else if len(words) == 2 {
		if bits.UintSize != 32 {
			panic("unexpected number of bits")
		}
		return (uint64(words[1]) << (64 / 2)) | uint64(words[0])
	}
	panic("unexpected word size")
}

// SetBits outputs a new posit with the bits set to those which you specify
func (p Posit64) SetBits(bits uint64) Posit64 {
	bi := new(big.Int).SetUint64(uint64(bits))
	return Posit64{impl: &SlowPosit{es: p.impl.es, nbits: p.impl.nbits, Bits: bi}}
}

// Clone makes a copy of a posit
func (p Posit64) Clone() Posit64 {
	return Posit64{impl: p.impl.Clone()}
}
