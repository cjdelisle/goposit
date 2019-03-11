package goposit

// Posit8x4 is a 4 bits wide vector of Posit8
type Posit8x4 struct{ impl []Posit8 }

// NewPosit8x4 makes a new vector of 4 Posit8
func NewPosit8x4(a Posit8) Posit8x4 {
	out := Posit8x4{impl: make([]Posit8, 4)}
	for i := 0; i < 4; i++ {
		out.impl[i] = a.Clone()
	}
	return out
}

// Add provides a thin wrapper around Posit8.Add
func (v Posit8x4) Add(x Posit8x4) Posit8x4 {
	out := Posit8x4{impl: make([]Posit8, 4)}
	for i, posit := range v.impl {
		out.impl[i] = posit.Add(x.impl[i])
	}
	return out
}

// AddExact is a thin wrapper around Posit8.AddExact
// One vector is returned which is a vector of sums and the second vector
// is a vector of remainders.
func (v Posit8x4) AddExact(x Posit8x4) (Posit8x4, Posit8x4) {
	out := Posit8x4{impl: make([]Posit8, len(v.impl))}
	diff := Posit8x4{impl: make([]Posit8, len(v.impl))}
	for i := 0; i < len(v.impl); i++ {
		out.impl[i], diff.impl[i] = v.impl[i].AddExact(x.impl[i])
	}
	return out, diff
}

// Sub provides a thin wrapper around Posit8.Sub
func (v Posit8x4) Sub(x Posit8x4) Posit8x4 {
	out := Posit8x4{impl: make([]Posit8, 4)}
	for i, posit := range v.impl {
		out.impl[i] = posit.Sub(x.impl[i])
	}
	return out
}

// SubExact is a thin wrapper around Posit8.SubExact
// One vector is returned which is a vector of sums and the second vector
// is a vector of remainders.
func (v Posit8x4) SubExact(x Posit8x4) (Posit8x4, Posit8x4) {
	out := Posit8x4{impl: make([]Posit8, len(v.impl))}
	diff := Posit8x4{impl: make([]Posit8, len(v.impl))}
	for i := 0; i < len(v.impl); i++ {
		out.impl[i], diff.impl[i] = v.impl[i].SubExact(x.impl[i])
	}
	return out, diff
}

// Mul provides a thin wrapper around Posit8.Mul
func (v Posit8x4) Mul(x Posit8x4) Posit8x4 {
	out := Posit8x4{impl: make([]Posit8, 4)}
	for i, posit := range v.impl {
		out.impl[i] = posit.Mul(x.impl[i])
	}
	return out
}

// Div provides a thin wrapper around Posit8.Div
func (v Posit8x4) Div(x Posit8x4) Posit8x4 {
	out := Posit8x4{impl: make([]Posit8, 4)}
	for i, posit := range v.impl {
		out.impl[i] = posit.Div(x.impl[i])
	}
	return out
}

// FromInt provides a thin wrapper around Posit8.FromInt
// if x is not 4 long, this function will panic
func (v Posit8x4) FromInt(x []int8) Posit8x4 {
	if len(x) != 4 {
		panic("unexpected length of input")
	}
	out := Posit8x4{impl: make([]Posit8, 4)}
	for i, posit := range v.impl {
		out.impl[i] = posit.FromInt(x[i])
	}
	return out
}

// FromUint provides a thin wrapper around Posit8.FromUint
// if x is not 4 long, this function will panic
func (v Posit8x4) FromUint(x []uint8) Posit8x4 {
	if len(x) != 4 {
		panic("unexpected length of input")
	}
	out := Posit8x4{impl: make([]Posit8, 4)}
	for i, posit := range v.impl {
		out.impl[i] = posit.FromUint(x[i])
	}
	return out
}

// Int provides a thin wrapper around Posit8.Int
func (v Posit8x4) Int() []int8 {
	out := make([]int8, 4)
	for i, posit := range v.impl {
		out[i] = posit.Int()
	}
	return out
}

