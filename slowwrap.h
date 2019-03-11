
__COMMENT__ POSIT_T is an NBITS bit posit with ES exponent bits
type POSIT_T struct{ impl *SlowPosit }

__COMMENT__ GLUE(NewPosit, NBITS) makes a new posit with NBITS bits and ES es bits, the initial value is zero
func GLUE(NewPosit, NBITS)() POSIT_T { return POSIT_T{impl: NewSlowPosit(NBITS, ES)} }


__COMMENT__ Add takes the sum of two posits
__COMMENT__ p.Add(x) creates a new posit z which is p+x, neither p nor x are altered
func (p POSIT_T) Add(x POSIT_T) POSIT_T { return POSIT_T{impl: p.impl.Add(x.impl)} }

__COMMENT__ AddExact returns exactly the sum of two posits, represented as two
__COMMENT__ more posits, the first result is a posit which is the sum truncated
__COMMENT__ (round to negative infinity) and the second posit is the difference between
__COMMENT__ the first posit and the true sum. That is to say true_sum = out2 + out1.
func (p POSIT_T) AddExact(x POSIT_T) (POSIT_T, POSIT_T) {
    res, diff := p.impl.AddExact(x.impl)
	return POSIT_T{impl: res}, POSIT_T{impl: diff}
}

__COMMENT__ Sub takes the difference of two posits
__COMMENT__ p.Sub(x) creates a new posit z which is p-x, neither p nor x are altered
func (p POSIT_T) Sub(x POSIT_T) POSIT_T { return POSIT_T{impl: p.impl.Sub(x.impl)} }


__COMMENT__ SubExact returns exactly the difference of two posits, represented as two
__COMMENT__ more posits, the first result is a posit which is the difference truncated
__COMMENT__ (round to negative infinity) and the second posit is the difference between
__COMMENT__ the first posit and the true sum. That is to say true_sum = out2 + out1.
func (p POSIT_T) SubExact(x POSIT_T) (POSIT_T, POSIT_T) {
    res, diff := p.impl.SubExact(x.impl)
	return POSIT_T{impl: res}, POSIT_T{impl: diff}
}

__COMMENT__ Mul takes the product of two posits
__COMMENT__ p.Mul(x) creates a new posit z which is x*y, neither p nor x are altered
func (p POSIT_T) Mul(x POSIT_T) POSIT_T { return POSIT_T{impl: p.impl.Mul(x.impl)} }

#ifdef BIGGER_T
__COMMENT__ MulPromote takes the product of two posits
__COMMENT__ p.Mul(x) creates a new posit z which is p*x represented as the next larger posit size
__COMMENT__ next larger means the bit width is doubled and the exponent size is increased by 1
__COMMENT__ this function is guaranteed not to round
func (p POSIT_T) MulPromote(x POSIT_T) BIGGER_T { return BIGGER_T{impl: p.impl.MulPromote(x.impl)} }
#endif

__COMMENT__ Div takes the quotent of two posits
__COMMENT__ p.Div(x) creates a new posit z which is p/x, neither p nor x are altered
func (p POSIT_T) Div(x POSIT_T) POSIT_T { return POSIT_T{impl: p.impl.Div(x.impl)} }

#ifdef BIGGER_T
__COMMENT__ DivPromote takes the quotent of two posits
__COMMENT__ p.DivPromote(x) creates a new posit z which is p/x, as the next larger posit size
__COMMENT__ next larger means the bit width is doubled and the exponent size is increased by 1
__COMMENT__ this function is guaranteed not to round
func (p POSIT_T) DivPromote(x POSIT_T) BIGGER_T { return BIGGER_T{impl: p.impl.DivPromote(x.impl)} }
#endif

__COMMENT__ Sqrt finds the square root of a posit
__COMMENT__ p.Sqrt() returns a new posit, p is not altered
func (p POSIT_T) Sqrt() POSIT_T { return POSIT_T{impl: p.impl.Sqrt()} }

//__COMMENT__ conversions ////

__COMMENT__ Size specific

__COMMENT__ FromInt creates a new posit which is set to the value of a integer
__COMMENT__ p.FromInt() outputs a new posit z, p is not altered
func (p POSIT_T) FromInt(i SWORD) POSIT_T { return POSIT_T{impl: p.impl.FromInt(int64(i))} }

