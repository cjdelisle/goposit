package goposit

// Posit128 is a holding place for a 128 bit result of a Posit64.MulPromote() or
// Posit64.DivPromote, it does not implement all of the features, just enough to
// break it down.
type Posit128 struct{ impl *SlowPosit }

// Exp outputs the exponent of the posit, specifically the number z for which
// 0.5*2**z <= positValue < 1*2**z
func (p Posit128) Exp() int32 {
	return p.impl.Exp()
}

// Mant outputs a posit with the same mantissa but an exponent of 0, that is a
// which is greater than or equal to 0.5 and less than 1
func (p Posit128) Mant() Posit128 { return Posit128{impl: p.impl.Mant()} }
