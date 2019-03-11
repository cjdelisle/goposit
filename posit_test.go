package goposit_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/cjdelisle/goposit"
)

func assert(x bool) {
	if !x {
		panic("assertion failure")
	}
}

func assertSame(t *testing.T, a, b *big.Int, name string) {
	if a.Cmp(b) != 0 {
		t.Errorf("%v test: resulting posits differ [%v] [%v]", name, a.Text(16), b.Text(16))
	}
}

func assertSameF(t *testing.T, a, b *big.Float, name string) {
	if a.Cmp(b) != 0 {
		t.Errorf("%v test: resulting floats differ", name)
		//t.Errorf("%v test: resulting floats differ [%v] [%v]", name, a, b)
	}
}

func p2f(t *testing.T, positClass *goposit.SlowPosit, pval *big.Int) *big.Float {
	b := positClass.Bits
	positClass.Bits = pval
	out := positClass.ToFloat()
	positClass.Bits = b
	if out.Acc() != big.Exact {
		t.Errorf("posit.ToFloat() returns an inexact float")
	}
	return out
}

func f2p(t *testing.T, positClass *goposit.SlowPosit, f *big.Float) (*big.Int, bool) {
	b := positClass.Bits
	f2 := big.NewFloat(0).Copy(f)
	isExact := positClass.FromFloat(f, false)
	if f.Cmp(f2) != 0 {
		t.Errorf("Posit.FromFloat modified it's input")
	}
	out := positClass.Bits
	positClass.Bits = b
	if isExact {
		f2 := p2f(t, positClass, out)
		assertSameF(t, f2, f, fmt.Sprintf("inverting %v (posit: [%v])", f, out.Text(16)))
	}
	return out, isExact
}

func f2pExact(t *testing.T, positClass *goposit.SlowPosit, f *big.Float, expectExact bool) *big.Int {
	out, exact := f2p(t, positClass, f)
	if exact != expectExact {
		t.Errorf("Expected exact conversion to be %v but it was %v", expectExact, exact)
	}
	return out
}

func conversionTest(t *testing.T, f *big.Float, expect *goposit.SlowPosit, name string, exact bool) {
	//p := expect.Clone()
	pint, isExact := f2p(t, expect, f)
	//isExact := positFromFloat(t, p, f)
	if exact != isExact {
		ex := " not"
		isEx := " not"
		if exact {
			ex = ""
		}
		if isExact {
			isEx = ""
		}
		t.Errorf("%v conversion should%v be exact and it is%v exact", name, ex, isEx)
	}
	assertSame(t, pint, expect.Bits, name)
	if exact {
		f1 := p2f(t, expect, pint)
		assertSameF(t, f1, f, "conversion (inverse)")
	}
}

func limitTest(t *testing.T, f *big.Float, p, expect *goposit.SlowPosit, name string) {
	// We need a highly precise float for the "almost" and "just past" tests
	f2 := big.NewFloat(0).Set(f)
	f2.SetPrec(512)

	// Almost
	goposit.Shf(f2, -int(expect.Nbits()))
	f2 = f2.Sub(f, f2)
	conversionTest(t, f2, expect, fmt.Sprintf("Almost %v", name), false)

	// Exact
	conversionTest(t, f, expect, fmt.Sprintf("Exact %v", name), true)

	// Past
	f2.Set(f)
	goposit.Shf(f2, -int(expect.Nbits()))
	f2 = f2.Add(f, f2)
	conversionTest(t, f2, expect, fmt.Sprintf("Just past %v", name), false)

	// Way past
	f2.Set(f)
	if f.MantExp(nil) < 1 {
		goposit.Shf(f2, -int(expect.Nbits()))
	} else {
		goposit.Shf(f2, int(expect.Nbits()))
	}
	conversionTest(t, f2, expect, fmt.Sprintf("Way past %v", name), false)
}

func simpleTestCase(t *testing.T, p *goposit.SlowPosit) {
	f := big.NewFloat(0)
	p2 := p.Clone()

	conversionTest(t, f.SetInf(true), p2.NaR(), "NaR", false)
	conversionTest(t, f.SetInt64(0), p2.Zero(), "zero", true)
	conversionTest(t, f.SetInt64(1), p2.One(), "one", true)
	conversionTest(t, f.SetInt64(-1), p2.NegOne(), "-one", true)

	f.SetInt64(1)
	goposit.Shf(f, p.Log2MaxVal())
	limitTest(t, f, p, p2.Max(), "maxval")

	f.SetInt64(-1)
	goposit.Shf(f, p.Log2MaxVal())
	limitTest(t, f, p, p2.MaxNeg(), "-maxval")

	f.SetInt64(1)
	goposit.Shf(f, -p.Log2MaxVal())
	limitTest(t, f, p, p2.Min(), "minval")

	f.SetInt64(-1)
	goposit.Shf(f, -p.Log2MaxVal())
	limitTest(t, f, p, p2.MinNeg(), "-minval")
}
func TestSimple(t *testing.T) {
	simpleTestCase(t, goposit.NewSlowPosit(8, 0))
	simpleTestCase(t, goposit.NewSlowPosit(16, 1))
	simpleTestCase(t, goposit.NewSlowPosit(32, 2))
	simpleTestCase(t, goposit.NewSlowPosit(64, 3))
	simpleTestCase(t, goposit.NewSlowPosit(128, 4))
}

const bisectDepth = 40

