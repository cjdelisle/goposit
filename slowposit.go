package goposit

import (
	"fmt"
	"math/big"
	"math/bits"
)

var bigf1 = big.NewFloat(1)
var bigi1 = big.NewInt(1)

func printFloat(f *big.Float) string {
	return f.Text('p', 40)
}

func assertExact(f *big.Float) {
	if f.Acc() != big.Exact {
		fmt.Printf("Float [%v] is not exact", printFloat(f))
		panic("float is not exact")
	}
}

// Shf is a utility function for "bit shifting" (multiplying by 2**i) a float
func Shf(x *big.Float, s int) {
	exp := x.MantExp(x)
	x.SetMantExp(x, exp+s)
	assertExact(x)
}

func checkAdd(f *big.Float, x *big.Float) {
	f.Add(f, x)
	assertExact(f)
}

func checkSub(f *big.Float, x *big.Float) {
	f.Sub(f, x)
	assertExact(f)
}

func popcnt(i *big.Int) int {
	cnt := 0
	for _, w := range i.Bits() {
		cnt += bits.OnesCount(uint(w))
	}
	return cnt
}

func negate(out *big.Int, size uint) {
	ones := big.NewInt(0)
	ones.SetBit(ones, int(size), 1)
	ones.Sub(ones, bigi1)

	out.Xor(out, ones)
	out.Add(out, bigi1)
	out.And(out, ones)
}

// SlowPosit is a structure which contains a float-like number
type SlowPosit struct {
	Bits  *big.Int
	nbits uint
	es    uint
}

// Nbits returns the bit width of the posit
func (p *SlowPosit) Nbits() uint { return p.nbits }

// Es returns exponent size
func (p *SlowPosit) Es() uint { return p.es }

// NewSlowPosit creates a bigint backed posit which converts everything to big.Float
func NewSlowPosit(nbits uint, es uint) *SlowPosit {
	if nbits <= es+3 {
		panic("es must be more than 3 bits less than nbits")
	}
	return &SlowPosit{nbits: nbits, es: es, Bits: big.NewInt(0)}
}

func (p *SlowPosit) signBitIdx() int { return int(p.nbits - 1) }
func (p *SlowPosit) regBitIdx() int  { return int(p.nbits - 2) }

func (p *SlowPosit) Clone() *SlowPosit {
	return &SlowPosit{
		nbits: p.nbits,
		es:    p.es,
		Bits:  p.Bits,
	}
}

func (p *SlowPosit) Log2MaxVal() int {
	return int(p.nbits*(1<<p.es)) - int(1<<(p.es+1))
}

func (p *SlowPosit) NaR() *SlowPosit {
	v := big.NewInt(0)
	v.SetBit(v, int(p.signBitIdx()), 1)
	p.Bits = v
	return p
}

func (p *SlowPosit) Zero() *SlowPosit {
	p.Bits = big.NewInt(0)
	return p
}

func (p *SlowPosit) One() *SlowPosit {
	v := big.NewInt(0)
	v.SetBit(v, int(p.regBitIdx()), 1)
	p.Bits = v
	return p
}

func (p *SlowPosit) NegOne() *SlowPosit {
	v := big.NewInt(0)
	v.SetBit(v, int(p.signBitIdx()), 1)
	v.SetBit(v, int(p.regBitIdx()), 1)
	p.Bits = v
	return p
}

func (p *SlowPosit) Max() *SlowPosit {
	v := big.NewInt(0)
	v.SetBit(v, int(p.signBitIdx()), 1)
	v.Sub(v, bigi1)
	p.Bits = v
	return p
}

func (p *SlowPosit) MaxNeg() *SlowPosit {
	v := big.NewInt(1)
	v.SetBit(v, int(p.signBitIdx()), 1)
	p.Bits = v
	return p
}

func (p *SlowPosit) Min() *SlowPosit {
	p.Bits = big.NewInt(1)
	return p
}

func (p *SlowPosit) MinNeg() *SlowPosit {
	v := big.NewInt(0)
	v.SetBit(v, int(p.nbits), 1)
	v.Sub(v, bigi1)
	p.Bits = v
	return p
}

func (p *SlowPosit) RawHex() string {
	return p.Bits.Text(16)
}

func (p *SlowPosit) Cmp(p2 *SlowPosit) int {
	if p2.nbits != p.nbits || p2.es != p.es {
		panic("can't compare posits of different sizes")
	}
	return p.Bits.Cmp(p2.Bits)
}

