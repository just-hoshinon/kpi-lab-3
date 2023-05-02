package lang

import (
	"errors"
	"github.com/roman-mazur/architecture-lab-3/painter"
	"github.com/stretchr/testify/assert"
	"image/color"
	"io"
	"strings"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	parser := Parser{}

	var whiteCmd io.Reader = strings.NewReader("white")
	whiteRes, whiteErr := parser.Parse(whiteCmd)

	if assert.Nil(t, whiteErr) {
		assert.Equal(t, painter.Fill{Color: color.White}, whiteRes[0])
	}

	var greenCmd io.Reader = strings.NewReader("green")
	greenRes, greenErr := parser.Parse(greenCmd)

	if assert.Nil(t, greenErr) {
		assert.Equal(t, painter.Fill{Color: color.RGBA{G: 0xff, A: 0xff}}, greenRes[0])
	}

	var updateCmd io.Reader = strings.NewReader("update")
	updateRes, updateErr := parser.Parse(updateCmd)

	if assert.Nil(t, updateErr) {
		assert.Equal(t, painter.UpdateOp, updateRes[0])
	}

	var bgRectCmd io.Reader = strings.NewReader("bgrect 0.25 0.25 0.75 0.75")
	bgRectRes, bgRectErr := parser.Parse(bgRectCmd)

	if assert.Nil(t, bgRectErr) {
		assert.Equal(t, painter.BgRect{X1: 0.25, Y1: 0.25, X2: 0.75, Y2: 0.75}, bgRectRes[0])
	}

	var figureCmd io.Reader = strings.NewReader("figure 0.5 0.5")
	figureRes, figureErr := parser.Parse(figureCmd)

	if assert.Nil(t, figureErr) {
		assert.Equal(t, painter.Figure{X: 0.5, Y: 0.5}, figureRes[0])
	}

	var moveCmd io.Reader = strings.NewReader("move 0.2 0.2")
	moveRes, moveErr := parser.Parse(moveCmd)

	if assert.Nil(t, moveErr) {
		assert.Equal(t, painter.Move{X: 0.2, Y: 0.2}, moveRes[0])
	}

	var resetCmd io.Reader = strings.NewReader("reset")
	resetRes, resetErr := parser.Parse(resetCmd)

	if assert.Nil(t, resetErr) {
		assert.Equal(t, painter.ResetOp, resetRes[0])
	}

	var wrongCmd io.Reader = strings.NewReader("some wrong command")
	_, wrongErr := parser.Parse(wrongCmd)

	if assert.NotNil(t, wrongErr) {
		assert.Equal(t, errors.New("unknown command"), wrongErr)
	}
}
