/*
 *
 *
 *
*/
package core

import "math"

/**** 2D VECTORS (X,Y). It's using for text, image, texture. ****/
type Vec2 struct {
	X,Y float32
}

// The `Add` return sum of vectors.
func (v Vec2) Add(other Vec2) Vec2 {
	return Vec2{v.X + other.X, v.Y + other.Y}
}

// The `Sub` return sub of vectors.
func (v Vec2) Sub(other Vec2) Vec2 {
	return Vec2{v.X - other.X, v.Y - other.Y}
}

// The `Dot` return dot of vectors.
func (v Vec2) Dot(other Vec2) float32 { 
	return v.X * other.X + v.Y * other.Y
}

// The `Length2` return the squared length of vector...
func (v Vec2) Length2() float32 {
	return v.X*v.X + v.Y*v.Y
}

// The `Length` return the Euclidean lenght of vector.
func (v Vec2) Length() float32 {
	return float32(math.Sqrt(float64(v.Length2())))
}

// The `Normalize` return unit-length vector.
func (v Vec2) Normalize() Vec2 {
	l := v.Length()
	if l == 0 {
		return Vec2{}
	}
	inv := 1.0 / l
	return Vec2{v.X * inv, v.Y * inv}
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

// Dot product.
func (v Vec3) Dot(other Vec3) float32 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// Length returns vector length.
func (v Vec3) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

// Normalize returns normalized vector.
func (v Vec3) Normalize() Vec3 {
	l := v.Length()
	if l == 0 {
		return Vec3{}
	}
	return Vec3{v.X / l, v.Y / l, v.Z / l}
}

// Cross Product.
func (v Vec3) Cross(other Vec3) Vec3 {
	return Vec3{
		v.Y*other.Z - v.Z*other.Y,
		v.Z*other.X - v.X*other.Z,
		v.X*other.Y - v.Y*other.X,
	}
}