// IsNaR returns true if the posit is a NaR value
func (p *SlowPosit) IsNaR() bool {
	if popcnt(p.Bits) != 1 {
		return false
	}
	return p.Bits.Bit(int(p.nbits-1)) == 1
}

// Uint64 returns the posit as a raw uint64
func (p *SlowPosit) Uint64() uint64 {
	return p.Bits.Uint64()
}

func (p *SlowPosit) SetBits(bits *big.Int) {
	p.Bits = bits
}

// FromFloat sets a posit from a floating point number as input
// The old value of the posit is lost
// this function returns true if the floating point number could be converted with no rounding
func (p *SlowPosit) FromFloat(input *big.Float, truncate bool) bool {
	if input.IsInf() {
		p.NaR()
		return false
	}
	sign := input.Sign()
	if sign == 0 {
		p.Bits = big.NewInt(0)
		return true
	}
	negative := sign < 0

	exponent := input.MantExp(nil)
	lessThan1 := exponent < 1
	scale := uint(exponent - 1)
	if lessThan1 {
		scale = uint(-exponent)
	}
	regimeBits := scale >> p.es
	subexponent := scale % (1 << p.es)
	if lessThan1 {
		subexponent = (1 << p.es) - 1 - subexponent
	}

	// Start with bigfloat 1, this is putting a 1 in the sign bit place
	// We're going to remove it later but it gives us a marker of the top
	// of the float so we don't need to rely on the exponent
	f := big.NewFloat(1)
	f.SetPrec((2 * p.nbits) - 1)
	f.SetMode(big.ToZero)
	assertExact(f)

	// create the regime
	inexact := false
	for i := uint(0); i < regimeBits+1; i++ {
		if i >= p.nbits-1 {
			inexact = true
			break
		}
		Shf(f, 1)
		if !lessThan1 {
			checkAdd(f, bigf1)
		}
	}

	// guard bit
	Shf(f, 1)
	if lessThan1 {
		checkAdd(f, bigf1)
	}

	// exponent
	Shf(f, int(p.es))
	checkAdd(f, big.NewFloat(float64(subexponent)))

	// subtract 1 because when we add the input, we want to lose the hidden bit
	checkSub(f, bigf1)

	// shift f so that it's 1's place lines up with the top bit of the input
	Shf(f, exponent-1)

	// Round 1
	if negative {
		f.Sub(f, input)
	} else {
		f.Add(f, input)
	}
	if f.Acc() != big.Exact {
		inexact = true
		// set the accuracy back to Exact
		f.SetMode(f.Mode())
	}

	// Put the sign bit just under the decimal point
	f.MantExp(f)
	assertExact(f)

	// Shift up by 2*Posit_SZ so that the 1's place is 1 bit below the bottom of the input
	// (recall the input was limited to 2*Posit_SZ-1 bits of precision)
	f.SetMantExp(f, int(p.nbits*2))
	assertExact(f)

	// Bump the precision up by 1 bit
	f.SetPrec(2 * p.nbits)
	assertExact(f)

	// Add an "inexact bit" which will be a tie-breaker, just as would have been the bits
	// that are now rounded off, we have a percision of 3*nbits
	if inexact {
		checkAdd(f, bigf1)
	}

	// Round 2 (fight!)
	// This time we round-to-nearest
	if truncate {
		// If we're asked to truncate the number, we should round to negative inf
		// because the number is inverted at the end.
		f.SetMode(big.ToNegativeInf)
	} else {
		f.SetMode(big.ToNearestEven)
	}
	f.SetPrec(p.nbits)
	if f.Acc() != big.Exact {
		inexact = true
		// set the accuracy back to Exact
		f.SetMode(f.Mode())
	}

	// Shift down so that the 1's place is the bottom bit of the actual output posit
	Shf(f, -int(p.nbits))

	out, accuracy := f.Int(nil)
	if accuracy != big.Exact {
		panic("conversion to int lost bits")
	}

	// If the sign bit is not set, something went wrong
	if out.Bit(int(p.nbits-1)) != 1 {
		panic("sign bit was not set")
	}

	// clear the sign bit
	out.SetBit(out, int(p.nbits-1), 0)

	if negative {
		negate(out, p.nbits)
	}

	p.Bits = out
	return !inexact
}