func bisect(t *testing.T, p *goposit.SlowPosit, lesser, greater *big.Float, lesserP, greaterP *big.Int) {
	diff := big.NewFloat(0).Sub(greater, lesser)
	ex := diff.MantExp(nil)
	oneToThe := big.NewFloat(1)
	oneToThe = oneToThe.SetMantExp(oneToThe, ex+1)

	//fmt.Printf("%v %v\n", diff, oneToThe)

	middle := big.NewFloat(0).Copy(lesser)
	middle.Add(middle, oneToThe)
	if middle.Cmp(greater) <= 0 {
		panic("bad test, wrong exponent")
	}
	middle.Sub(middle, oneToThe)
	if middle.Acc() != big.Exact {
		t.Errorf("Not enough precision")
		return
	}

	smallestFind := big.NewFloat(0).SetInf(true)
	for i := 0; i < bisectDepth; i++ {
		middle.Add(middle, oneToThe)
		i2, _ := f2p(t, p, middle)
		if i2.Cmp(lesserP) > 0 {
			if i2.Cmp(greaterP) > 0 {
				// too high
			} else {
				smallestFind.Copy(middle)
			}
			middle.Sub(middle, oneToThe)
		}
		goposit.Shf(oneToThe, -1)
	}
	//fmt.Printf("Smallest: %v\n", smallestFind)
}

func roundingTest(t *testing.T, p *goposit.SlowPosit, hexmin, hexmax string) {
	i0, s := big.NewInt(0).SetString(hexmin, 16)
	assert(s)
	rangeEnd, s := big.NewInt(0).SetString(hexmax, 16)
	assert(s)

	big1 := big.NewInt(1)
	big2 := big.NewInt(2)
	pbigger := goposit.NewSlowPosit(p.Nbits()+1, p.Es())

	f0 := p2f(t, p, i0)
	for i0.Cmp(rangeEnd) != 0 {
		i1, exact := f2p(t, p, f0)
		if i1.Cmp(i0) != 0 || !exact {
			t.Errorf("Inexact encoding of %v to %v original was %v", f0.Text('p', 16), i1.Text(16), i0.Text(16))
			// assertSame(i1, i0)
			// assert(!mpz_cmp(fc->i1, fc->i0));
			// assert(exact);
		}

		i1.Add(i1, big1)
		f1 := p2f(t, p, i1)
		i0_5x := big.NewInt(0)
		i0_5x.Add(i0_5x, i0)
		i0_5x.Mul(i0_5x, big2)
		i0_5x.Add(i0_5x, big1)

		f0_5 := p2f(t, pbigger, i0_5x)
		fsmall := big.NewFloat(0).Copy(f0_5)
		goposit.Shf(fsmall, -30)

		f0_51 := big.NewFloat(0).Add(f0_5, fsmall)
		f0_49 := big.NewFloat(0).Sub(f0_5, fsmall)
		if f0_51.Acc() != big.Exact || f0_49.Acc() != big.Exact {
			t.Errorf("not enough precision")
			return
		}

		i0_5 := f2pExact(t, p, f0_5, false)
		i0_51 := f2pExact(t, p, f0_51, false)
		i0_49 := f2pExact(t, p, f0_49, false)
		if f0.Sign() > 0 {
			assertSame(t, i0_49, i0, "rounding test (positive) assert i0_49 = i0")
			assertSame(t, i0_51, i1, "rounding test (positive) assert i0_51 = i1")
		} else {
			assertSame(t, i0_49, i1, "rounding test (negative) assert i0_49 = i1")
			assertSame(t, i0_51, i0, "rounding test (negative) assert i0_51 = i0")
		}
		if i0.Bit(0) == 1 {
			assertSame(t, i0_5, i1, "dead middle round-to-even i0_5 = i1")
		} else {
			assertSame(t, i0_5, i0, "dead middle round-to-even i0_5 = i0")
		}

		// swappity
		_i1 := i0
		i0 = i1
		i1 = _i1
		_f1 := f0
		f0 = f1
		f1 = _f1
	}
}

func TestRounding(t *testing.T) {
	// exhaustive testing of posit8
	roundingTest(t, goposit.NewSlowPosit(8, 0), "01", "7f")
	roundingTest(t, goposit.NewSlowPosit(8, 0), "81", "ff")

	// exhaustive testing of posit16
	roundingTest(t, goposit.NewSlowPosit(16, 1), "0001", "7fff")
	roundingTest(t, goposit.NewSlowPosit(16, 1), "8001", "ffff")

	// test only edge cases in posit32
	roundingTest(t, goposit.NewSlowPosit(32, 2), "00000001", "0000000f")
	roundingTest(t, goposit.NewSlowPosit(32, 2), "7ffffff0", "7fffffff")
	roundingTest(t, goposit.NewSlowPosit(32, 2), "80000001", "8000000f")
	roundingTest(t, goposit.NewSlowPosit(32, 2), "fffffff0", "ffffffff")

	// test only edge cases in posit64
	roundingTest(t, goposit.NewSlowPosit(64, 3), "0000000000000001", "000000000000000f")
	roundingTest(t, goposit.NewSlowPosit(64, 3), "7ffffffffffffff0", "7fffffffffffffff")
	roundingTest(t, goposit.NewSlowPosit(64, 3), "8000000000000001", "800000000000000f")
	roundingTest(t, goposit.NewSlowPosit(64, 3), "fffffffffffffff0", "ffffffffffffffff")

	// test only edge cases in posit128
	roundingTest(t, goposit.NewSlowPosit(128, 4), "00000000000000000000000000000001", "0000000000000000000000000000000f")
	roundingTest(t, goposit.NewSlowPosit(128, 4), "7ffffffffffffffffffffffffffffff0", "7ffffffffffffffffffffffffffffff1")
	roundingTest(t, goposit.NewSlowPosit(128, 4), "80000000000000000000000000000001", "8000000000000000000000000000000f")
	roundingTest(t, goposit.NewSlowPosit(128, 4), "fffffffffffffffffffffffffffffff0", "ffffffffffffffffffffffffffffffff")
}
