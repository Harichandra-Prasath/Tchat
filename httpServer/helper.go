package httpserver

import "github.com/google/uuid"

func getChannel(id uuid.UUID) chan *message {
	chn, ok := ChnMapper.Load(id)
	if !ok {
		chn = make(chan *message, 100)
		ChnMapper.Store(id, chn)
	}
	return chn.(chan *message)
}
