package command

import "github.com/abzicht/gpsplit/gpxtransform/config"

type TransformCommand interface {
	GetConfiguration() (config.TransformConfig, error)
}

type CommandError struct {
	Msg string
}

func (ce CommandError) Error() string {
	return ce.Msg
}
