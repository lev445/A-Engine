package core

import (
	"fmt"
	"image"
	_ "image/draw"
	"os"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	_ "golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Font struct {
	ttfFont  *truetype.Font
	atlas    *image.RGBA
	texture  uint32
	glyphs   map[rune]*Glyph
	lineHeight float32
}
type Glyph struct {
	Width, Height float32
	BearingX, BearingY float32
	Advance float32
	U, V, U2, V2 float32
}

func LoadFont(path string, fontSize float32) (*Font, error) {
	fontBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("\033[41mERR\033[0m :: app:ReadFile: Failed to read file: %v", err)
	}
	ttfFont, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, fmt.Errorf("\033[41mERR\033[0m :: app:ParseFont: Failed to parse font: %v", err)
	}

	f := &Font{
		ttfFont:  ttfFont,
		glyphs:   make(map[rune]*Glyph),
		lineHeight: fontSize,
	}
	atlasSize := 512
	f.atlas = image.NewRGBA(image.Rect(0, 0, atlasSize, atlasSize))
	fc := freetype.NewContext()
	fc.SetDPI(72)
	fc.SetFont(ttfFont)
	fc.SetFontSize(float64(fontSize))
	fc.SetSrc(image.White)
	fc.SetDst(f.atlas)
	var x, y float32 = 0, 0
	var maxRowHeight float32 = 0.0

	for ch := rune(32); ch <= 126; ch++ {
		glyphIndex := ttfFont.Index(ch)
		if glyphIndex == 0 {
			continue
		}
		pt := freetype.Pt(int(x), int(y)+int(fontSize))
		_, err := fc.DrawString(string(ch), pt)
		if err != nil {
			continue
			return nil, fmt.Errorf("\033[41mWRN\033[0m :: app:DrawString: Perhaps some symbol could not be drawn properly. This is not a error, but this event is catched: %v", err)
		}
		advance := ttfFont.HMetric(fixed.I(int(fontSize)), glyphIndex).AdvanceWidth
		glyphWidth := float32(advance) / 64.0
		bounds := ttfFont.Bounds(fixed.I(int(fontSize)))
		glyphHeight := float32(bounds.Max.Y-bounds.Min.Y) / 64.0
		//advance := ttfFont.HMetric(fixed.I(int(fontSize)), glyphIndex).AdvanceWidth
		//bounds := ttfFont.Bounds(fixed.I(int(fontSize)))
		glyphWidth = float32(advance) / float32(1<<16)
		glyphHeight = float32(bounds.Max.Y-bounds.Min.Y) / float32(1<<16)
		f.glyphs[ch] = &Glyph{
			Width:  glyphWidth,
			Height: glyphHeight,
			BearingX: 0,
			BearingY: 0,
			Advance: glyphWidth,
			U:  x / float32(atlasSize),
			V:  y / float32(atlasSize),
			U2: (x + glyphWidth) / float32(atlasSize),
			V2: (y + glyphHeight) / float32(atlasSize),
		}
		x += glyphWidth + 2
		if x+glyphWidth > float32(atlasSize) {
			x = 0
			y += float32(maxRowHeight) + 2
			maxRowHeight = 0
		}

		if glyphHeight > maxRowHeight {
			maxRowHeight = glyphHeight
		}
	}
	gl.GenTextures(1, &f.texture)
	gl.BindTexture(gl.TEXTURE_2D, f.texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(
		gl.TEXTURE_2D, 0, gl.RGBA,
		int32(atlasSize), int32(atlasSize), 0,
		gl.RGBA, gl.UNSIGNED_BYTE,
		gl.Ptr(f.atlas.Pix),
	)

	return f, nil
}
func (f *Font) GetGlyph(ch rune) *Glyph {
	if g, ok := f.glyphs[ch]; ok {
		return g
	}
	return f.glyphs['?']
}
func (f *Font) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, f.texture)
}
func (f *Font) Cleanup() {
	gl.DeleteTextures(1, &f.texture)
}
