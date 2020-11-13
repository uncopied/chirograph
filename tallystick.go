package main

import (
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/eps"
	"github.com/tdewolff/canvas/pdf"
	"github.com/tdewolff/canvas/rasterizer"
	"github.com/tdewolff/canvas/svg"
	qrcode "github.com/uncopied/go-qrcode"
	"log"
	"math"
)

const (
	fontFamily              = "Montserrat"
	fontFileRegular         = "fonts/Montserrat-Regular.ttf"
	fontFileBold            = "fonts/Montserrat-Bold.ttf"
	fontSize                = 8.0
	pageWidth       float64 = 297 //210 // A4 = 297 x 210
	pageHeight      float64 = 210 //297 // A4 = 297 x 210

	// blocks
	hBlocks     = 9
	vBlocks     = 7

	// landscape
	tallyWidth  = pageHeight
	tallyHeight = pageWidth / 2
	tallyX      = (pageWidth - tallyWidth) / 2
	tallyY      = (pageHeight - tallyHeight) / 2
	blockWidth  = tallyWidth / hBlocks
	blockHeight = tallyHeight / vBlocks

	// portrait
	tallyWidthPortrait  = tallyHeight
	tallyHeightPortrait = tallyWidth
	tallyXPortrait      = tallyY
	tallyYPortrait      = tallyX
	blockWidthPortrait  = blockHeight
	blockHeightPortrait = blockWidth

	blockInnerMargin = 0.05
)

func drawUncopiedLogo(fontFamily *canvas.FontFamily, ctx *canvas.Context, hBlock float64, vBlock float64) {
	// Draw a comprehensive text box
	face := fontFamily.Face(fontSize, canvas.Black, canvas.FontBold, canvas.FontNormal)
	rich := canvas.NewRichText()
	rich.Add(face, "un©opied\n"+
		"uncopied.art	")
	//metrics := face.Metrics()

	text := rich.ToText(blockWidth, blockHeight, canvas.Center, canvas.Center, 0.0, 0.0)

	ctx.SetFillColor(canvas.Lightpink)
	ctx.DrawPath(tallyX+hBlock*blockWidth, tallyY+vBlock*blockHeight, canvas.RoundedRectangle(blockWidth, blockHeight, 1))
	ctx.DrawText(tallyX+hBlock*blockWidth, tallyY+(vBlock+1)*blockHeight, text)
}

func drawHText(fontFamily *canvas.FontFamily, ctx *canvas.Context, hBlock float64, vBlock float64, textContent string) {
	// Draw a comprehensive text box
	face := fontFamily.Face(fontSize, canvas.Black, canvas.FontBold, canvas.FontNormal)
	rich := canvas.NewRichText()
	rich.Add(face, textContent)
	//metrics := face.Metrics()

	text := rich.ToText(3*blockWidth, blockHeight, canvas.Center, canvas.Center, 0.0, 0.0)

	ctx.SetFillColor(canvas.Lightblue)
	ctx.DrawPath(tallyX+hBlock*blockWidth, tallyY+vBlock*blockHeight, canvas.RoundedRectangle(blockWidth*3, blockHeight, 1))
	ctx.DrawText(tallyX+hBlock*blockWidth, tallyY+(vBlock+1)*blockHeight, text)
}

func drawVText(fontFamily *canvas.FontFamily, ctx *canvas.Context, hBlock float64, vBlock float64, textContent string) {
	// Draw a comprehensive text box
	face := fontFamily.Face(fontSize, canvas.Black, canvas.FontBold, canvas.FontNormal)
	text := canvas.NewTextBox(face,textContent,3*blockHeightPortrait, blockWidthPortrait,canvas.Center, canvas.Center, 0,0)

	ctx.RotateAbout(+90,pageWidth / 2, pageHeight / 2)
	ctx.SetFillColor(canvas.Red)
	ctx.DrawPath(tallyXPortrait*2+hBlock*blockWidthPortrait, tallyYPortrait/2+vBlock*blockHeightPortrait, canvas.RoundedRectangle(blockWidthPortrait*3, blockHeightPortrait, 1))
	ctx.DrawText(tallyXPortrait*2+hBlock*blockWidthPortrait, tallyYPortrait/2+(vBlock+1)*blockHeightPortrait, text)
	ctx.RotateAbout(-90,pageWidth / 2, pageHeight / 2)
}

