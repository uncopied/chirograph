package tallystick

import (
	"github.com/boombuler/barcode/code128"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/pdf"
	"github.com/tdewolff/canvas/svg"
	qrcode "github.com/uncopied/go-qrcode"
	"io"
	"log"
	"math"
	"math/rand"
)

const (
	assetsFontsPath = "assets/fonts"
	fontFamily                 = "Montserrat"
	fontFileRegular            = assetsFontsPath+"/Montserrat-Regular.ttf"
	fontFileBold               = assetsFontsPath+"/Montserrat-Bold.ttf"
	fontSizeNormal             = 8.0
	fontSizeSmall              = 6.0
	pageWidth          float64 = 297 //210 // A4 = 297 x 210
	pageHeight         float64 = 210 //297 // A4 = 297 x 210

	// random cutline bands
	randCutWidth=0.1
	randCutHeight=0.1
	randSideStep=0.1
	cutLineWidth=0.2

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
	tallyXPortrait      = -pageHeight + tallyY
	tallyYPortrait      = tallyX
	blockWidthPortrait  = blockHeight
	blockHeightPortrait = blockWidth

	// qrblock margin
	qrBlockInnerMargin = 0.2
	// code128 margin
	codeBlockInnerMargin = 0.1
	// Draw tally
	drawTallyYesNo = true
	// Draw grids
	drawGrid = false
	// colorize
	colorize = false
)

type Tallystick struct {

	CertificateLabel string

	PrimaryLinkURL string
	SecondaryLinkURL string

	IssuerTokenURL string
	OwnerTokenURL string
	PrimaryAssetVerifierTokenURL string
	SecondaryAssetVerifierTokenURL string
	PrimaryOwnerVerifierTokenURL string
	SecondaryOwnerVerifierTokenURL string
	PrimaryIssuerVerifierTokenURL string
	SecondaryIssuerVerifierTokenURL string

	MailToContentLeft string
	MailToContentRight string

}

func drawUncopiedLogo(fontFamily *canvas.FontFamily, ctx *canvas.Context, hBlock float64, vBlock float64, rotate bool, mailToContent string ) {
	// Draw a comprehensive text box
	face := fontFamily.Face(fontSizeSmall, canvas.Black, canvas.FontBold, canvas.FontNormal)
	rich := canvas.NewRichText()
	rich.Add(face, mailToContent)
	//metrics := face.Metrics()

	if colorize {
		ctx.SetFillColor(canvas.Lightpink)
	} else {
		ctx.SetFillColor(canvas.White)
	}
	if rotate {
		// portrait
		ctx.Rotate(-90)
		text := rich.ToText(blockWidthPortrait, blockHeightPortrait, canvas.Center, canvas.Center, 0.0, 0.0)
		ctx.DrawPath(tallyXPortrait+(vBlocks-vBlock-1)*blockWidthPortrait, tallyYPortrait+hBlock*blockHeightPortrait, canvas.RoundedRectangle(blockWidthPortrait, blockHeightPortrait, 1))
		ctx.DrawText(tallyXPortrait+(vBlocks-vBlock-1)*blockWidthPortrait, tallyYPortrait+(hBlock+1)*blockHeightPortrait, text)
		ctx.Rotate(90)
	} else {
		// landscape
		text := rich.ToText(blockWidth, blockHeight, canvas.Center, canvas.Center, 0.0, 0.0)
		ctx.DrawPath(tallyX+hBlock*blockWidth, tallyY+vBlock*blockHeight, canvas.RoundedRectangle(blockWidth, blockHeight, 1))
		ctx.DrawText(tallyX+hBlock*blockWidth, tallyY+(vBlock+1)*blockHeight, text)
	}
}

