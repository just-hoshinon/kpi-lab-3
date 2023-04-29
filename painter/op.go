package painter

import (
	"github.com/roman-mazur/architecture-lab-3/ui"
	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"
)

// Operation описує всі можливі операції
type Operation interface {
	Update(state *TextureState)
}

// TextureOperation змінюють текстуру
type TextureOperation interface {
	// Do виконує зміну операції, повертаючи true, якщо текстура вважається готовою для відображення.
	Do(t screen.Texture)
	Update(state *TextureState)
}

// OperationList групує список операції в одну.
type OperationList []Operation

// UpdateOp операція, яка не змінює текстуру, але сигналізує, що текстуру потрібно розглядати як готову.
var UpdateOp = Update{}

type Update struct{}

func (op Update) Update(state *TextureState) {}

// Fill зафарбовує текстуру у відповідний колір
type Fill struct {
	Color color.Color
}

func (op Fill) Do(t screen.Texture) {
	t.Fill(t.Bounds(), op.Color, screen.Src)
}

func (op Fill) Update(state *TextureState) {
	state.backgroundColor = op
}

type Reset struct{}

// ResetOp операція очищує вікно
var ResetOp = Reset{}

func (op Reset) Update(state *TextureState) {
	state.backgroundColor = Fill{Color: color.Black}
	state.backgroundRect = nil
	state.figureCenters = nil
}

// BgRect операція додає чорний прямокутник на екран в певних координатах
type BgRect struct {
	X1 float32
	Y1 float32
	X2 float32
	Y2 float32
}

func (op BgRect) Do(t screen.Texture) {
	t.Fill(
		image.Rect(
			int(op.X1*float32(t.Size().X)),
			int(op.Y1*float32(t.Size().Y)),
			int(op.X2*float32(t.Size().X)),
			int(op.Y2*float32(t.Size().Y)),
		),
		color.Black,
		screen.Src,
	)
}

func (op BgRect) Update(state *TextureState) {
	state.backgroundRect = op
}

// Figure операція додає фігуру варіанту на вказані координати
type Figure struct {
	X float32
	Y float32
}

func (op Figure) Do(t screen.Texture) {
	ui.DrawT(
		t,
		image.Pt(
			int(op.X*float32(t.Size().X)),
			int(op.Y*float32(t.Size().Y)),
		),
	)
}

func (op Figure) Update(state *TextureState) {
	state.figureCenters = append(state.figureCenters, &op)
}

// Move операція переміщує усі на відповідну кількість пікселів
type Move struct {
	X float32
	Y float32
}

func (op Move) Update(state *TextureState) {
	for _, fig := range state.figureCenters {
		fig.X = op.X
		fig.Y = op.Y
	}
}