// ToFloat takes a posit and outputs a big.Float
// conversion is always exact
func (p *SlowPosit) ToFloat() *big.Float {
	// NaR
	if p.IsNaR() {
		out := big.NewFloat(0)
		out.SetInf(true)
		return out
	}

	// 0
	if p.Bits.Sign() == 0 {
		return big.NewFloat(0)
	}

	i1 := big.NewInt(0)
	i1.Set(p.Bits)

	negative := i1.Bit(p.signBitIdx()) == 1
	if negative {
		negate(i1, p.nbits)
	}

	lessThan1 := i1.Bit(p.regBitIdx()) == 0

	guardBitIdx := p.regBitIdx() - 1
	regimeBits := uint(0)
	for ; guardBitIdx >= 0; guardBitIdx-- {
		if (i1.Bit(guardBitIdx) == 1) == lessThan1 {
			break
		}
		regimeBits++
	}

	if guardBitIdx >= 0 {
		// mask off the regime
		i2 := big.NewInt(0)
		i2.SetBit(i2, guardBitIdx, 1)
		i2.Sub(i2, bigi1)
		i1.And(i1, i2)

		// but we're going to flag the guard bit itself so that
		// we can find the top of the subexponent
		i1.SetBit(i1, guardBitIdx, 1)
	} else {
		i1.SetUint64(0)
	}

	f1 := big.NewFloat(0)
	neededPrecision := (guardBitIdx / 64) * 64
	if neededPrecision < guardBitIdx {
		neededPrecision += 64
	}
	if f1.Prec() < uint(neededPrecision) {
		// Make sure we have enough precision
		f1.SetPrec(uint(neededPrecision))
	}
	for _, w := range i1.Bits() {
		fw := big.NewFloat(0)
		fw.SetPrec(bits.UintSize)
		fw.SetUint64(uint64(w))
		assertExact(fw)
		checkAdd(f1, fw)
		Shf(f1, bits.UintSize)
	}

	// Move the guard bit up to just below the decimal point...
	f1.MantExp(f1)
	assertExact(f1)

	// Grab off the subexponent and the guard bit
	Shf(f1, int(p.es+1))
	subexp64, _ := f1.Int64()
	checkSub(f1, big.NewFloat(0).SetInt64(subexp64))
	subexp := uint(uint64(subexp64) % (1 << p.es))

	// Add the hidden bit
	checkAdd(f1, bigf1)
	Shf(f1, -1)

	// calculate the exponent
	exponent := 0
	if lessThan1 {
		exponent = 1 - int((regimeBits+1)<<p.es) + int(subexp%(uint(1)<<p.es))
	} else {
		exponent = 1 + int(regimeBits<<p.es) + int(subexp%(uint(1)<<p.es))
	}

	f1.SetMantExp(f1, exponent)
	assertExact(f1)

	if negative {
		f1.Neg(f1)
	}

	return f1
}

func assertCompat(x, y *SlowPosit) {
	if x.nbits != y.nbits || x.es != y.es {
		panic("Attempted to do math on posits with different parameters")
	}
}

func getFloats(x, y *SlowPosit) (xf, yf *big.Float) {
	assertCompat(x, y)
	xf = x.ToFloat()
	yf = y.ToFloat()
	xf.SetPrec(x.nbits * 2)
	yf.SetPrec(x.nbits * 2)
	xf.SetMode(big.ToNearestEven)
	yf.SetMode(big.ToNearestEven)
	return
}

func outPosit(template *SlowPosit, f *big.Float) *SlowPosit {
	out := SlowPosit{nbits: template.nbits, es: template.es}
	out.FromFloat(f, false)
	return &out
}

func addExact(p *SlowPosit, x *big.Float, y *big.Float) (*SlowPosit, *SlowPosit) {
	// We can have loss of precision during the add OR during the serialziation
	z := new(big.Float)
	z.SetMode(big.ToNegativeInf)
	z.Add(x, y)
	zp := SlowPosit{nbits: p.nbits, es: p.es}
	zp.FromFloat(z, true)

	z2 := zp.ToFloat()
	z2.SetMode(big.ToNegativeInf)
	if x.Cmp(y) > 0 {
		z2.Sub(z2, x)
		z2.Sub(z2, y)
	} else {
		z2.Sub(z2, y)
		z2.Sub(z2, x)
	}
	remp := SlowPosit{nbits: p.nbits, es: p.es}
	if exact := remp.FromFloat(z, false); !exact {
		panic("inexact conversion when using addExact")
	}
	return &zp, &remp
}