func drawText(fontFamily *canvas.FontFamily, ctx *canvas.Context, hBlock float64, vBlock float64, wBlock float64, textContent string, rotate bool, hAlign canvas.TextAlign, vAlign canvas.TextAlign, whiteOrTransparent bool, fontSize float64) {
	// Draw a comprehensive text box
	face := fontFamily.Face(fontSize, canvas.Black, canvas.FontBold, canvas.FontNormal)
	rich := canvas.NewRichText()
	rich.Add(face, textContent)
	//metrics := face.Metrics()
	if whiteOrTransparent {
		ctx.SetFillColor(canvas.White)
	} else {
		ctx.SetFillColor(canvas.Transparent)
	}
	if rotate {
		ctx.Rotate(-90)
		text := rich.ToText(wBlock*blockWidthPortrait, blockHeightPortrait, hAlign, vAlign, 0.0, 0.0)
		ctx.DrawPath(tallyXPortrait+(vBlocks-vBlock-1)*blockWidthPortrait, tallyYPortrait+hBlock*blockHeightPortrait, canvas.RoundedRectangle(blockWidthPortrait*wBlock, blockHeightPortrait, 1))
		ctx.DrawText(tallyXPortrait+(vBlocks-vBlock-1)*blockWidthPortrait, tallyYPortrait+(hBlock+1)*blockHeightPortrait, text)
		ctx.Rotate(90)
	} else {
		text := rich.ToText(wBlock*blockWidth, blockHeight,  hAlign, vAlign, 0.0, 0.0)
		ctx.DrawPath(tallyX+hBlock*blockWidth, tallyY+vBlock*blockHeight, canvas.RoundedRectangle(blockWidth*wBlock, blockHeight, 1))
		ctx.DrawText(tallyX+hBlock*blockWidth, tallyY+(vBlock+1)*blockHeight, text)
	}
}

func drawBarCode128(ctx *canvas.Context, hBlock float64, vBlock float64, wBlock float64, widthRatio float64,  textContent string, rotate bool ) {
	barCode, _ := code128.Encode(textContent)
	if rotate {
		ctx.Rotate(-90)
		pixelWidth := blockWidthPortrait *wBlock/float64(barCode.Bounds().Max.X)
		for i := 0; i < barCode.Bounds().Max.X; i++ {
			c := barCode.At(i, 0)
			ctx.SetFillColor(c)
			ctx.DrawPath(tallyXPortrait+(vBlocks-vBlock-1)*blockWidthPortrait+float64(i)*pixelWidth, tallyYPortrait+hBlock*blockHeightPortrait, canvas.Rectangle(pixelWidth, blockHeightPortrait*widthRatio))
		}
		ctx.Rotate(90)
	} else {
		pixelWidth := blockWidth *wBlock/float64(barCode.Bounds().Max.X)
		for i := 0; i < barCode.Bounds().Max.X; i++ {
			c := barCode.At(i, 0)
			ctx.SetFillColor(c)
			ctx.DrawPath(tallyX+hBlock*blockWidth+float64(i)*pixelWidth, tallyY+vBlock*blockHeight, canvas.Rectangle(pixelWidth, blockHeight*widthRatio))
		}
	}
}

func addRandomCenteredPoint(p *canvas.Polyline, hBlock float64, vBlock float64, topOrRight bool) [2]float64 {
	var xy [2]float64
	xy[0] = tallyX +hBlock*blockWidth +(0.5-randCutWidth*rand.Float64())*blockWidth
	xy[1] = tallyY +vBlock*blockHeight +(0.5-randCutHeight*rand.Float64())*blockHeight
	p.Add(xy[0], xy[1])
	return xy
}

func addRandomZigZagPoints(p *canvas.Polyline, hBlock float64, vBlock float64, topOrRight bool) [2]float64 {
	if topOrRight {
		var xy [2]float64
		xy[0] = tallyX +hBlock*blockWidth +(0.5+randSideStep+randCutWidth*randCutWidth*rand.Float64())*blockWidth
		xy[1] = tallyY +vBlock*blockHeight +(0.25-randCutHeight*rand.Float64())*blockHeight
		p.Add(xy[0], xy[1])
		xy[0] = tallyX +hBlock*blockWidth +(0.5-randSideStep-randCutWidth*randCutWidth*rand.Float64())*blockWidth
		xy[1] = tallyY +vBlock*blockHeight +(0.75+randCutHeight*rand.Float64())*blockHeight
		p.Add(xy[0], xy[1])
		return xy
	} else {
		var xy [2]float64
		xy[0] = tallyX +hBlock*blockWidth +(0.25-randCutWidth*rand.Float64())*blockWidth
		xy[1] = tallyY +vBlock*blockHeight +(0.5+randSideStep+randCutHeight*rand.Float64())*blockHeight
		p.Add(xy[0], xy[1])
		xy[0] = tallyX +hBlock*blockWidth +(0.75+randCutWidth*rand.Float64())*blockWidth
		xy[1] = tallyY +vBlock*blockHeight +(0.5-randSideStep-randCutHeight*rand.Float64())*blockHeight
		p.Add(xy[0], xy[1])
		return xy
	}
}