// Uint provides a thin wrapper around Posit8.Uint
func (v Posit8x4) Uint() []uint8 {
	out := make([]uint8, 4)
	for i, posit := range v.impl {
		out[i] = posit.Uint()
	}
	return out
}

// Exp provides a thin wrapper around Posit8.Exp
func (v Posit8x4) Exp() []int8 {
	out := make([]int8, 4)
	for i, posit := range v.impl {
		out[i] = posit.Exp()
	}
	return out
}

// Sqrt provides a thin wrapper around Posit8.Sqrt
func (v Posit8x4) Sqrt() Posit8x4 {
	out := Posit8x4{impl: make([]Posit8, 4)}
	for i, posit := range v.impl {
		out.impl[i] = posit.Sqrt()
	}
	return out
}

// ExpAdd provides a thin wrapper around Posit8.ExpAdd
// if x is not 4 long, this function will panic
func (v Posit8x4) ExpAdd(x []int8) Posit8x4 {
	if len(x) != 4 {
		panic("unexpected length of input")
	}
	out := Posit8x4{impl: make([]Posit8, 4)}
	for i, posit := range v.impl {
		out.impl[i] = posit.ExpAdd(x[i])
	}
	return out
}

// Bits provides a thin wrapper around Posit8.Bits
func (v Posit8x4) Bits() []uint8 {
	out := make([]uint8, 4)
	for i, posit := range v.impl {
		out[i] = posit.Bits()
	}
	return out
}

// Clone provides a thin wrapper around Posit8.Clone
func (v Posit8x4) Clone() Posit8x4 {
	out := Posit8x4{impl: make([]Posit8, 4)}
	for i, posit := range v.impl {
		out.impl[i] = posit.Clone()
	}
	return out
}

// Get provides access to one of the posits in the vector
func (v *Posit8x4) Get(i int) Posit8 { return v.impl[i] }

// Put updates one of the posits in the vector
func (v *Posit8x4) Put(i int, x Posit8) { v.impl[i] = x }

// Posit16x2 is a 2 bits wide vector of Posit16
type Posit16x2 struct{ impl []Posit16 }

// NewPosit16x2 makes a new vector of 2 Posit16
func NewPosit16x2(a Posit16) Posit16x2 {
	out := Posit16x2{impl: make([]Posit16, 2)}
	for i := 0; i < 2; i++ {
		out.impl[i] = a.Clone()
	}
	return out
}

// Add provides a thin wrapper around Posit16.Add
func (v Posit16x2) Add(x Posit16x2) Posit16x2 {
	out := Posit16x2{impl: make([]Posit16, 2)}
	for i, posit := range v.impl {
		out.impl[i] = posit.Add(x.impl[i])
	}
	return out
}

// AddExact is a thin wrapper around Posit16.AddExact
// One vector is returned which is a vector of sums and the second vector
// is a vector of remainders.
func (v Posit16x2) AddExact(x Posit16x2) (Posit16x2, Posit16x2) {
	out := Posit16x2{impl: make([]Posit16, len(v.impl))}
	diff := Posit16x2{impl: make([]Posit16, len(v.impl))}
	for i := 0; i < len(v.impl); i++ {
		out.impl[i], diff.impl[i] = v.impl[i].AddExact(x.impl[i])
	}
	return out, diff
}

// Sub provides a thin wrapper around Posit16.Sub
func (v Posit16x2) Sub(x Posit16x2) Posit16x2 {
	out := Posit16x2{impl: make([]Posit16, 2)}
	for i, posit := range v.impl {
		out.impl[i] = posit.Sub(x.impl[i])
	}
	return out
}

// SubExact is a thin wrapper around Posit16.SubExact
// One vector is returned which is a vector of sums and the second vector
// is a vector of remainders.
func (v Posit16x2) SubExact(x Posit16x2) (Posit16x2, Posit16x2) {
	out := Posit16x2{impl: make([]Posit16, len(v.impl))}
	diff := Posit16x2{impl: make([]Posit16, len(v.impl))}
	for i := 0; i < len(v.impl); i++ {
		out.impl[i], diff.impl[i] = v.impl[i].SubExact(x.impl[i])
	}
	return out, diff
}

