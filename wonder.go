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

type Wonder struct {
	inBuf                   []uint8
	outBuf                  bytes.Buffer
	onImgLoadCb, shutdownCb js.Func
	onProcessVideoStream    js.Func
	colorViewCb             js.Func
	sourceImg               image.Image

	console js.Value
	done    chan struct{}
}

// New returns a new instance of shimmer
func New() *Wonder {
	return &Wonder{
		console: js.Global().Get("console"),
		done:    make(chan struct{}),
	}
}

// Start sets up all the callbacks and waits for the close signal
// to be sent from the browser.
func (w *Wonder) Start() {
	// Setup callbacks
	w.setupOnImgLoadCb()
	w.setupProcessVideoStream()

	js.Global().Set("loadImage", w.onImgLoadCb)
	js.Global().Set("processVideoStream", w.onProcessVideoStream)

	w.setupColorViewCb()
	js.Global().Get("document").
		Call("getElementById", "colorViewRadioForm").
		Call("addEventListener", "change", w.colorViewCb)

	w.setupShutdownCb()
	js.Global().Get("document").
		Call("getElementById", "close").
		Call("addEventListener", "click", w.shutdownCb)

	<-w.done
	w.log("Shutting down app")
	w.onImgLoadCb.Release()
	w.onProcessVideoStream.Release()
	w.colorViewCb.Release()
	w.shutdownCb.Release()
}

// updateImage writes the image to a byte buffer and then converts it to base64.
// Then it sets the value to the src attribute of the target image.
func (w *Wonder) updateImage(img *image.RGBA, start time.Time) {
	enc := imgio.JPEGEncoder(90)
	err := enc(&w.outBuf, img)
	if err != nil {
		w.log(err.Error())
		return
	}

	dst := js.Global().Get("Uint8Array").New(len(w.outBuf.Bytes()))
	n := js.CopyBytesToJS(dst, w.outBuf.Bytes())
	w.console.Call("log", "bytes copied:", strconv.Itoa(n))
	js.Global().Call("displayImage", dst)
	w.console.Call("log", "time taken:", time.Now().Sub(start).String())
	w.outBuf.Reset()
}

// utility function to log a msg to the UI from inside a callback
func (w *Wonder) log(msg string) {
	js.Global().Get("document").
		Call("getElementById", "status").
		Set("innerText", msg)
}

func (w *Wonder) setupShutdownCb() {
	w.shutdownCb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		w.done <- struct{}{}
		return nil
	})
}
