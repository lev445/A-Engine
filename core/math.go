/* ~ Welcome to math module in A+ Engine!
 * This module is manage of math in game engine.
 * Job:
 * ++ Count 2D vectors
 * ++ Count 3D vectors
 * ++ Scalar mul, normalize, cross, lerp/dist helpers
 */
package core

import "math"

/* epsSq — a "zero-length" threshold for the SQUARED length; it protects Normalize
from division by zero and denormal numbers. */
const epsSq = 1e-12

/**** 2D VECTORS (X,Y). It's using for text, image, texture. ****/
type Vec2 struct {
	X, Y float32
}

// The `Add` return sum of vectors.
func (v Vec2) Add(other Vec2) Vec2 {
	return Vec2{v.X + other.X, v.Y + other.Y}
}

// The `Sub` return sub of vectors.
func (v Vec2) Sub(other Vec2) Vec2 {
	return Vec2{v.X - other.X, v.Y - other.Y}
}

// `Mul` multiplies the vector by a scalar value.
func (v Vec2) Mul(scalar float32) Vec2 {
	return Vec2{v.X * scalar, v.Y * scalar}
}

// The `Dot` return dot of vectors.
func (v Vec2) Dot(other Vec2) float32 {
	return v.X*other.X + v.Y*other.Y
}

// LengthSq returns the SQUARED length of the vector (without sqrt).
func (v Vec2) LengthSq() float32 {
	return v.X*v.X + v.Y*v.Y
}

// The `Length` return the Euclidean lenght of vector.
func (v Vec2) Length() float32 {
	return float32(math.Sqrt(float64(v.LengthSq())))
}

// The `Normalize` return unit-length vector.
func (v Vec2) Normalize() Vec2 {
	l2 := v.LengthSq()
	if l2 < epsSq {
		return Vec2{}
	}
	inv := 1.0 / float32(math.Sqrt(float64(l2)))
	return Vec2{v.X * inv, v.Y * inv}
}

// The `Lerp` return linear interpolation between v and other (t in [0..1]).
func (v Vec2) Lerp(other Vec2, t float32) Vec2 {
	return v.Add(other.Sub(v).Mul(t))
}

// The `Dist` return distance between two points.
func (v Vec2) Dist(other Vec2) float32 {
	return v.Sub(other).Length()
}

/**** 3D VECTORS (X,Y,Z). ****/
// Vec3 description.
type Vec3 struct {
	X, Y, Z float32
}

// Add vectors.
func (v Vec3) Add(other Vec3) Vec3 {
	return Vec3{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

// Sub vectors.
func (v Vec3) Sub(other Vec3) Vec3 {
	return Vec3{v.X - other.X, v.Y - other.Y, v.Z - other.Z}
}

// `Mul` multiplies the vector by a scalar value.
func (v Vec3) Mul(scalar float32) Vec3 {
	return Vec3{v.X * scalar, v.Y * scalar, v.Z * scalar}
}

// Dot product.
func (v Vec3) Dot(other Vec3) float32 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// LengthSq returns the SQUARED length of the vector (without sqrt).
func (v Vec3) LengthSq() float32 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

// Length returns vector length.
func (v Vec3) Length() float32 {
	return float32(math.Sqrt(float64(v.LengthSq())))
}

// Normalize returns normalized vector.
func (v Vec3) Normalize() Vec3 {
	l2 := v.LengthSq()
	if l2 < epsSq {
		return Vec3{}
	}
	inv := 1.0 / float32(math.Sqrt(float64(l2)))
	return Vec3{v.X * inv, v.Y * inv, v.Z * inv}
}

// The `Lerp` return linear interpolation between v and other (t in [0..1]).
func (v Vec3) Lerp(other Vec3, t float32) Vec3 {
	return v.Add(other.Sub(v).Mul(t))
}

// The `Dist` return distance between two points.
func (v Vec3) Dist(other Vec3) float32 {
	return v.Sub(other).Length()
}

// Cross Product.
func (v Vec3) Cross(other Vec3) Vec3 {
	return Vec3{
		v.Y*other.Z - v.Z*other.Y,
		v.Z*other.X - v.X*other.Z,
		v.X*other.Y - v.Y*other.X,
	}
}

// Clampf restricts x to the [lo, hi] range.
func Clampf(x, lo, hi float32) float32 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}