/// math ////

// Add takes the sum of two posits
// p.Add(x) creates a new posit z which is p+x, neither p nor x are altered
func (p *SlowPosit) Add(x *SlowPosit) *SlowPosit {
	pf, xf := getFloats(p, x)
	xf.Add(pf, xf)
	return outPosit(p, xf)
}

// AddExact returns exactly the sum of two posits, represented as two
// more posits, the first result is a posit which is the sum truncated
// (round to negative infinity) and the second posit is the difference between
// the first posit and the true sum. That is to say true_sum = out2 + out1.
func (p *SlowPosit) AddExact(x *SlowPosit) (*SlowPosit, *SlowPosit) {
	pf, xf := getFloats(p, x)
	return addExact(p, pf, xf)
}

// Sub takes the difference of two posits
// p.Sub(x) creates a new posit z which is p-x, neither p nor x are altered
func (p *SlowPosit) Sub(x *SlowPosit) *SlowPosit {
	pf, xf := getFloats(p, x)
	xf.Sub(pf, xf)
	return outPosit(p, xf)
}

// SubExact returns exactly the difference of two posits, represented as two
// more posits, the first result is a posit which is the difference truncated
// (round to negative infinity) and the second posit is the difference between
// the first posit and the true sum. That is to say true_sum = out2 + out1.
func (p *SlowPosit) SubExact(x *SlowPosit) (*SlowPosit, *SlowPosit) {
	pf, xf := getFloats(p, x)
	xf.Neg(xf)
	return addExact(p, pf, xf)
}

// Mul takes the product of two posits
// p.Mul(x) creates a new posit z which is x*y, neither p nor x are altered
func (p *SlowPosit) Mul(x *SlowPosit) *SlowPosit {
	pf, xf := getFloats(p, x)
	xf.Mul(pf, xf)
	return outPosit(p, xf)
}

// MulPromote takes the product of two posits
// p.Mul(x) creates a new posit z which is p*x represented as the next larger posit size
// next larger means the bit width is doubled and the exponent size is increased by 1
// this function is guaranteed not to round
func (p *SlowPosit) MulPromote(x *SlowPosit) *SlowPosit {
	pf, xf := getFloats(p, x)
	pf.SetPrec(x.nbits * 4)
	xf.SetPrec(x.nbits * 4)
	xf.Mul(xf, pf)
	bigger := &SlowPosit{nbits: p.nbits * 2, es: p.es + 1}
	return outPosit(bigger, xf)
}

// Div takes the quotent of two posits
// p.Div(x) creates a new posit z which is p/x, neither p nor x are altered
func (p *SlowPosit) Div(x *SlowPosit) *SlowPosit {
	pf, xf := getFloats(p, x)
	xf.Quo(pf, xf)
	return outPosit(p, xf)
}

// DivPromote takes the quotent of two posits
// p.DivPromote(x) creates a new posit z which is p/x, as the next larger posit size
// next larger means the bit width is doubled and the exponent size is increased by 1
// this function is guaranteed not to round
func (p *SlowPosit) DivPromote(x *SlowPosit) *SlowPosit {
	pf, xf := getFloats(p, x)
	pf.SetPrec(x.nbits * 4)
	xf.SetPrec(x.nbits * 4)
	xf.Quo(xf, pf)
	bigger := &SlowPosit{nbits: p.nbits * 2, es: p.es + 1}
	return outPosit(bigger, xf)
}

// Sqrt finds the square root of a posit
// p.Sqrt() returns a new posit, p is not altered
func (p *SlowPosit) Sqrt() *SlowPosit {
	xf := p.ToFloat()
	xf.SetPrec(p.nbits * 2)
	xf.Sqrt(xf)
	return outPosit(p, xf)
}

//// conversions ////

// FromInt creates a new posit which is set to the value of a integer
// p.FromInt() outputs a new posit z, p is not altered
func (p *SlowPosit) FromInt(i int64) *SlowPosit {
	f := big.NewFloat(0)
	f.SetPrec(64)
	f.SetInt64(i)
	return outPosit(p, f)
}

