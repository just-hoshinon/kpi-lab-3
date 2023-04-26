package painter

import (
	"image/color"

	"golang.org/x/exp/shiny/screen"
)

// Operation змінює вхідну текстуру.
type Operation interface {
	// Do виконує зміну операції, повертаючи true, якщо текстура вважається готовою для відображення.
	Do(t screen.Texture) (ready bool)
}

// OperationList групує список операції в одну.
type OperationList []Operation

func (ol OperationList) Do(t screen.Texture) (ready bool) {
	for _, o := range ol {
		ready = o.Do(t) || ready
	}
	return
}

// UpdateOp операція, яка не змінює текстуру, але сигналізує, що текстуру потрібно розглядати як готову.
var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(t screen.Texture) bool { return true }

// OperationFunc використовується для перетворення функції оновлення текстури в Operation.
type OperationFunc func(t screen.Texture)

func (f OperationFunc) Do(t screen.Texture) bool {
	f(t)
	return false
}

// WhiteFill зафарбовує текстуру у білий колір. Може бути використана як Operation через OperationFunc(WhiteFill).
func WhiteFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.White, screen.Src)
}

// GreenFill зафарбовує текстуру у зелений колір. Може бути використана як Operation через OperationFunc(GreenFill).
func GreenFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.RGBA{G: 0xff, A: 0xff}, screen.Src)
}

type resetOp struct{}

// ResetOp операція очищує вікно
var ResetOp = resetOp{}

func (op resetOp) Do(t screen.Texture) bool { return true }

// BgRect операція додає чорний прямокутник на екран в певних координатах
type BgRect struct {
	X1 float32
	Y1 float32
	X2 float32
	Y2 float32
}

func (op BgRect) Do(t screen.Texture) bool {
	return false
}

// Figure операція додає фігуру варіанту на вказані координати
type Figure struct {
	X float32
	Y float32
}

func (op Figure) Do(t screen.Texture) bool {
	return false
}

// Move операція переміщує усі фігури у вказані координати
type Move struct {
	X float32
	Y float32
}

func (op Move) Do(t screen.Texture) bool {
	return false
}
