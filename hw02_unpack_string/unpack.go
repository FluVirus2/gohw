package hw02unpackstring

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	builder := strings.Builder{}

	reader := strings.NewReader(s)

	// attempt to read the first rune
	firstRune, _, rerr := reader.ReadRune()

	// if there is no runes, EOF err will be set
	if rerr == io.EOF {
		return "", nil
	}

	// if first rune is digit it cannot be legal string
	if unicode.IsDigit(firstRune) {
		return "", fmt.Errorf("invalid character spotted: %w", ErrInvalidString)
	}

	// holds previous rune by pointer
	prevRune := &firstRune

	// while all runes are not read
	for r, _, err := reader.ReadRune(); err == nil; r, _, err = reader.ReadRune() {
		isDigit := unicode.IsDigit(r)

		// if a letter was read and no letter is being held
		if !isDigit && prevRune == nil {
			rn := r
			prevRune = &rn

			continue
		}

		// if a letter was read and any letter is being held
		if !isDigit {
			_, werr := builder.WriteRune(*prevRune)
			if werr != nil {
				return "", fmt.Errorf("error while writing in strings.Builder: %w", werr)
			}

			// hold new letter
			rn := r
			prevRune = &rn
		}

		// if a digit was read and no letter is being held is mistake
		if isDigit && prevRune == nil {
			return "", fmt.Errorf("invalid digit position spotted: %w", ErrInvalidString)
		}

		// if a digit was read and any letter is being held
		if isDigit {
			// r times write rune
			for range r - '0' {
				_, werr := builder.WriteRune(*prevRune)
				if werr != nil {
					return "", fmt.Errorf("error while writing in strings.Builder: %w", werr)
				}
			}

			// flush prev rune (set it to nil)
			prevRune = nil
		}
	}

	// if there is a rune not flushed yet
	if prevRune != nil {
		_, werr := builder.WriteRune(*prevRune)
		if werr != nil {
			return "", fmt.Errorf("error while writing in srings.Builder: %w", werr)
		}
	}

	return builder.String(), nil
}
