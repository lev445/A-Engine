package core

import (
	"math"
	"testing"
)

func approx(a, b float32) bool {
	return float32(math.Abs(float64(a-b))) < 1e-6
}

func TestVec2Arithmetics(t *testing.T) {
	a := Vec2{1, 2}
	b := Vec2{3, -4}
	if got := a.Add(b); got != (Vec2{4, -2}) {
		t.Errorf("Add = %v", got)
	}
	if got := a.Sub(b); got != (Vec2{-2, 6}) {
		t.Errorf("Sub = %v", got)
	}
	if got := a.Mul(2); got != (Vec2{2, 4}) {
		t.Errorf("Mul = %v", got)
	}
	if got := a.Dot(b); !approx(got, -5) {
		t.Errorf("Dot = %v, want -5", got)
	}
	if got := a.Dist(b); !approx(got, float32(math.Sqrt(4+36))) {
		t.Errorf("Dist = %v", got)
	}
	if got := a.Lerp(b, 0.5); !approx(got.X, 2) || !approx(got.Y, -1) {
		t.Errorf("Lerp = %v, want {2 -1}", got)
	}
}

func TestVec2LengthSq(t *testing.T) {
	v := Vec2{3, 4}
	if got := v.LengthSq(); !approx(got, 25) {
		t.Errorf("LengthSq = %v, want 25", got)
	}
	if got := v.Length(); !approx(got, 5) {
		t.Errorf("Length = %v, want 5", got)
	}
}

func TestVec3LengthSq(t *testing.T) {
	v := Vec3{3, 4, 0}
	if got := v.LengthSq(); !approx(got, 25) {
		t.Errorf("LengthSq = %v, want 25", got)
	}
	if got := v.Length(); !approx(got, 5) {
		t.Errorf("Length = %v, want 5", got)
	}
	if got := v.LengthSq(); approx(got, float32(math.Sqrt(5))) {
		t.Errorf("LengthSq looks like √⁴ (fourth root): %v", got)
	}
}

func TestVec3Normalize(t *testing.T) {
	cases := []struct {
		in, want Vec3
	}{
		{Vec3{3, 0, 0}, Vec3{1, 0, 0}},
		{Vec3{0, 0, -4}, Vec3{0, 0, -1}},
		{Vec3{0, 0, 0}, Vec3{0, 0, 0}},
		{Vec3{1e-9, 0, 0}, Vec3{0, 0, 0}}, // denormal -> zero, without NaN
	}
	for _, c := range cases {
		if got := c.in.Normalize(); got != c.want {
			t.Errorf("Normalize(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestNormalizeIdempotent(t *testing.T) {
	v := Vec3{1, -2, 2}
	n1 := v.Normalize()
	n2 := n1.Normalize()
	if !approx(n1.X, n2.X) || !approx(n1.Y, n2.Y) || !approx(n1.Z, n2.Z) {
		t.Errorf("normalize(normalize(v)) != normalize(v): %v != %v", n1, n2)
	}
	if l := n1.Length(); !approx(l, 1) {
		t.Errorf("|normalize(v)| = %v, want 1", l)
	}
}

func TestCrossBasis(t *testing.T) {
	x := Vec3{1, 0, 0}
	y := Vec3{0, 1, 0}
	z := Vec3{0, 0, 1}
	if got := x.Cross(y); got != z {
		t.Errorf("X x Y = %v, want Z", got)
	}
	if got := y.Cross(z); got != x {
		t.Errorf("Y x Z = %v, want X", got)
	}
	if got := z.Cross(x); got != y {
		t.Errorf("Z x X = %v, want Y", got)
	}
	if got := x.Cross(x); got != (Vec3{}) {
		t.Errorf("X x X = %v, want 0", got)
	}
}

func TestClampf(t *testing.T) {
	if got := Clampf(5, 0, 1); got != 1 {
		t.Errorf("Clampf(5,0,1) = %v", got)
	}
	if got := Clampf(-5, 0, 1); got != 0 {
		t.Errorf("Clampf(-5,0,1) = %v", got)
	}
	if got := Clampf(0.5, 0, 1); got != 0.5 {
		t.Errorf("Clampf(0.5,0,1) = %v", got)
	}
}
