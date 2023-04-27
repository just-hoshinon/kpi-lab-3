package painter

import (
	"image"

	"golang.org/x/exp/shiny/screen"
)

// Receiver отримує текстуру, яка була підготовлена в результаті виконання команд у циелі подій.
type Receiver interface {
	Update(t screen.Texture)
}

// Loop реалізує цикл подій для формування текстури отриманої через виконання операцій отриманих з внутрішньої черги.
type Loop struct {
	Receiver Receiver

	next screen.Texture // Текстура, яка зараз формується
	prev screen.Texture // Текстура, яка була відправлена останнього разу у Receiver

	mq    MessageQueue
	state TextureState
}

var size = image.Pt(400, 400)

// Start запускає цикл подій. Цей метод потрібно запустити до того, як викликати на ньому будь-які інші методи.
func (l *Loop) Start(s screen.Screen) {
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)
	l.mq = MessageQueue{queue: make(chan Operation, 15)}
	l.state = TextureState{}

	go func() {
		for {
			e := l.mq.Pull()

			switch e.(type) {
			case Figure, BgRect, Move, Fill, Reset:
				e.Update(l.state)
			case Update:
				t, _ := s.NewTexture(size)
				l.state.backgroundColor.Do(t)
				l.state.backgroundRect.Do(t)
				for _, fig := range l.state.figureCenters {
					fig.Do(t)
				}
				l.Receiver.Update(t)
			}
		}
	}()
}

// Post додає нову операцію у внутрішню чергу.
func (l *Loop) Post(ol OperationList) {

	for _, op := range ol {
		l.mq.Push(op)
	}
}

// StopAndWait сигналізує
func (l *Loop) StopAndWait() {

}

// MessageQueue черга повідомлень
type MessageQueue struct {
	queue chan Operation
}

func (mq *MessageQueue) Push(op Operation) {
	mq.queue <- op
}

func (mq *MessageQueue) Pull() Operation {
	return <-mq.queue
}
