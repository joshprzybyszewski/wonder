// +build js,wasm

package wonder

import (
	"bytes"
	"image"
	"strconv"
	"syscall/js"
	"time"

	"github.com/anthonynsimon/bild/imgio"
)

type Shimmer struct {
	inBuf                    []uint8
	outBuf                   bytes.Buffer
	onImgLoadCb, shutdownCb  js.Func
	brightnessCb, contrastCb js.Func
	hueCb, satCb             js.Func
	sourceImg                image.Image

	console js.Value
	done    chan struct{}
}

// New returns a new instance of shimmer
func New() *Shimmer {
	return &Shimmer{
		console: js.Global().Get("console"),
		done:    make(chan struct{}),
	}
}

// Start sets up all the callbacks and waits for the close signal
// to be sent from the browser.
func (s *Shimmer) Start() {
	// Setup callbacks
	s.setupOnImgLoadCb()
	js.Global().Set("loadImage", s.onImgLoadCb)

	s.setupBrightnessCb()
	js.Global().Get("document").
		Call("getElementById", "brightness").
		Call("addEventListener", "change", s.brightnessCb)

	s.setupContrastCb()
	js.Global().Get("document").
		Call("getElementById", "contrast").
		Call("addEventListener", "change", s.contrastCb)

	s.setupHueCb()
	js.Global().Get("document").
		Call("getElementById", "hue").
		Call("addEventListener", "change", s.hueCb)

	s.setupSatCb()
	js.Global().Get("document").
		Call("getElementById", "sat").
		Call("addEventListener", "change", s.satCb)

	s.setupShutdownCb()
	js.Global().Get("document").
		Call("getElementById", "close").
		Call("addEventListener", "click", s.shutdownCb)

	<-s.done
	s.log("Shutting down app")
	s.onImgLoadCb.Release()
	s.brightnessCb.Release()
	s.contrastCb.Release()
	s.hueCb.Release()
	s.satCb.Release()
	s.shutdownCb.Release()
}

// updateImage writes the image to a byte buffer and then converts it to base64.
// Then it sets the value to the src attribute of the target image.
func (s *Shimmer) updateImage(img *image.RGBA, start time.Time) {
	enc := imgio.JPEGEncoder(90)
	err := enc(&s.outBuf, img)
	if err != nil {
		s.log(err.Error())
		return
	}

	dst := js.Global().Get("Uint8Array").New(len(s.outBuf.Bytes()))
	n := js.CopyBytesToJS(dst, s.outBuf.Bytes())
	s.console.Call("log", "bytes copied:", strconv.Itoa(n))
	js.Global().Call("displayImage", dst)
	s.console.Call("log", "time taken:", time.Now().Sub(start).String())
	s.outBuf.Reset()
}

// utility function to log a msg to the UI from inside a callback
func (s *Shimmer) log(msg string) {
	js.Global().Get("document").
		Call("getElementById", "status").
		Set("innerText", msg)
}

func (s *Shimmer) setupShutdownCb() {
	s.shutdownCb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		s.done <- struct{}{}
		return nil
	})
}