__COMMENT__ FromUint creates a new posit which is set to the value of an unsigned integer
__COMMENT__ p.FromUint() outputs a new posit z, p is not altered
func (p POSIT_T) FromUint(i UWORD) POSIT_T { return POSIT_T{impl: p.impl.FromUint(uint64(i))} }

__COMMENT__ Int outputs an SWORD representation of the value of the posit
__COMMENT__ lost bits are rounded to nearest even and if the number is greater than or equal
__COMMENT__ to two to the NBITS-1 power, the maximum SMAX for positive input or
__COMMENT__ UMAX for negative input will return.
__COMMENT__ Likewise if the number is less than or equal to 1/2 it will round down to zero.
func (p POSIT_T) Int() SWORD {
    x := p.impl.Int()
    if x > SMAX {
        return SMAX
    }
    if x < -SMAX {
        return -SMAX
    }
    return SWORD(x)
}

__COMMENT__ Uint outputs a UWORD representation of the value of the posit, if the posit
__COMMENT__ is less than or equal to 1/2, 0 will be returned and if the value is greater
__COMMENT__ than two to the maxBits power, 0xffffffffffffffff will be returned.
func (p POSIT_T) Uint() UWORD {
    x := p.impl.Int()
    #if UMAX < 0xffffffffffffffff
    if x > UMAX {
        return UMAX
    }
    #endif
    return UWORD(x)
}

//__COMMENT__ binary manipulation ////

__COMMENT__ Exp outputs the exponent of the posit, specifically the number z for which
__COMMENT__ 0.5*2**z <= positValue < 1*2**z
func (p POSIT_T) Exp() SWORD {
    x := p.impl.Exp()
    #if UMAX < 0xffffffffffffffff
    if x > SMAX || x < -SMAX {
        panic("posit with bigger exponent than should be possible")
    }
    #endif
    return SWORD(x)
}

__COMMENT__ Mant outputs a posit with the same mantissa but an exponent of 0, that is a
__COMMENT__ number which is greater than or equal to 0.5 and less than 1
func (p POSIT_T) Mant() POSIT_T { return POSIT_T{impl: p.impl.Mant()}; }

__COMMENT__ ExpAdd returns a new posit with x added to the exponent, effectively multiplying the
__COMMENT__ posit value by 2**x. It is equivilent to a bit-shift in integer math. If bits is
__COMMENT__ negative then the exponent will be decreased, if bits causes the exponent to increase
__COMMENT__ in magnitude, it might cause the posit to round.
func (p POSIT_T) ExpAdd(x SWORD) POSIT_T { return POSIT_T{impl: p.impl.ExpAdd(int32(x))} }

#ifdef BIGGER_T
__COMMENT__ Up returns the same number up-casted to the next larger posit, in this case larger means
__COMMENT__ twice the nbits and an exponent size which is one greater.
func (p POSIT_T) Up() BIGGER_T {
    return BIGGER_T{impl:p.impl.Up()}
}
#endif

#ifdef SMALLER_T
__COMMENT__ Down returns the same number down-casted to the next smaller posit, see Up() for the
__COMMENT__ definition of larger and smaller.
func (p POSIT_T) Down() SMALLER_T {
    return SMALLER_T{impl:p.impl.Up()}
}
#endif

__COMMENT__ Bits outputs a UWORD containing the raw binary format of the posit
func (p POSIT_T) Bits() UWORD {
    words := p.impl.Bits.Bits()
    if len(words) == 1 {
        return UWORD(words[0])
    } else if len(words) == 2 {
        if bits.UintSize != 32 {
            panic("unexpected number of bits")
        }
        return (UWORD(words[1]) << (NBITS/2)) | UWORD(words[0])
    }
    panic("unexpected word size")
}

__COMMENT__ SetBits outputs a new posit with the bits set to those which you specify
func (p POSIT_T) SetBits(bits UWORD) POSIT_T {
    bi := new(big.Int).SetUint64(uint64(bits))
    return POSIT_T{impl: &SlowPosit{es:p.impl.es, nbits:p.impl.nbits, Bits:bi}}
}

__COMMENT__ Clone makes a copy of a posit
func (p POSIT_T) Clone() POSIT_T {
    return POSIT_T{impl: p.impl.Clone()}
}
