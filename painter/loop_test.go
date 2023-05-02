package painter

import (
	"errors"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"image"
	"image/color"
	"testing"
)

type MockReceiver struct {
	calls int
}

func (rec *MockReceiver) Update(_ screen.Texture) {
	rec.calls++
}

type MockScreen struct{}

func (s MockScreen) NewBuffer(_ image.Point) (screen.Buffer, error) {
	return nil, errors.New("nothing")
}
func (s MockScreen) NewTexture(_ image.Point) (screen.Texture, error) {
	return nil, errors.New("nothing")
}
func (s MockScreen) NewWindow(_ *screen.NewWindowOptions) (screen.Window, error) {
	return nil, errors.New("nothing")
}

type ConChecker struct {
	length int
	res    chan struct{}
}

func (c ConChecker) done() {
	c.res <- struct{}{}
}

func (c ConChecker) check() {
	for i := 0; i < c.length; i++ {
		<-c.res
	}
}

func makeChecker(length int) ConChecker {
	return ConChecker{length: length, res: make(chan struct{}, length)}
}

func TestBlackFill(t *testing.T) {
	ops := OperationList{
		Fill{Color: color.RGBA{G: 0xff, A: 0xff}},
		Fill{Color: color.Black},
	}

	c := makeChecker(len(ops))
	loop := Loop{Receiver: &MockReceiver{}, doneFunc: c.done}

	loop.Start(MockScreen{})
	loop.Post(ops)

	c.check()

	if loop.state.backgroundColor.Color != color.Black {
		t.Error("Incorrect color")
	}
}

func TestDefault(t *testing.T) {
	ops := OperationList{}

	loop := Loop{Receiver: &MockReceiver{}}

	loop.Start(MockScreen{})
	loop.Post(ops)

	if loop.state.backgroundColor.Color != color.White || loop.state.backgroundRect != nil || loop.state.figureCenters != nil {
		t.Error("Incorrect color")
	}
}

func TestManyFills(t *testing.T) {
	ops := OperationList{
		Fill{Color: color.RGBA{G: 0x4f, A: 0xfb}}, Fill{Color: color.Black},
		Fill{Color: color.RGBA{G: 0xff, A: 0xff}}, Fill{Color: color.Gray{Y: 40}},
		Fill{Color: color.RGBA{G: 0xff, A: 0xfa}}, Fill{Color: color.White},
		Fill{Color: color.RGBA{G: 0xfc, A: 0x4f}}, Fill{Color: color.Black},
	}

	c := makeChecker(len(ops))
	loop := Loop{Receiver: &MockReceiver{}, doneFunc: c.done}

	loop.Start(MockScreen{})
	loop.Post(ops)

	c.check()

	if loop.state.backgroundColor.Color != color.Black {
		t.Error("Incorrect color")
	}
}

func TestSaveLastRect(t *testing.T) {
	ops := OperationList{
		BgRect{
			X1: 0.4,
			Y1: 0.3,
			X2: 0.5,
			Y2: 0.7,
		},
		BgRect{
			X1: 0.1,
			Y1: 0.1,
			X2: 0.2,
			Y2: 0.3,
		},
	}

	c := makeChecker(len(ops))
	loop := Loop{Receiver: &MockReceiver{}, doneFunc: c.done}
	loop.Start(MockScreen{})
	loop.Post(ops)
	c.check()

	last := BgRect{
		X1: 0.1,
		Y1: 0.1,
		X2: 0.2,
		Y2: 0.3,
	}

	if *loop.state.backgroundRect != last {
		t.Error("Incorrect rect")
	}
}

func TestAddFigures(t *testing.T) {
	ops := OperationList{
		Figure{
			X: 0.4,
			Y: 0.6,
		},
		Figure{
			X: 0.1,
			Y: 0.2,
		},
	}

	c := makeChecker(len(ops))
	loop := Loop{Receiver: &MockReceiver{}, doneFunc: c.done}
	loop.Start(MockScreen{})
	loop.Post(ops)
	c.check()

	first := Figure{
		X: 0.4,
		Y: 0.6,
	}
	second := Figure{
		X: 0.1,
		Y: 0.2,
	}

	if *loop.state.figureCenters[0] != first || *loop.state.figureCenters[1] != second {
		t.Error("Incorrect figures")
	}
}

func TestMoveBothFigures(t *testing.T) {
	ops := OperationList{
		Figure{
			X: 0.4,
			Y: 0.6,
		},
		Figure{
			X: 0.1,
			Y: 0.2,
		},
		Move{
			X: 0.3,
			Y: 0.1,
		},
	}

	c := makeChecker(len(ops))
	loop := Loop{Receiver: &MockReceiver{}, doneFunc: c.done}
	loop.Start(MockScreen{})
	loop.Post(ops)
	c.check()

	moved := Figure{
		X: 0.3,
		Y: 0.1,
	}

	if *loop.state.figureCenters[0] != moved || *loop.state.figureCenters[1] != moved {
		t.Error("Incorrect figures")
	}
}

