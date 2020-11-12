package main

import (
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/eps"
	"github.com/tdewolff/canvas/pdf"
	"github.com/tdewolff/canvas/rasterizer"
	"github.com/tdewolff/canvas/svg"
)

const (
	font_family = "Montserrat"
	font_file_regular = "fonts/Montserrat-Regular.ttf"
	font_file_bold = "fonts/Montserrat-Bold.ttf"
	text_content = "Origin from Wikidata\nElian Carsenat, 11-2020 (1/15)"
	page_width  float64 = 297 //210 // A4 = 297 x 210
	page_height float64 = 210 //297 // A4 = 297 x 210

	tally_width = page_height
	tally_height = page_width/2
	tally_x = (page_width-tally_width)/2
	tally_y = (page_height-tally_height)/2
	h_blocks = 9
	v_blocks = 7
	block_width = tally_width/h_blocks
	block_height = tally_height/v_blocks
)

func main() {
	c := canvas.New(page_width, page_height)
	ctx := canvas.NewContext(c)
	fontFamily := canvas.NewFontFamily(font_family)
	fontFamily.Use(canvas.CommonLigatures)
	if err := fontFamily.LoadFontFile(font_file_regular, canvas.FontRegular); err != nil {
		panic(err)
	}
	if err := fontFamily.LoadFontFile(font_file_bold, canvas.FontBold); err != nil {
		panic(err)
	}


	ctx.SetFillColor(canvas.Lightgray)
	for i := 0; i < h_blocks; i++ {
		for j := 0; j < v_blocks; j++ {
			if (i+j)%2 == 0 {
				ctx.DrawPath(tally_x+float64(i)*block_width, tally_y+float64(j)*block_height, canvas.RoundedRectangle(block_width, block_height, 1))
			}
		}
	}

	// Draw a comprehensive text box
	fontSize := 8.0
	face := fontFamily.Face(fontSize, canvas.Black, canvas.FontBold, canvas.FontNormal)
	rich := canvas.NewRichText()
	rich.Add(face, text_content)
	//metrics := face.Metrics()

	text := rich.ToText(3*block_width, block_height, canvas.Center, canvas.Center, 0.0, 0.0)

	ctx.SetFillColor(canvas.Lightblue)
	ctx.DrawPath(tally_x+3.0*block_width, tally_y+3.0*block_height, canvas.RoundedRectangle(3*block_width, block_height, 1))

	ctx.DrawText(tally_x+3.0*block_width,tally_y+4.0*block_height, text)

	//p:=canvas.MustParseSVG(qrSVG)
	//ctx.DrawPath(10,10,p)
	c.WriteFile("canvas_out.svg", svg.Writer)
	c.WriteFile("canvas_out.pdf", pdf.Writer)
	c.WriteFile("canvas_out.eps", eps.Writer)
	c.WriteFile("canvas_out.png", rasterizer.PNGWriter(3.2))
}


