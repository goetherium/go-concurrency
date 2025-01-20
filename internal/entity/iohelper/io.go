package iohelper

import (
	"io"
)

func Read(r io.Reader, b []byte) ([]byte, error) {
	for {
		rest := cap(b) - len(b)

		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			return b, err
		}

		if n < rest {
			// Прочитали всё что клиент прислал, выходим чтобы
			// не заблокироваться на следующем чтении.
			return b, nil
		}

		if len(b) == cap(b) {
			// Add more capacity (let append pick how much).
			b = append(b, 0)[:len(b)]
		}
	}
}