func drawCutLine(ctx *canvas.Context) {
	ctx.SetFillColor(canvas.Transparent)
	ctx.SetStrokeColor(canvas.Blue)
	ctx.SetStrokeWidth(cutLineWidth)

	// Draw around
	polyline := &canvas.Polyline{}
	polyline.Add(tallyX, tallyY)
	polyline.Add(tallyX+hBlocks*blockWidth, tallyY)
	polyline.Add(tallyX+hBlocks*blockWidth, tallyY+vBlocks*blockHeight)
	polyline.Add(tallyX, tallyY+vBlocks*blockHeight)
	polyline.Add(tallyX, tallyY)
	ctx.DrawPath(0, 0, polyline.ToPath())
	// create the 4 random points

	//Draw left
	polyline = &canvas.Polyline{}
	polyline.Add(tallyX+blockWidth*1.5, tallyY)
	addRandomZigZagPoints(polyline, 1,0, true)
	addRandomCenteredPoint(polyline, 1,1,true)
	cbd1 := addRandomZigZagPoints(polyline, 1,2, true)
	addRandomCenteredPoint(polyline, 1,3,true)
	abd1 := addRandomZigZagPoints(polyline, 1,4, true)
	addRandomCenteredPoint(polyline, 1,5,true)
	addRandomZigZagPoints(polyline, 1,6,true)
	polyline.Add(tallyX+blockWidth*1.5, tallyY+blockHeight*vBlocks)
	ctx.DrawPath(0, 0, polyline.ToPath())

	//Draw right
	polyline = &canvas.Polyline{}
	polyline.Add(tallyX+blockWidth*7.5, tallyY)
	addRandomZigZagPoints(polyline, 7,0,true)
	addRandomCenteredPoint(polyline, 7,1,true)
	cbd2 := addRandomZigZagPoints(polyline,7,2, true)
	addRandomCenteredPoint(polyline, 7,3,true)
	abd2 := addRandomZigZagPoints(polyline,7,4, true)
	addRandomCenteredPoint(polyline, 7,5,true)
	addRandomZigZagPoints(polyline, 7,6,true)
	polyline.Add(tallyX+blockWidth*7.5, tallyY+blockHeight*vBlocks)
	ctx.DrawPath(0, 0, polyline.ToPath())

	//Draw bottom
	polyline = &canvas.Polyline{}
	polyline.Add(cbd1[0], cbd1[1])
	addRandomZigZagPoints(polyline, 2,2,false)
	addRandomZigZagPoints(polyline, 3,2,false)
	addRandomCenteredPoint(polyline,4,2, false)
	addRandomZigZagPoints(polyline,5,2, false)
	addRandomZigZagPoints(polyline,6,2, false)
	polyline.Add(cbd2[0], cbd2[1])
	ctx.DrawPath(0, 0, polyline.ToPath())

	//Draw bottom
	polyline = &canvas.Polyline{}
	polyline.Add(abd1[0], abd1[1])
	addRandomZigZagPoints(polyline, 2,4,false)
	addRandomZigZagPoints(polyline, 3,4,false)
	addRandomCenteredPoint(polyline,4,4, false)
	addRandomZigZagPoints(polyline,5,4, false)
	addRandomZigZagPoints(polyline,6,4, false)
	polyline.Add(abd2[0], abd2[1])
	ctx.DrawPath(0, 0, polyline.ToPath())

}