func drawQRCode(ctx *canvas.Context, hBlock float64, vBlock float64, content string) {
	innerSquareWidth := math.Min(blockWidth, blockHeight) * (1 - blockInnerMargin)
	hBlockMargin := (blockWidth - innerSquareWidth) / 2
	vBlockMargin := (blockHeight - innerSquareWidth) / 2
	//ctx.SetFillColor(canvas.Red)
	//ctx.DrawPath(tallyX+float64(hBlock)*blockWidth+hBlockMargin, tallyY+float64(vBlock)*blockHeight+vBlockMargin, canvas.Rectangle(innerSquareWidth, innerSquareWidth))

	q, err := qrcode.New(content, qrcode.Highest)
	if err != nil {
		log.Fatal(err)
	}
	q.DrawQRCode(ctx, tallyX+hBlock*blockWidth+hBlockMargin, tallyY+vBlock*blockHeight+vBlockMargin, innerSquareWidth)
}

func main() {
	c := canvas.New(pageWidth, pageHeight)
	ctx := canvas.NewContext(c)
	fontFamily := canvas.NewFontFamily(fontFamily)
	fontFamily.Use(canvas.CommonLigatures)
	if err := fontFamily.LoadFontFile(fontFileRegular, canvas.FontRegular); err != nil {
		panic(err)
	}
	if err := fontFamily.LoadFontFile(fontFileBold, canvas.FontBold); err != nil {
		panic(err)
	}

	// show the grid
	ctx.SetFillColor(canvas.Lightgray)
	for i := 0; i < hBlocks; i++ {
		for j := 0; j < vBlocks; j++ {
			if (i+j)%2 == 0 {
				ctx.DrawPath(tallyX+float64(i)*blockWidth, tallyY+float64(j)*blockHeight, canvas.RoundedRectangle(blockWidth, blockHeight, 1))
			}
		}
	}

	drawUncopiedLogo(fontFamily, ctx, 0, 1)

	myTextContent := "Origin from Wikidata\nElian Carsenat, 11-2020 (1/15)"
	drawHText(fontFamily, ctx, 3, 1, myTextContent)
	drawHText(fontFamily, ctx, 3, 3, myTextContent)
	drawHText(fontFamily, ctx, 3, 5, myTextContent)

	drawVText(fontFamily, ctx, 0, 2, myTextContent)

	drawQRCode(ctx, 2, 3, "uncopied-B1")
	drawQRCode(ctx, 6, 3, "uncopied-B2")

	drawQRCode(ctx, 2, 5, "uncopied-A1")
	drawQRCode(ctx, 6, 5, "uncopied-A2")

	drawQRCode(ctx, 2, 1, "uncopied-C1")
	drawQRCode(ctx, 6, 1, "uncopied-C2")

	drawQRCode(ctx, 1, 5, "uncopied-D1A")
	drawQRCode(ctx, 1, 3, "uncopied-D1B")
	drawQRCode(ctx, 1, 1, "uncopied-D1C")

	drawQRCode(ctx, 7, 5, "uncopied-D2A")
	drawQRCode(ctx, 7, 3, "uncopied-D2B")
	drawQRCode(ctx, 7, 1, "uncopied-D2C")

	drawQRCode(ctx, 0, 6, "uncopied-D11")
	drawQRCode(ctx, 0, 0, "uncopied-D12")

	drawQRCode(ctx, 8, 6, "uncopied-D21")
	drawQRCode(ctx, 8, 0, "uncopied-D22")

	//p:=canvas.MustParseSVG(qrSVG)
	//ctx.DrawPath(10,10,p)
	c.WriteFile("canvas_out.svg", svg.Writer)
	c.WriteFile("canvas_out.pdf", pdf.Writer)
	c.WriteFile("canvas_out.eps", eps.Writer)
	c.WriteFile("canvas_out.png", rasterizer.PNGWriter(3.2))
}