func TestMoveFirstFigure(t *testing.T) {
	ops := OperationList{
		Figure{
			X: 0.4,
			Y: 0.6,
		},
		Move{
			X: 0.3,
			Y: 0.1,
		},
		Figure{
			X: 0.1,
			Y: 0.2,
		},
	}

	c := makeChecker(len(ops))
	loop := Loop{Receiver: &MockReceiver{}, doneFunc: c.done}
	loop.Start(MockScreen{})
	loop.Post(ops)
	c.check()

	first := Figure{
		X: 0.3,
		Y: 0.1,
	}
	second := Figure{
		X: 0.1,
		Y: 0.2,
	}

	if *loop.state.figureCenters[0] != first || *loop.state.figureCenters[1] != second {
		t.Error("Incorrect figures")
	}
}

func TestDontMoveFigures(t *testing.T) {
	ops := OperationList{
		Move{
			X: 0.3,
			Y: 0.1,
		},
		Figure{
			X: 0.4,
			Y: 0.6,
		},
		Figure{
			X: 0.1,
			Y: 0.2,
		},
	}

	c := makeChecker(len(ops))
	loop := Loop{Receiver: &MockReceiver{}, doneFunc: c.done}
	loop.Start(MockScreen{})
	loop.Post(ops)
	c.check()

	first := Figure{
		X: 0.4,
		Y: 0.6,
	}
	second := Figure{
		X: 0.1,
		Y: 0.2,
	}

	if *loop.state.figureCenters[0] != first || *loop.state.figureCenters[1] != second {
		t.Error("Incorrect figures")
	}
}

func TestReset(t *testing.T) {
	ops := OperationList{
		Figure{
			X: 0.4,
			Y: 0.6,
		},
		Figure{
			X: 0.1,
			Y: 0.2,
		},
		BgRect{
			X1: 0.4,
			Y1: 0.3,
			X2: 0.5,
			Y2: 0.7,
		},
		Move{
			X: 0.3,
			Y: 0.1,
		},
		Fill{Color: color.RGBA{G: 0xff, A: 0xff}},
		BgRect{
			X1: 0.1,
			Y1: 0.1,
			X2: 0.2,
			Y2: 0.3,
		},
		Fill{Color: color.Black},
		Reset{},
	}

	c := makeChecker(len(ops))
	loop := Loop{Receiver: &MockReceiver{}, doneFunc: c.done}
	loop.Start(MockScreen{})
	loop.Post(ops)
	c.check()

	if loop.state.figureCenters != nil || loop.state.backgroundRect != nil || loop.state.backgroundColor.Color != color.Black {
		t.Error("Reset works incorrectly")
	}
}

func TestInChaoticOrder(t *testing.T) {
	ops := OperationList{
		Figure{
			X: 0.4,
			Y: 0.6,
		},
		BgRect{
			X1: 0.4,
			Y1: 0.3,
			X2: 0.5,
			Y2: 0.7,
		},
		Move{
			X: 0.3,
			Y: 0.1,
		},
		Fill{Color: color.Black},
		BgRect{
			X1: 0.1,
			Y1: 0.1,
			X2: 0.2,
			Y2: 0.3,
		},
		Figure{
			X: 0.1,
			Y: 0.2,
		},
		Fill{Color: color.RGBA{G: 0xff, A: 0xff}},
	}

	c := makeChecker(len(ops))
	loop := Loop{Receiver: &MockReceiver{}, doneFunc: c.done}
	loop.Start(MockScreen{})
	loop.Post(ops)
	c.check()

	figure1 := Figure{
		X: 0.3,
		Y: 0.1,
	}
	figure2 := Figure{
		X: 0.1,
		Y: 0.2,
	}
	rect := BgRect{
		X1: 0.1,
		Y1: 0.1,
		X2: 0.2,
		Y2: 0.3,
	}
	fill := Fill{
		Color: color.RGBA{G: 0xff, A: 0xff},
	}

	if *loop.state.figureCenters[0] != figure1 || *loop.state.figureCenters[1] != figure2 || *loop.state.backgroundRect != rect || *loop.state.backgroundColor != fill {
		t.Error("Chaotic order works incorrectly")
	}
}

func TestUpdate(t *testing.T) {
	var src screen.Screen
	done := make(chan struct{})
	go driver.Main(func(s screen.Screen) {
		src = s
		done <- struct{}{}
		select {}
	})
	<-done

	ops := OperationList{
		Figure{
			X: 0.4,
			Y: 0.6,
		},
		BgRect{
			X1: 0.4,
			Y1: 0.3,
			X2: 0.5,
			Y2: 0.7,
		},
		Move{
			X: 0.3,
			Y: 0.1,
		},
		Fill{Color: color.Black},
		BgRect{
			X1: 0.1,
			Y1: 0.1,
			X2: 0.2,
			Y2: 0.3,
		},
		Figure{
			X: 0.1,
			Y: 0.2,
		},
		Fill{Color: color.RGBA{G: 0xff, A: 0xff}},
		Update{},
	}

	rec := &MockReceiver{}

	c := makeChecker(len(ops))
	loop := Loop{Receiver: rec, doneFunc: c.done}
	loop.Start(src)
	loop.Post(ops)
	c.check()

	if rec.calls != 1 {
		t.Error("Update works incorrectly")
	}
}
