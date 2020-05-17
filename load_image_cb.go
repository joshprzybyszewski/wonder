// +build js,wasm

package wonder

import (
	"bytes"
	"image"
	"syscall/js"
)

func (w *Wonder) setupOnImgLoadCb() {
	w.onImgLoadCb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		array := args[0]
		w.inBuf = make([]uint8, array.Get("byteLength").Int())
		js.CopyBytesToGo(w.inBuf, array)

		reader := bytes.NewReader(w.inBuf)
		var err error
		w.sourceImg, _, err = image.Decode(reader)
		if err != nil {
			w.log(err.Error())
			return nil
		}
		w.log("Ready for operations")

		// reset brightness and contrast sliders
		// js.Global().Get("document").
		// 	Call("getElementById", "brightness").
		// 	Set("value", 0)

		// js.Global().Get("document").
		// 	Call("getElementById", "contrast").
		// 	Set("value", 0)
		return nil
	})
}
