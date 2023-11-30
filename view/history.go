package view

import "github.com/rs/zerolog"

type History struct {
	stack       []func()
	HistorySize int
	Logger      *zerolog.Logger
}

func (h *History) push(back func()) {
	h.stack = append(h.stack, back)
	if len(h.stack) > h.HistorySize {
		h.stack = h.stack[1:]
	}
	h.Logger.Debug().Msgf("History stack: %v", h.stack)

}

func (h *History) pop() {
	h.Logger.Debug().Msgf("History stack: %v", h.stack)
	if len(h.stack) > 1 {
		last := h.stack[len(h.stack)-2]
		last()
		h.stack = h.stack[:len(h.stack)-2]
	}
}
