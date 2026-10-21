package core
import (
	"math"
	mgl32 "github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	Position       Vec3
	Yaw, Pitch     float32
	Front,Up,Right Vec3
	Speed          float32
}

// create a new camera
func NewCam(pos Vec3) *Camera {
	c:=&Camera {
		Position: pos,
		Yaw:     -90.0,
		Pitch:    0.0,
		Speed:    5.0,
		Up:       Vec3{0, 1, 0},
	}
	c.updateVectors()
	return c
}

// count `direction vectors` based on Yaw and Pitch angles.
func (c *Camera) updateVectors() {
	yawRad := float32(math.Cos(float64(mgl32.DegToRad(c.Yaw))))
	pitchRad := float32(math.Cos(float64(mgl32.DegToRad(c.Pitch))))
	c.Front = Vec3{
		X: yawRad * pitchRad,
		Y: float32(math.Sin(float64(mgl32.DegToRad(c.Pitch)))),
		Z: float32(math.Sin(float64(mgl32.DegToRad(c.Yaw)))) * pitchRad,
	}
	c.Front = c.Front.Normalize()
	c.Right = c.Front.Cross(c.Up).Normalize()
}

// return `matrix view` for shaders.
func (c *Camera) GetViewMatrix() mgl32.Mat4 {
	target := Vec3{
		c.Position.X + c.Front.X,
		c.Position.Y + c.Front.Y,
		c.Position.Z + c.Front.Z,
	}
	return mgl32.LookAtV(
		mgl32.Vec3{c.Position.X, c.Position.Y, c.Position.Z},
		mgl32.Vec3{target.X, target.Y, target.Z},
		mgl32.Vec3{c.Up.X, c.Up.Y, c.Up.Z},
	)
}

// handles clicks
func (c *Camera) ProcessKeyboard(direction string, deltaTime float32) {
	velocity := c.Speed * deltaTime
	var moveDir Vec3
	switch direction {
	case "FORWARD":
		moveDir = c.Front
	case "BACKWARD":
		moveDir = Vec3{-c.Front.X, -c.Front.Y, -c.Front.Z}
	case "LEFT":
		moveDir = Vec3{-c.Right.X, -c.Right.Y, -c.Right.Z}
	case "RIGHT":
		moveDir = c.Right
	}
	c.Position = c.Position.Add(Vec3{
		moveDir.X * velocity,
		moveDir.Y * velocity,
		moveDir.Z * velocity,
	})
}

func (c *Camera) ProcessMouseMove(xOffset, yOffset float32, constrainPitch bool) {
	sensitivity := float32(0.1)
	xOffset *= sensitivity
	yOffset *= sensitivity
	c.Yaw += xOffset
	c.Pitch += yOffset
	if constrainPitch {
		if c.Pitch > 89.0 {
			c.Pitch = 89.0
		}
		if c.Pitch < -89.0 {
			c.Pitch = -89.0
		}
	}
	c.updateVectors()
	
}
