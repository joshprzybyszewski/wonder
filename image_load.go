// +build js,wasm

package wonder

import (
	"bytes"
	"image"
	"strconv"
	"syscall/js"
	"time"
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

		return nil
	})
}

func (w *Wonder) setupProcessVideoStream() {
	w.onProcessVideoStream = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		w.log("got stream")

		v := &videoProcessor{}
		v.load()

		return nil
	})
}

type videoProcessor struct {
	width, height int
	video         js.Value
	ctx2d         js.Value
	console       js.Value

	dataBuf []uint8
}

func (v *videoProcessor) timerCB() {
	v.computeFrame()

	time.AfterFunc(16*time.Millisecond, v.timerCB)
}

func (v *videoProcessor) updateWidthHeight() {
	v.width = v.video.Get(`width`).Int()
	v.height = v.video.Get(`height`).Int()
}

func (v *videoProcessor) load() {
	v.console = js.Global().Get("console")
	v.video = js.Global().Get("document").
		Call("getElementById", "videoElement")

	v.updateWidthHeight()

	v.ctx2d = js.Global().Get("document").
		Call("getElementById", "videoCanvas").
		Call("getContext", "2d")

	v.video.Call("addEventListener", "play", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		v.updateWidthHeight()
		v.timerCB()

		return nil
	}), false)
}

func (v *videoProcessor) computeFrame() {
	if v.width <= 0 || v.height <= 0 {
		v.updateWidthHeight()
		return
	}
	v.ctx2d.Call(`drawImage`, v.video, 0, 0, v.width, v.height)

	frame := v.ctx2d.Call(`getImageData`, 0, 0, v.width, v.height)
	dataSlice := frame.Get(`data`)
	dataLen := dataSlice.Get(`length`).Int()
	if dataLen <= 0 {
		v.console.Call("log", `dataSlice string: `+dataSlice.Type().String()+` with len: `+strconv.Itoa(dataLen))
		return
	}

	if cap(v.dataBuf) < dataLen {
		v.dataBuf = make([]uint8, dataLen)
	} else if len(v.dataBuf) < dataLen {
		v.dataBuf = v.dataBuf[:dataLen]
	}
	js.CopyBytesToGo(v.dataBuf, dataSlice)
	v.dataBuf = v.dataBuf[:dataLen]

	for i4 := 0; i4+2 < len(v.dataBuf); i4 += 4 {
		grey := (v.dataBuf[i4+0] + v.dataBuf[i4+1] + v.dataBuf[i4+2]) / 3

		dataSlice.SetIndex(i4+0, grey) // as opposed to v.dataBuf[i4+0] = grey
		dataSlice.SetIndex(i4+1, grey)
		dataSlice.SetIndex(i4+2, grey)
		// v.dataBuf[i4+0] = grey
		// v.dataBuf[i4+1] = grey
		// v.dataBuf[i4+2] = grey
	}

	v.ctx2d.Call(`putImageData`, frame, 0, 0)
}
