package core
import (
	"github.com/go-gl/gl/v4.1-core/gl"
	mgl32 "github.com/go-gl/mathgl/mgl32"
)

type TextRenderer struct {
	program uint32
	vao, vbo uint32
	projLoc int32
	colorLoc int32
	texLoc int32
	orthoProjection mgl32.Mat4
}

func NewTextRenderer(othoProjection mgl32.Mat4) *TextRenderer {
	return &TextRenderer{
		orthoProjection: othoProjection,
	}
}

func (tr *TextRenderer) SetProgram(program uint32) {
	tr.program = program
	tr.projLoc = gl.GetUniformLocation(program, gl.Str("uProjection\x00"))
	tr.colorLoc = gl.GetUniformLocation(program, gl.Str("uTextColor\x00"))
	tr.texLoc = gl.GetUniformLocation(program, gl.Str("uFontTexture\x00"))
}

func (tr *TextRenderer) SetupVAO() {
	gl.GenVertexArrays(1, &tr.vao)
	gl.GenBuffers(1, &tr.vbo)
	gl.BindVertexArray(tr.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, tr.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 6*4*4, nil, gl.DYNAMIC_DRAW)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 4*4, gl.PtrOffset(2*4))
	gl.EnableVertexAttribArray(1)
	gl.BindVertexArray(0)
}

// RenderText -- this function draw text via font
func (tr *TextRenderer) RenderText(font *Font, text string, x, y, scale float32, r, g, b float32) {
	if tr.program == 0 || font == nil {
		return
	}
	gl.UseProgram(tr.program)
	gl.UniformMatrix4fv(tr.projLoc, 1, false, &tr.orthoProjection[0])
	gl.Uniform4f(tr.colorLoc, r, g, b, 1.0)
	gl.Uniform1i(tr.texLoc, 0)
	font.Bind()
	gl.BindVertexArray(tr.vao)
	currentX := x
	for _, ch := range text {
		glyph := font.GetGlyph(ch)
		charX := currentX + glyph.BearingX*scale
		charY := y - glyph.BearingY*scale
		charW := glyph.Width * scale
		charH := glyph.Height * scale
		vertices := []float32 {
			// POSITION            UV
			charX, charY, glyph.U, glyph.V,
			charX + charW, charY, glyph.U2, glyph.V,
			charX + charW, charY + charH, glyph.U2, glyph.V2,
			charX, charY, glyph.U, glyph.V,
			charX + charW, charY + charH, glyph.U2, glyph.V2,
			charX, charY + charH, glyph.U, glyph.V2,
		}
		gl.BindBuffer(gl.ARRAY_BUFFER, tr.vbo)
		gl.BufferSubData(gl.ARRAY_BUFFER,0,len(vertices)*4, gl.Ptr(vertices))
		gl.DrawArrays(gl.TRIANGLES,0,6)
	}
	gl.BindVertexArray(0)
}

func (tr *TextRenderer) Cleanup() {
	gl.DeleteVertexArrays(1, &tr.vao)
	gl.DeleteBuffers(1, &tr.vbo)
}
