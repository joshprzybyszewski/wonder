// +build js,wasm

package wonder

import (
	"syscall/js"
	"time"

	"github.com/anthonynsimon/bild/adjust"
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
		args[0].Call("preventDefault")
		return nil
	})
}

func (w *Wonder) setupBrightnessCb() {
	w.brightnessCb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// quick return if no source image is yet uploaded
		if w.sourceImg == nil {
			return nil
		}
		delta := args[0].Get("target").Get("valueAsNumber").Float()
		start := time.Now()
		res := adjust.Brightness(w.sourceImg, delta)
		w.updateImage(res, start)
		args[0].Call("preventDefault")
		return nil
	})
}

func (w *Wonder) setupContrastCb() {
	w.contrastCb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// quick return if no source image is yet uploaded
		if w.sourceImg == nil {
			return nil
		}
		delta := args[0].Get("target").Get("valueAsNumber").Float()
		start := time.Now()
		res := adjust.Contrast(w.sourceImg, delta)
		w.updateImage(res, start)
		args[0].Call("preventDefault")
		return nil
	})
}

func (w *Wonder) setupHueCb() {
	w.hueCb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// quick return if no source image is yet uploaded
		if w.sourceImg == nil {
			return nil
		}
		delta := args[0].Get("target").Get("valueAsNumber").Int()
		start := time.Now()
		res := adjust.Hue(w.sourceImg, delta)
		w.updateImage(res, start)
		args[0].Call("preventDefault")
		return nil
	})
}

func (w *Wonder) setupSatCb() {
	w.satCb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// quick return if no source image is yet uploaded
		if w.sourceImg == nil {
			return nil
		}
		delta := args[0].Get("target").Get("valueAsNumber").Float()
		start := time.Now()
		res := adjust.Saturation(w.sourceImg, delta)
		w.updateImage(res, start)
		args[0].Call("preventDefault")
		return nil
	})
}