func drawQRCode(ctx *canvas.Context, hBlock float64, vBlock float64, content string, rotate bool ) {
	q, err := qrcode.New(content, qrcode.Highest)
	if err != nil {
		log.Fatal(err)
	}
	if rotate {
		ctx.Rotate(-90)
		innerSquareWidth := math.Min(blockWidthPortrait, blockHeightPortrait) * (1 - qrBlockInnerMargin)
		hBlockMargin := (blockWidthPortrait - innerSquareWidth) / 2
		vBlockMargin := (blockHeightPortrait - innerSquareWidth) / 2
		q.DrawQRCode(ctx, tallyXPortrait+(vBlocks-vBlock-1)*blockWidthPortrait+hBlockMargin, tallyYPortrait+hBlock*blockHeightPortrait+vBlockMargin, innerSquareWidth)
		ctx.Rotate(90)
	} else {
		innerSquareWidth := math.Min(blockWidth, blockHeight) * (1 - qrBlockInnerMargin)
		hBlockMargin := (blockWidth - innerSquareWidth) / 2
		vBlockMargin := (blockHeight - innerSquareWidth) / 2
		q.DrawQRCode(ctx, tallyX+hBlock*blockWidth+hBlockMargin, tallyY+vBlock*blockHeight+vBlockMargin, innerSquareWidth)
	}
}

func drawTally(fontFamily *canvas.FontFamily, ctx *canvas.Context, t *Tallystick) {

	drawBarCode128(ctx,0,6,7, 2-codeBlockInnerMargin, t.PrimaryLinkURL,true)
	drawBarCode128(ctx,7+codeBlockInnerMargin,6,7, 2-codeBlockInnerMargin, t.PrimaryLinkURL,true)
	drawBarCode128(ctx,2,0,5, 1, t.PrimaryLinkURL,false)
	drawBarCode128(ctx,2,2,5, 1, t.PrimaryLinkURL,false)
	drawBarCode128(ctx,2,4,5, 1, t.PrimaryLinkURL,false)
	drawBarCode128(ctx,2,6,5, 1, t.PrimaryLinkURL,false)

	drawText(fontFamily, ctx, 3, 1, 3, t.CertificateLabel, false, canvas.Center, canvas.Center, true, fontSizeNormal)
	drawText(fontFamily, ctx, 3, 3, 3, t.CertificateLabel, false, canvas.Center, canvas.Center, true, fontSizeNormal)
	drawText(fontFamily, ctx, 3, 5, 3, t.CertificateLabel,false, canvas.Center, canvas.Center, true, fontSizeNormal)

	drawText(fontFamily, ctx, 0, 4, 3, t.CertificateLabel, true, canvas.Center, canvas.Center,true, fontSizeNormal)
	drawText(fontFamily, ctx, 8, 4, 3, t.CertificateLabel, true, canvas.Center, canvas.Center,true, fontSizeNormal)


	drawQRCode(ctx, 2, 5, t.PrimaryLinkURL, false)
	drawQRCode(ctx, 6, 5, t.SecondaryLinkURL, false)

	drawQRCode(ctx, 4, 4, t.IssuerTokenURL, false)

	drawQRCode(ctx, 2, 3, t.PrimaryLinkURL, false)
	drawQRCode(ctx, 6, 3, t.SecondaryLinkURL, false)

	drawQRCode(ctx, 4, 2, t.OwnerTokenURL, false)

	drawQRCode(ctx, 2, 1, t.PrimaryLinkURL, false)
	drawQRCode(ctx, 6, 1, t.SecondaryLinkURL, false)

	drawQRCode(ctx, 1, 5, t.PrimaryIssuerVerifierTokenURL, false)
	drawQRCode(ctx, 1, 3, t.PrimaryAssetVerifierTokenURL, false)
	drawQRCode(ctx, 1, 1, t.PrimaryOwnerVerifierTokenURL, false)

	drawQRCode(ctx, 7, 5, t.SecondaryIssuerVerifierTokenURL, false)
	drawQRCode(ctx, 7, 3, t.SecondaryAssetVerifierTokenURL, false)
	drawQRCode(ctx, 7, 1, t.SecondaryOwnerVerifierTokenURL, false)

	drawQRCode(ctx, 0, 6, t.PrimaryLinkURL, true)
	drawQRCode(ctx, 0, 0, t.SecondaryLinkURL, true)

	drawQRCode(ctx, 8, 6, t.PrimaryLinkURL, true)
	drawQRCode(ctx, 8, 0, t.SecondaryLinkURL, true)

	drawUncopiedLogo(fontFamily, ctx, 0, 1, false, t.MailToContentLeft)
	drawUncopiedLogo(fontFamily, ctx, 0, 5, true, t.MailToContentLeft)
	drawUncopiedLogo(fontFamily, ctx, 8, 1, false,t.MailToContentRight)
	drawUncopiedLogo(fontFamily, ctx, 8, 5, true, t.MailToContentRight)

	drawText(fontFamily, ctx, 2, 1, 1, t.PrimaryLinkURL, false, canvas.Left, canvas.Top, false, fontSizeSmall)
	drawText(fontFamily, ctx, 2, 3, 1, t.PrimaryLinkURL, false, canvas.Left, canvas.Top,false, fontSizeSmall)
	drawText(fontFamily, ctx, 2, 5, 1, t.PrimaryLinkURL,false,  canvas.Left, canvas.Top,false, fontSizeSmall)

	drawText(fontFamily, ctx, 6, 1, 1, t.SecondaryLinkURL, false, canvas.Right, canvas.Bottom, false, fontSizeSmall)
	drawText(fontFamily, ctx, 6, 3, 1, t.SecondaryLinkURL, false, canvas.Right, canvas.Bottom,false, fontSizeSmall)
	drawText(fontFamily, ctx, 6, 5, 1, t.SecondaryLinkURL,false,  canvas.Right, canvas.Bottom,false, fontSizeSmall)

	drawCutLine(ctx)
}

