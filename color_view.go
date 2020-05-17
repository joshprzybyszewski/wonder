// +build js,wasm

package wonder

import (
	"image"
	"image/color"
	"syscall/js"
	"time"

	"github.com/anthonynsimon/bild/adjust"
	"github.com/anthonynsimon/bild/clone"
)

const (
	normalViewStyle    = `norm`
	redOnlyViewStyle   = `ro`
	greenOnlyViewStyle = `go`
	blueOnlyViewStyle  = `bo`
	shiftingViewStyle  = `sc`
)

var (
	colorViewStyle = normalViewStyle
)

func isValidViewStyle(val string) bool {
	switch val {
	case normalViewStyle,
		redOnlyViewStyle, greenOnlyViewStyle, blueOnlyViewStyle,
		shiftingViewStyle:
		return true
	}
	return false
}

func (w *Wonder) setupColorViewCb() {
	w.colorViewCb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// quick return if no source image is yet uploaded
		val := args[0].Get("target").Get("value").String()
		if !isValidViewStyle(val) {
			w.log(`developer error: viewStyle ` + val + ` is invalid`)
			return nil
		}
		colorViewStyle = val

		w.updateImageForViewStyle()

		args[0].Call("preventDefault")
		return nil
	})
}

func (w *Wonder) updateImageForViewStyle() {
	if w.sourceImg != nil {
		start := time.Now()
		var res *image.RGBA

		switch colorViewStyle {
		case normalViewStyle:
			res = clone.AsRGBA(w.sourceImg)
		default:
			res = adjust.Apply(w.sourceImg, applyViewStyle)
		}

		w.updateImage(res, start)
	}
}

func applyViewStyle(input color.RGBA) color.RGBA {
	switch colorViewStyle {
	case redOnlyViewStyle:
		return color.RGBA{input.R, 0, 0, input.A}
	case greenOnlyViewStyle:
		return color.RGBA{0, input.G, 0, input.A}
	case blueOnlyViewStyle:
		return color.RGBA{0, 0, input.B, input.A}

	case normalViewStyle, shiftingViewStyle:
		fallthrough
	default:
		return input // color.RGBA{input.R, input.G, input.B, input.A}
	}
}