// FromUint creates a new posit which is set to the value of an unsigned integer
// p.FromUint() outputs a new posit z, p is not altered
func (p *SlowPosit) FromUint(i uint64) *SlowPosit {
	f := big.NewFloat(0)
	f.SetPrec(64)
	f.SetUint64(i)
	return outPosit(p, f)
}

// Int outputs an int64 representation of the value of the posit
// lost bits are rounded to nearest even and if the number is greater than or equal
// to two to the 63rd power, the maximum 0x7fffffffffffffff for positive input or
// 0xffffffffffffffff for negative input will return.
// Likewise if the number is less than or equal to 1/2 it will round down to zero.
func (p *SlowPosit) Int() (out int64) {
	f := p.ToFloat()
	if f.IsInf() {
		return 0x7fffffffffffffff
	}
	if f.Sign() == 0 {
		return 0
	}
	exp := f.MantExp(nil)
	if exp >= 64 {
		out = 0x7fffffffffffffff
		if f.Sign() < 0 {
			out = -out
		}
		return
	}
	twoToThe64 := big.NewFloat(18446744073709551616.0)
	if f.Sign() < 0 {
		twoToThe64 = twoToThe64.Neg(twoToThe64)
	}
	f.Add(f, twoToThe64)
	if f.MantExp(nil) != 64 {
		panic("internal: exponent calculation error")
	}

	f.SetMode(big.ToNearestEven)
	f.SetPrec(64)

	f.Sub(f, twoToThe64)
	out, acc := f.Int64()
	if acc != big.Exact {
		panic("Accuracy should be Exact")
	}
	return
}

// Uint outputs a uint64 representation of the value of the posit, if the posit
// is less than or equal to 1/2, 0 will be returned and if the value is greater
// than two to the 64 power, 0xffffffffffffffff will be returned.
func (p *SlowPosit) Uint() (out uint64) {
	f := p.ToFloat()
	if f.IsInf() {
		return 0xffffffffffffffff
	}
	if f.Sign() == 0 {
		return 0
	}
	exp := f.MantExp(nil)
	if exp > 64 {
		out = 0xffffffffffffffff
		return
	}
	twoToThe65 := big.NewFloat(36893488147419103232.0)
	if f.Sign() < 0 {
		twoToThe65 = twoToThe65.Neg(twoToThe65)
	}
	f.Add(f, twoToThe65)
	if f.MantExp(nil) != 65 {
		panic("internal: exponent calculation error")
	}

	f.SetMode(big.ToNearestEven)
	f.SetPrec(65)

	f.Sub(f, twoToThe65)
	out, acc := f.Uint64()
	if acc != big.Exact {
		panic("Accuracy should be Exact")
	}
	return
}

//// binary manipulation ////

// Exp outputs the exponent of the posit, specifically the number z for which
// 0.5*2**z <= positValue < 1*2**z
func (p *SlowPosit) Exp() int32 {
	f := p.ToFloat()
	return int32(f.MantExp(nil))
}

// Mant outputs a posit with the same mantissa but an exponent of 0, that is a
// number which is greater than or equal to 0.5 and less than 1
func (p *SlowPosit) Mant() *SlowPosit {
	f := p.ToFloat()
	f.MantExp(f)
	return outPosit(p, f)
}

// ExpAdd returns a new posit with x added to the exponent, effectively multiplying the
// posit value by 2**x. It is equivilent to a bit-shift in integer math. If bits is
// negative then the exponent will be decreased, if bits causes the exponent to increase
// in magnitude, it might cause the posit to round.
func (p *SlowPosit) ExpAdd(x int32) *SlowPosit {
	f := p.ToFloat()
	Shf(f, int(x))
	return outPosit(p, f)
}

// casting

// Up returns the same number up-casted to the next larger posit, in this case larger means
// twice the nbits and an exponent size which is one greater.
func (p *SlowPosit) Up() *SlowPosit {
	f := p.ToFloat()
	bigger := SlowPosit{nbits: p.nbits * 2, es: p.es + 1}
	return outPosit(&bigger, f)
}

// Down returns the same number down-casted to the next smaller posit, see Up() for the
// definition of larger and smaller.
func (p *SlowPosit) Down() *SlowPosit {
	if p.es < 1 {
		panic("There is no smaller posit size")
	}
	f := p.ToFloat()
	smaller := SlowPosit{nbits: p.nbits / 2, es: p.es - 1}
	return outPosit(&smaller, f)
}