func DrawSVG(t *Tallystick, w io.Writer) error {
	c := Draw(t)
	svg := svg.New(w, c.W, c.H)
	c.Render(svg)
	return svg.Close()
}

//c.WriteFile("canvas_out.svg", svg.Writer)
//c.WriteFile("canvas_out.pdf", pdf.Writer)
func DrawPDF(t *Tallystick, w io.Writer) error {
	c := Draw(t)
	pdf := pdf.New(w, c.W, c.H)
	c.Render(pdf)
	return pdf.Close()
}

func Draw(t *Tallystick) *canvas.Canvas {
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
	// paint a grid on landscape
	if drawGrid {
		ctx.SetFillColor(canvas.Lightgray)
		for i := 0; i < hBlocks; i++ {
			for j := 0; j < vBlocks; j++ {
				if (i+j)%2 == 0 {
					ctx.DrawPath(tallyX+float64(i)*blockWidth, tallyY+float64(j)*blockHeight, canvas.RoundedRectangle(blockWidth, blockHeight, 1))
				}
			}
		}

		ctx.SetFillColor(canvas.Darkcyan)
		ctx.Rotate(-90)
		// paint a grid on portrait
		for i := 0; i < vBlocks; i++ {
			for j := 0; j < hBlocks; j++ {
				if (i+j)%2 == 0 {
					//ctx.DrawPath(tallyX+float64(i)*blockWidth, tallyY+float64(j)*blockHeight, canvas.RoundedRectangle(blockWidth, blockHeight, 1))
					ctx.DrawPath(tallyXPortrait+float64(i)*blockWidthPortrait+1, tallyYPortrait+1+float64(j)*blockHeightPortrait, canvas.RoundedRectangle(blockWidthPortrait-2, blockHeightPortrait-2, 1))
				}
			}
		}
		ctx.Rotate(90)
	}

	if drawTallyYesNo {
		drawTally(fontFamily,ctx,t)
	}
	return c
	//p:=canvas.MustParseSVG(qrSVG)
	//ctx.DrawPath(10,10,p)
	//c.WriteFile("canvas_out.svg", svg.Writer)
	//c.WriteFile("canvas_out.pdf", pdf.Writer)
	//c.WriteFile("canvas_out.eps", eps.Writer)
	//c.WriteFile("canvas_out.png", rasterizer.PNGWriter(3.2))
}