// Mul provides a thin wrapper around Posit16.Mul
func (v Posit16x2) Mul(x Posit16x2) Posit16x2 {
	out := Posit16x2{impl: make([]Posit16, 2)}
	for i, posit := range v.impl {
		out.impl[i] = posit.Mul(x.impl[i])
	}
	return out
}

// Div provides a thin wrapper around Posit16.Div
func (v Posit16x2) Div(x Posit16x2) Posit16x2 {
	out := Posit16x2{impl: make([]Posit16, 2)}
	for i, posit := range v.impl {
		out.impl[i] = posit.Div(x.impl[i])
	}
	return out
}

// FromInt provides a thin wrapper around Posit16.FromInt
// if x is not 2 long, this function will panic
func (v Posit16x2) FromInt(x []int16) Posit16x2 {
	if len(x) != 2 {
		panic("unexpected length of input")
	}
	out := Posit16x2{impl: make([]Posit16, 2)}
	for i, posit := range v.impl {
		out.impl[i] = posit.FromInt(x[i])
	}
	return out
}

// FromUint provides a thin wrapper around Posit16.FromUint
// if x is not 2 long, this function will panic
func (v Posit16x2) FromUint(x []uint16) Posit16x2 {
	if len(x) != 2 {
		panic("unexpected length of input")
	}
	out := Posit16x2{impl: make([]Posit16, 2)}
	for i, posit := range v.impl {
		out.impl[i] = posit.FromUint(x[i])
	}
	return out
}

// Int provides a thin wrapper around Posit16.Int
func (v Posit16x2) Int() []int16 {
	out := make([]int16, 2)
	for i, posit := range v.impl {
		out[i] = posit.Int()
	}
	return out
}

// Uint provides a thin wrapper around Posit16.Uint
func (v Posit16x2) Uint() []uint16 {
	out := make([]uint16, 2)
	for i, posit := range v.impl {
		out[i] = posit.Uint()
	}
	return out
}

// Exp provides a thin wrapper around Posit16.Exp
func (v Posit16x2) Exp() []int16 {
	out := make([]int16, 2)
	for i, posit := range v.impl {
		out[i] = posit.Exp()
	}
	return out
}

// Sqrt provides a thin wrapper around Posit16.Sqrt
func (v Posit16x2) Sqrt() Posit16x2 {
	out := Posit16x2{impl: make([]Posit16, 2)}
	for i, posit := range v.impl {
		out.impl[i] = posit.Sqrt()
	}
	return out
}

// ExpAdd provides a thin wrapper around Posit16.ExpAdd
// if x is not 2 long, this function will panic
func (v Posit16x2) ExpAdd(x []int16) Posit16x2 {
	if len(x) != 2 {
		panic("unexpected length of input")
	}
	out := Posit16x2{impl: make([]Posit16, 2)}
	for i, posit := range v.impl {
		out.impl[i] = posit.ExpAdd(x[i])
	}
	return out
}

// Bits provides a thin wrapper around Posit16.Bits
func (v Posit16x2) Bits() []uint16 {
	out := make([]uint16, 2)
	for i, posit := range v.impl {
		out[i] = posit.Bits()
	}
	return out
}

// Clone provides a thin wrapper around Posit16.Clone
func (v Posit16x2) Clone() Posit16x2 {
	out := Posit16x2{impl: make([]Posit16, 2)}
	for i, posit := range v.impl {
		out.impl[i] = posit.Clone()
	}
	return out
}

// Get provides access to one of the posits in the vector
func (v *Posit16x2) Get(i int) Posit16 { return v.impl[i] }

// Put updates one of the posits in the vector
func (v *Posit16x2) Put(i int, x Posit16) { v.impl[i] = x }
