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

	dataBuf      []uint8
	scaledData   [101][256]uint8
	numFrames    int
	maxNumFrames int
	rP, gP, bP   []int
}

func (v *videoProcessor) timerCB() {
	t0 := time.Now()
	v.computeFrame()

	dur := time.Since(t0)
	waitTime := (16 * time.Millisecond) - dur
	if waitTime < time.Millisecond {
		waitTime = time.Millisecond
	}
	go v.log(`computeFrame took: ` + strconv.Itoa(int(dur)/1000) + `us`)
	time.AfterFunc(waitTime, v.timerCB)
}

func (v *videoProcessor) updateWidthHeight() {
	v.width = v.video.Get(`width`).Int()
	v.height = v.video.Get(`height`).Int()
}

func (v *videoProcessor) load() {
	v.scaledData = [101][256]uint8{}
	for perc := 0; perc < 101; perc++ {
		for val := 0; val < 256; val++ {
			v.scaledData[perc][val] = uint8((float64(val) * float64(perc)) / 100.0)
		}
	}
	v.buildShifter()

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

func (v *videoProcessor) buildShifter() {
	// here's the goal. we're going to set percentages for each RGB
	// value in the following manner. As you move from left to right,
	// time is progressing
	// r: 100 - 100 -  0  -  0  - 100 - 100 -  0  - 100
	// g: 100 -  0  - 100 -  0  -  0  - 100 - 100 - 100
	// b: 100 -  0  -  0  - 100 - 100 -  0  - 100 - 100

	const minArrVal = 0
	const maxArrVal = 100

	vLong := 512
	v.rP = make([]int, vLong)
	v.gP = make([]int, vLong)
	v.bP = make([]int, vLong)
	fill := func(startI, numEntries, val int, arr []int) {
		if val < minArrVal {
			val = minArrVal
		} else if val > maxArrVal {
			val = maxArrVal
		}
		for i := 0; i < numEntries; i++ {
			arr[i+startI] = val
		}
	}
	dec := func(startI, numEntries, startVal, decVal int, arr []int) {
		arr[startI] = startVal
		var newVal int
		for i := 1; i < numEntries; i++ {
			arrI := startI + i
			newVal = arr[arrI-1] - decVal
			if newVal < minArrVal {
				newVal = minArrVal
			}
			arr[arrI] = newVal
		}
	}
	inc := func(startI, numEntries, startVal, incVal int, arr []int) {
		arr[startI] = startVal
		var newVal int
		for i := 1; i < numEntries; i++ {
			arrI := startI + i
			newVal = arr[arrI-1] + incVal
			if newVal > maxArrVal {
				newVal = maxArrVal
			}
			arr[arrI] = newVal
		}
	}

	numFrames := 50 // this is easier if we have a divisor of 100
	interval := maxArrVal / numFrames
	i := 0
	fill(i, numFrames, maxArrVal, v.rP)
	dec(i, numFrames, maxArrVal, interval, v.gP)
	dec(i, numFrames, maxArrVal, interval, v.bP)

	i += numFrames
	dec(i, numFrames, maxArrVal, interval, v.rP)
	inc(i, numFrames, minArrVal, interval, v.gP)
	fill(i, numFrames, minArrVal, v.bP)

	i += numFrames
	fill(i, numFrames, minArrVal, v.rP)
	dec(i, numFrames, maxArrVal, interval, v.gP)
	inc(i, numFrames, minArrVal, interval, v.bP)

	i += numFrames
	inc(i, numFrames, minArrVal, interval, v.rP)
	fill(i, numFrames, minArrVal, v.gP)
	fill(i, numFrames, maxArrVal, v.bP)

	i += numFrames
	fill(i, numFrames, maxArrVal, v.rP)
	inc(i, numFrames, minArrVal, interval, v.gP)
	dec(i, numFrames, maxArrVal, interval, v.bP)

	i += numFrames
	dec(i, numFrames, maxArrVal, interval, v.rP)
	fill(i, numFrames, maxArrVal, v.gP)
	inc(i, numFrames, minArrVal, interval, v.bP)

	i += numFrames
	inc(i, numFrames, minArrVal, interval, v.rP)
	fill(i, numFrames, maxArrVal, v.gP)
	fill(i, numFrames, maxArrVal, v.bP)

	v.maxNumFrames = i + numFrames
	if v.maxNumFrames > vLong {
		v.log(`should have panicked`)
	}
	v.rP = v.rP[:v.maxNumFrames]
	v.gP = v.gP[:v.maxNumFrames]
	v.bP = v.bP[:v.maxNumFrames]
}

func (v *videoProcessor) log(msg string) {
	v.console.Call("log", msg)
}

func (v *videoProcessor) computeFrame() {
	if v.width <= 0 || v.height <= 0 {
		v.updateWidthHeight()
		return
	}
	v.ctx2d.Call(`drawImage`, v.video, 0, 0, v.width, v.height)

	frame := v.ctx2d.Call(`getImageData`, 0, 0, v.width, v.height)

	var dataSlice js.Value
	switch colorViewStyle {
	case normalViewStyle:
		// do not populate the dataBuf
	default:
		dataSlice = frame.Get(`data`)

		dataLen := dataSlice.Get(`length`).Int()
		if dataLen <= 0 {
			v.log(`dataSlice string: ` + dataSlice.Type().String() + ` with len: ` + strconv.Itoa(dataLen))
			return
		}

		if cap(v.dataBuf) < dataLen {
			v.dataBuf = make([]uint8, dataLen)
		} else if len(v.dataBuf) < dataLen {
			v.dataBuf = v.dataBuf[:dataLen]
		}
		js.CopyBytesToGo(v.dataBuf, dataSlice)
		v.dataBuf = v.dataBuf[:dataLen]
	}

	switch colorViewStyle {
	case redOnlyViewStyle:
		for i4 := 0; i4+2 < len(v.dataBuf); i4 += 4 {
			// dataSlice.SetIndex(i4+0, v.dataBuf[i4])
			dataSlice.SetIndex(i4+1, 0)
			dataSlice.SetIndex(i4+2, 0)
		}
	case greenOnlyViewStyle:
		for i4 := 0; i4+2 < len(v.dataBuf); i4 += 4 {
			dataSlice.SetIndex(i4+0, 0)
			// dataSlice.SetIndex(i4+1, v.dataBuf[i4+1])
			dataSlice.SetIndex(i4+2, 0)
		}
	case blueOnlyViewStyle:
		for i4 := 0; i4+2 < len(v.dataBuf); i4 += 4 {
			dataSlice.SetIndex(i4+0, 0)
			dataSlice.SetIndex(i4+1, 0)
			// dataSlice.SetIndex(i4+2, v.dataBuf[i4+2])
		}

	case shiftingViewStyle:
		rPerc := v.rP[v.numFrames]
		gPerc := v.gP[v.numFrames]
		bPerc := v.bP[v.numFrames]

		var rVal, gVal, bVal uint8
		for i4 := 0; i4+2 < len(v.dataBuf); i4 += 4 {
			rVal = v.dataBuf[i4]
			dataSlice.SetIndex(i4, v.scaledData[rPerc][rVal])
			gVal = v.dataBuf[i4+1]
			dataSlice.SetIndex(i4+1, v.scaledData[gPerc][gVal])
			bVal = v.dataBuf[i4+1]
			dataSlice.SetIndex(i4+2, v.scaledData[bPerc][bVal])
		}
		v.numFrames++
		if v.numFrames == v.maxNumFrames {
			v.numFrames = 0
		}

	case normalViewStyle:
		// do nothing
	default:
		// oof developer error. return early
		return
	}

	v.ctx2d.Call(`putImageData`, frame, 0, 0)
}
