# GoPosit

Status: incomplete/unmaintained - msg me cjd on mastodon.social if you want to take over

Posit Math library in golang

**Caution**: This is not very well tested yet, use at your own risk

This library is currently backed by the golang big.Float implementation in order to avoid risk
of error. In the future, faster implementations will be written but the API will stay the same.

## API

This library provides support for the standard Posit8, Posit16, Posit32 and Posit64.
Each type has 21 supported functions, except for Posit8 which are missing the downcast function
because there is no smaller size.

Supported functions:

* `p.Add(x Posit<T>) (z Posit<T>)` Add two posits of the same type, output another one, round to
nearest even.
* `p.AddExact(x Posit) (z Posit<T>, r Posit<T>)` Add two posits of the same type, output 2 posits,
one being a truncated sum (with posit encoding, truncation means round-to-negative-infinity) and
another which contains the exact difference between the result and the actual sum. If the two
resulting posits can be added together, they will yield the exact sum.
* `p.Sub(x Posit<T>) (z Posit<T>)` Same as Add but x is subtracted from p.
* `p.SubExact(x Posit<T>) (z Posit<T>, r Posit<T>)` Same as SubExact except x is subtracted from p.
* `p.Mul(x Posit<T>) (z Posit<T>)` Multiply two posits, round to nearest even.
* `p.MulPromote(x Posit<T>) (x Posit<T+1>)` Multiply two posits, returns a posit of the next larger
size. No rounding is needed.
* `p.Div(x Posit<T>) (z Posit<T>)` Divides two posits, round to nearest even.
* `p.DivPromote(x Posit<T>) (z Posit<T+1>)` Divides two posits, returns a posit of the next smaller
size. No rounding is needed.
* `p.Sqrt() (z Posit<T>)` Take the square root of a posit, round to nearest even.
* `p.FromInt(i int<nbits>) (z Posit<T>)` Convert a signed integer of size n to a posit of size n,
for example you can convert an 8 bit integer to a Posit8 and so on.
* `p.FromUint(ui uint<nbits>) (z Posit<T>)` Convert an unsigned integer of size n to a posit of the
same size.
* `p.Int() (i int<nbits>)` Get the posit value as an integer, if the posit is larger than maximum
or smaller than the negative maximum for the integer size, it saturates. The number is rounded to
nearest even.
* `p.Uint() (ui int<nbits>)` Get the posit value as an unsigned integer, if the posit is larger
than the maxumim for the integer size, it saturates. The number is rounded to nearest even.
* `p.Exp() (e int<nbits>)` Get the exponent for the posit as a signed integer the same size as the
posit.
* `p.Mant() (z Posit<T>)` Get a new posit which has the same mantissa but whose exponent is set
to 0.
* `p.ExpAdd(i int<nbits>) (z Posit<T>)` Get a new posit which has i added to it's exponent, if the
regime grows forcing the posit to lose bits, it will round to nearest even. The input is signed
so it can both increase or reduce the exponent.
* `p.Up() (z Posit<T+1>)` Convert the posit to the next larger size.
* `p.Down() (z Posit<T-1>)` Convert the posit to the next smaller size.
* `p.Bits() (ui uint<nbits>)` Get the binary representation of the posit as an unsigned integer
of the same size as the posit.
* `p.SetBits(ui uint<nbits>) Posit<T>` Create a new posit with the specified bits.
* `p.Clone() (z Posit<T>)` Make a copy of the posit.


### Limited Posit128

This library is missing a lot of "soft" posit functionality because it is intended to model what
is expected to be provided by a hardware posit implementation, Posit64.MulPromote() and
Posit64.DivPromote() return a Posit128 which is bigger than the word size on most computers.

In order to make it possible to do Posit64 multiplication without losing bits (and to implement)
a quire in software, there is minimal support for a Posit128, just enough to be able to efficiently
break it back down into smaller numbers and work with them.

* `Posit128.Exp() int32` Get the exponent for the posit as a signed 32 bit integer.
* `Posit128.Mant() Posit128` Get a new posit which has the same mantissa but whose exponent is set
to 0.

With these two functions, you can extract the exponent and then "normalize" the Posit128 such that
the header is a known bit width and can be efficiently masked off, then what you have left is a
bunch of bits that can be added to a quire or multi-precision integer.


### Vector operations

Two Posit vectors are currently supported: `Posit8x4` and `Posit16x2`. Posit vectors can be created
using `NewPosit8x4(posit8)` which will create a vector of 4 copies of the same `Posit8`. All of the
20 functions work the same way on vectors, functions which take a signed or unsigned integer take
a slice of the same integer type instead.

2 additional functions which are available for vectors are:

* `v.Get(i int) Posit<T>` returns one posit from the vector
* `v.Put(i int, p Posit<T>)` sets one posit in the vector


## SlowPosit

There is a big.Float backed posit called SlowPosit, you can use this for emulating posits of any
size that you want, it provides all of the same functions as the normal posit sizes and is the
backing for any which are not otherwise implemented.

* `NewSlowPosit(nbits uint, exponentSize uint) *SlowPosit` Create a new SlowPosit
