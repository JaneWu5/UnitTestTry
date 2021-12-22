package main

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"os"
)

func C() error {
	_, err := os.Open("err")
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func B() error {
	return C()
}
func A() error {
	return B()
}
func init() {
	log.Logger = zerolog.New(os.Stderr).
		With().Timestamp().Stack().CallerWithSkipFrameCount(2).Logger()
	// 使用官方提供的，输出更友好
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

func main() {
	err := A()
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
	}
}
