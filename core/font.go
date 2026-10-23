package core

import (
	"fmt"
	"image"
	_ "image/draw"
	"math"
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
	// ВАЖНО: без SetClip контекст имеет пустой clip-rectangle (Rect(0,0,0,0)),
	// и DrawString не нарисует НИ ОДНОГО глифа (атлас останется пустым).
	fc.SetClip(f.atlas.Bounds())
	// Пакуем атлас строго по целочисленной сетке: глиф растеризуется ровно
	// в тех координатах, из которых потом берётся UV-окно. Иначе глиф
	// смещён относительно своего UV-окна на долю текселя и биллинейный
	// фильтр размывает каждый символ.
	bounds := ttfFont.Bounds(fixed.I(int(fontSize)))
	cellH := int(math.Ceil(float64(bounds.Max.Y-bounds.Min.Y) / 64.0))
	x, y := 0, 0

	for ch := rune(32); ch <= 126; ch++ {
		glyphIndex := ttfFont.Index(ch)
		if glyphIndex == 0 {
			continue
		}
		pt := freetype.Pt(x, y+int(fontSize))
		_, err := fc.DrawString(string(ch), pt)
		if err != nil {
			continue
		}
		// HMetric возвращает 26.6 fixed-point: делим на 64 и округляем
		// вверх до целого числа текселей.
		advance := ttfFont.HMetric(fixed.I(int(fontSize)), glyphIndex).AdvanceWidth
		advPx := int(math.Ceil(float64(advance) / 64.0))
		f.glyphs[ch] = &Glyph{
			Width:    float32(advPx),
			Height:   float32(cellH),
			BearingX: 0,
			BearingY: 0,
			Advance:  float32(advPx),
			U:        float32(x) / float32(atlasSize),
			V:        float32(y) / float32(atlasSize),
			U2:       float32(x+advPx) / float32(atlasSize),
			V2:       float32(y+cellH) / float32(atlasSize),
		}
		x += advPx + 2
		if x+advPx > atlasSize {
			x = 0
			y += cellH + 2
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
