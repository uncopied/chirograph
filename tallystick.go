package main

import (
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/eps"
	"github.com/tdewolff/canvas/pdf"
	"github.com/tdewolff/canvas/rasterizer"
	"github.com/tdewolff/canvas/svg"
	"image/color"
)

const (
	page_width  float64 = 297 //210 // A4 = 297 x 210
	page_height float64 = 210 //297 // A4 = 297 x 210
)

func main() {
	c := canvas.New(page_width, page_height)
	ctx := canvas.NewContext(c)
	fontFamily := canvas.NewFontFamily("Kurinto Sans")
	fontFamily.Use(canvas.CommonLigatures)
	if err := fontFamily.LoadFontFile("fonts/KurintoSans-Rg.ttf", canvas.FontRegular); err != nil {
		panic(err)
	}


	ctx.SetFillColor(canvas.Blue)
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if (i+j)%2 == 0 {
				ctx.DrawPath(float64(i)*page_width/10, float64(j)*page_height/10, canvas.RoundedRectangle(page_width/10, page_height/10, 1))
			}
		}
	}

	// Draw a comprehensive text box
	fontSize := 14.0
	face := fontFamily.Face(fontSize, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	rich := canvas.NewRichText()
	rich.Add(face, "\"Lorem ")
	rich.Add(face, " dolor ")
	rich.Add(face, "ipsum")
	rich.Add(face, "\". Confiscator")
	rich.Add(face, " Curabitur mattis dui tellus vel.")
	rich.Add(fontFamily.Face(fontSize, canvas.Black, canvas.FontBold, canvas.FontNormal), " faux bold")
	rich.Add(face, " ")
	rich.Add(fontFamily.Face(fontSize, canvas.Black, canvas.FontItalic, canvas.FontNormal), "faux italic")
	rich.Add(face, " ")
	rich.Add(fontFamily.Face(fontSize, canvas.Black, canvas.FontRegular, canvas.FontNormal, canvas.FontUnderline), "underline")
	rich.Add(face, " ")
	rich.Add(fontFamily.Face(fontSize, canvas.White, canvas.FontRegular, canvas.FontNormal, canvas.FontDoubleUnderline), "double underline")
	rich.Add(face, " ")
	rich.Add(fontFamily.Face(fontSize, canvas.Black, canvas.FontRegular, canvas.FontNormal, canvas.FontSineUnderline), "sine")
	rich.Add(face, " ")
	rich.Add(fontFamily.Face(fontSize, canvas.Black, canvas.FontRegular, canvas.FontNormal, canvas.FontSawtoothUnderline), "sawtooth")
	rich.Add(face, " ")
	rich.Add(fontFamily.Face(fontSize, canvas.Black, canvas.FontRegular, canvas.FontNormal, canvas.FontDottedUnderline), "dotted")
	rich.Add(face, " ")
	rich.Add(fontFamily.Face(fontSize, canvas.Black, canvas.FontRegular, canvas.FontNormal, canvas.FontDashedUnderline), "dashed")
	rich.Add(face, " ")
	rich.Add(fontFamily.Face(fontSize, canvas.Black, canvas.FontRegular, canvas.FontNormal, canvas.FontOverline), "overline ")
	rich.Add(fontFamily.Face(fontSize, canvas.Black, canvas.FontItalic, canvas.FontNormal, canvas.FontStrikethrough, canvas.FontSineUnderline, canvas.FontOverline), "combi")
	rich.Add(face, ".")
	drawText(ctx, 5, 95, face, rich)
	//p:=canvas.MustParseSVG(qrSVG)
	//ctx.DrawPath(10,10,p)
	c.WriteFile("canvas_out.svg", svg.Writer)
	c.WriteFile("canvas_out.pdf", pdf.Writer)
	c.WriteFile("canvas_out.eps", eps.Writer)
	c.WriteFile("canvas_out.png", rasterizer.PNGWriter(3.2))
}


func drawText(c *canvas.Context, x, y float64, face canvas.FontFace, rich *canvas.RichText) {
	metrics := face.Metrics()
	width, height := 90.0, 35.0

	text := rich.ToText(width, height, canvas.Justify, canvas.Top, 0.0, 0.0)

	c.SetFillColor(color.RGBA{192, 0, 64, 255})
	c.DrawPath(x, y, text.Bounds().ToPath())
	c.SetFillColor(color.RGBA{50, 50, 50, 50})
	c.DrawPath(x, y, canvas.Rectangle(width, -metrics.LineHeight))
	c.SetFillColor(color.RGBA{0, 0, 0, 50})
	c.DrawPath(x, y+metrics.CapHeight-metrics.Ascent, canvas.Rectangle(width, -metrics.CapHeight-metrics.Descent))
	c.DrawPath(x, y+metrics.XHeight-metrics.Ascent, canvas.Rectangle(width, -metrics.XHeight))

	c.SetFillColor(canvas.Black)
	c.DrawPath(x, y, canvas.Rectangle(width, -height).Stroke(0.2, canvas.RoundCap, canvas.RoundJoin))
	c.DrawText(x, y, text)
}