
__COMMENT__ VEC_T is a WIDTH bits wide vector of ONE_T
type VEC_T struct{ impl []ONE_T }

__COMMENT__ GLUE(New, VEC_T) makes a new vector of WIDTH ONE_T
func GLUE(New, VEC_T)(a ONE_T) VEC_T {
    out := VEC_T{impl: make([]ONE_T, WIDTH)}
    for i := 0; i < WIDTH; i++ {
        out.impl[i] = a.Clone()
    }
    return out
}

#define VECTOR_OP(_OP_) \
    __COMMENT__ _OP_ provides a thin wrapper around ONE_T._OP_ __NEWLINE__\
    func (v VEC_T) _OP_(x VEC_T) VEC_T { \
        out := VEC_T{impl: make([]ONE_T, WIDTH)}; \
        for i, posit := range v.impl { out.impl[i] = posit._OP_(x.impl[i]) }; \
        return out; \
    }

#define NOARG_WORD_OP(_OP_, _WORD_) \
    __COMMENT__ _OP_ provides a thin wrapper around ONE_T._OP_ __NEWLINE__\
    func (v VEC_T) _OP_() []_WORD_ { \
        out := make([]_WORD_, WIDTH); \
        for i, posit := range v.impl { out[i] = posit._OP_() }; \
        return out; \
    }

#define WORD_OP(_OP_, _WORD_) \
    __COMMENT__ _OP_ provides a thin wrapper around ONE_T._OP_ __NEWLINE__\
    __COMMENT__ if x is not WIDTH long, this function will panic __NEWLINE__\
    func (v VEC_T) _OP_(x []_WORD_) VEC_T { \
        if len(x) != WIDTH { panic("unexpected length of input"); }; \
        out := VEC_T{impl: make([]ONE_T, WIDTH)}; \
        for i, posit := range v.impl { out.impl[i] = posit._OP_(x[i]) }; \
        return out; \
    }

#define NOARG_OP(_OP_) \
    __COMMENT__ _OP_ provides a thin wrapper around ONE_T._OP_ __NEWLINE__\
    func (v VEC_T) _OP_() VEC_T { \
        out := VEC_T{impl: make([]ONE_T, WIDTH)}; \
        for i, posit := range v.impl { out.impl[i] = posit._OP_() }; \
        return out; \
    }


VECTOR_OP(Add)

__COMMENT__ AddExact is a thin wrapper around ONE_T.AddExact
__COMMENT__ One vector is returned which is a vector of sums and the second vector
__COMMENT__ is a vector of remainders.
func (v VEC_T) AddExact(x VEC_T) (VEC_T, VEC_T) {
	out := VEC_T{impl: make([]ONE_T, len(v.impl))}
	diff := VEC_T{impl: make([]ONE_T, len(v.impl))}
	for i := 0; i < len(v.impl); i++ {
		out.impl[i], diff.impl[i] = v.impl[i].AddExact(x.impl[i])
	}
	return out, diff
}

VECTOR_OP(Sub)

__COMMENT__ SubExact is a thin wrapper around ONE_T.SubExact
__COMMENT__ One vector is returned which is a vector of sums and the second vector
__COMMENT__ is a vector of remainders.
func (v VEC_T) SubExact(x VEC_T) (VEC_T, VEC_T) {
	out := VEC_T{impl: make([]ONE_T, len(v.impl))}
	diff := VEC_T{impl: make([]ONE_T, len(v.impl))}
	for i := 0; i < len(v.impl); i++ {
		out.impl[i], diff.impl[i] = v.impl[i].SubExact(x.impl[i])
	}
	return out, diff
}

VECTOR_OP(Mul)

VECTOR_OP(Div)

WORD_OP(FromInt, SWORD)

WORD_OP(FromUint, UWORD)

NOARG_WORD_OP(Int, SWORD)

NOARG_WORD_OP(Uint, UWORD)

NOARG_WORD_OP(Exp, SWORD)

NOARG_OP(Sqrt)

WORD_OP(ExpAdd, SWORD)

NOARG_WORD_OP(Bits, UWORD)

WORD_OP(SetBits, UWORD)

NOARG_OP(Clone)

__COMMENT__ Get provides access to one of the posits in the vector
func (v* VEC_T) Get(i int) ONE_T { return v.impl[i] }

__COMMENT__ Put updates one of the posits in the vector
func (v* VEC_T) Put(i int, x ONE_T) { v.impl[i] = x }

#undef VECTOR_OP
#undef NOARG_WORD_OP
#undef WORD_OP
#undef NOARG_OP
