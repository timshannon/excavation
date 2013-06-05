package engine

import (
	"bitbucket.org/tshannon/gohorde/horde3d"
	"code.google.com/p/freetype-go/freetype"
	"errors"
	"image"
	"image/draw"
	"strconv"
)

const (
	freeTypeDPI = 72
)

var textIncrementer = 0

//Text is a freetype rasterized text to a horde overlay
// an slice of strings is rasterized using 72 dpi
// Each entry in the slice is a separate line spaced
// according to the lineSpacing value 2 = doublespaced
type Text struct {
	text        []string
	area        *ScreenArea
	size        float64
	lineSpacing float64
	fontFile    string
	overlay     *Overlay
	textureRes  *Texture
	context     *freetype.Context
	background  *image.RGBA
}

func NewText(text []string, fontFile string, size float64,
	color *Color, area *ScreenArea) *Text {

	newText := &Text{
		text:        text,
		area:        area,
		size:        size,
		lineSpacing: 1,
		fontFile:    fontFile,
	}

	textIncrementer++
	name := "TextVirtualOverlay_" + strconv.Itoa(textIncrementer)
	materialData := `<Material>
		<Shader source="shaders/overlay.shader"/>
		
		<Sampler name="albedoMap" map="` + name + `" />
		</Material>`

	newText.overlay = &Overlay{
		Dimensions: area,
		Color:      color,
		Material: &Material{
			NewVirtualResource(name+".material.xml",
				ResTypeMaterial, []byte(materialData))},
	}

	fontData, err := loadEngineData(fontFile)
	if err != nil {
		RaiseError(err)
		return nil
	}

	font, err := freetype.ParseFont(fontData)
	if err != nil {
		RaiseError(err)
		return nil
	}

	newText.context = freetype.NewContext()
	newText.context.SetDPI(freeTypeDPI)
	newText.context.SetFont(font)
	newText.context.SetFontSize(size)

	newText.textureRes = NewVirtualTexture(name, newText.area.PixelWidth(),
		newText.area.PixelHeight(), horde3d.Formats_TEX_BGRA8, horde3d.ResFlags_NoTexMipmaps)

	newText.overlay.Material.Load()
	newText.overlay.Material.SetResParamI(horde3d.MatRes_SamplerElem, 0, horde3d.MatRes_SampTexResI,
		int(newText.textureRes.H3DRes))

	newText.rasterize()

	newText.overlay.Material.Load()
	return newText
}

func (t *Text) rasterize() {
	c := t.context
	back := image.NewRGBA(image.Rect(0, 0, t.area.PixelWidth(), t.area.PixelHeight()))
	draw.Draw(back, back.Bounds(), image.Transparent, image.ZP, draw.Src)

	c.SetClip(back.Bounds())
	c.SetDst(back)
	c.SetSrc(image.White)
	t.background = back

	//TODO: Make resolution independent
	err := errors.New("")
	pt := freetype.Pt(5, int(c.PointToFix32(t.size)>>8))
	for _, s := range t.text {
		_, err = c.DrawString(s, pt)
		if err != nil {
			RaiseError(err)
			return
		}
		pt.Y += c.PointToFix32(t.size * t.lineSpacing)
	}

	t.textureRes.SetData(t.background)
	t.textureRes.Load()

}

func (t *Text) Place() {
	t.overlay.Place()
}

func (t *Text) Overlay() *Overlay { return t.overlay }
func (t *Text) Text() []string    { return t.text }
func (t *Text) SetText(text []string) {
	t.text = text
	t.rasterize()
}

func (t *Text) Size() float64 { return t.size }
func (t *Text) SetSize(size float64) {
	t.size = size
	t.rasterize()
}

func (t *Text) LineSpacing() float64 { return t.lineSpacing }
func (t *Text) SetLineSpacing(spacing float64) {
	t.lineSpacing = spacing
	t.rasterize()
}

func (t *Text) FontFile() string { return t.fontFile }
func (t *Text) SetFontFile(file string) {
	t.fontFile = file
	fontData, err := loadEngineData(file)
	if err != nil {
		RaiseError(err)
		return
	}

	font, err := freetype.ParseFont(fontData)
	if err != nil {
		RaiseError(err)
		return
	}

	t.context.SetFont(font)
	t.rasterize()
}

func (t *Text) Color() *Color         { return t.overlay.Color }
func (t *Text) SetColor(color *Color) { t.overlay.Color = color }

func (t *Text) Area() *ScreenArea { return t.area }
func (t *Text) SetArea(area *ScreenArea) {
	t.area = area
	t.rasterize()
}

func (t *Text) Unload() {
	t.overlay.Material.Remove()
	t.textureRes.Remove()
}
