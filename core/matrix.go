package core

import mgl32 "github.com/go-gl/mathgl/mgl32"

// Mat4 is a wrapper around mgl32.Mat4
type Mat4 mgl32.Mat4

// CreateMVP creates Model-View-Projection matrix for shader.
// Accepts object position, screen size and camera parameters.
func CreateMVP(objPos Vec3, width, height float32, camPos, camTarget Vec3) mgl32.Mat4 {
	// Model Matrix 
	model := mgl32.Translate3D(objPos.X, objPos.Y, objPos.Z)

	// Camera Matrix
	view := mgl32.LookAtV(
		mgl32.Vec3{camPos.X, camPos.Y, camPos.Z},
		mgl32.Vec3{camTarget.X, camTarget.Y, camTarget.Z},
		mgl32.Vec3{0, 1, 0},
	)

	// Perspective matrix
	aspect := width / height
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), aspect, 0.1, 100.0)
	return projection.Mul4(view).Mul4(model)
}
