// Copyright 2014, David Minor. All rights reserved.
// Use of this source code is governed by the MIT
// license which can be found in the LICENSE file.

package uint128

import (
	"fmt"
)

type Uint128 struct {
	H, L uint64
}

func (i Uint128) Equals(o Uint128) bool {
	return i.L == o.L && i.H == o.H
}

func (i Uint128) ShiftLeft(bits uint) Uint128 {
	if bits >= 128 {
		i.H = 0
		i.L = 0
	} else if bits >= 64 {
		i.H = i.L << (bits - 64)
		i.L = 0
	} else {
		i.H <<= bits
		i.H |= i.L >> (64 - bits)
		i.L <<= bits
	}
	return i
}

func (i Uint128) ShiftRight(bits uint) Uint128 {
	if bits >= 128 {
		i.H = 0
		i.L = 0
	} else if bits >= 64 {
		i.L = i.H >> (bits - 64)
		i.H = 0
	} else {
		i.L >>= bits
		i.L |= i.H << (64 - bits)
		i.H >>= bits
	}
	return i
}

func (x Uint128) And(y Uint128) Uint128 {
	x.H &= y.H
	x.L &= y.L
	return x
}

func (x Uint128) Xor(y Uint128) Uint128 {
	x.H ^= y.H
	x.L ^= y.L
	return x
}

func (x Uint128) Or(y Uint128) Uint128 {
	x.H |= y.H
	x.L |= y.L
	return x
}

func (augend Uint128) Add(addend Uint128) Uint128 {
	origlow := augend.L
	augend.L += addend.L
	augend.H += addend.H
	if augend.L < origlow { // wrapping occurred, so carry the 1
		augend.H += 1
	}
	return augend
}

// (Adapted from go's math/big)
// z1<<64 + z0 = x*y
// Adapted from Warren, Hacker's Delight, p. 132.
func mult(x, y uint64) (z1, z0 uint64) {
	z0 = x * y // lower 64 bits are easy
	// break the multiplication into (x1 << 32 + x0)(y1 << 32 + y0)
	// which is x1*y1 << 64 + (x0*y1 + x1*y0) << 32 + x0*y0
	// so now we can do 64 bit multiplication and addition and
	// shift the results into the right place
	x0, x1 := x&0x00000000ffffffff, x>>32
	y0, y1 := y&0x00000000ffffffff, y>>32
	w0 := x0 * y0
	t := x1*y0 + w0>>32
	w1 := t & 0x00000000ffffffff
	w2 := t >> 32
	w1 += x0 * y1
	z1 = x1*y1 + w2 + w1>>32
	return
}

func (multiplicand Uint128) Mult(multiplier Uint128) Uint128 {
	hl := multiplicand.H*multiplier.L + multiplicand.L*multiplier.H
	multiplicand.H, multiplicand.L = mult(multiplicand.L, multiplier.L)
	multiplicand.H += hl
	return multiplicand
}

func (u Uint128) Mod(n uint64) uint64 {
	// uses (ab mod n) == (a mod n)(b mod n) mod n
	// and (a + b) mod n == (a mod n) + (b mod n)

	// u = u.H * (2^64) + u.L
	// u % n == (u.H * (2^64)) mod n + u.L mod n
	// u % n == (u.H mod n)(2^64 mod n) + u.L mod n
	// u % n == (u.H mod n)(2^32 mod n)(2^32 mod n) + u.L mod n
	if n > 2097152 {
		fmt.Println("N should be smaller than 2^21 otherwise results may be undefined")
	}
	highModMul := ((4294967296) % n) * ((4294967296) % n)
	return (highModMul*(u.H%n) + (u.L % n)) % n
}

func (u Uint128) Div(n uint64) Uint128 {
	// u = ho + l, o = 2^64
	// u / 6 = (ho / 6) + (l / 6)
	return Uint128{}
}
