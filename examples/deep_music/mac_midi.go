// note: go mod tidy gets confused about this somehow -- so excluding when doing it.

//go:build not_darwin

package main

import (
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)
