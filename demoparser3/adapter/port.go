package adapter

import (
	"github.com/markus-wa/demoinfocs-golang/v5/pkg/demoinfocs"
	dp "github.com/markus-wa/godispatch"
)

// demoParser abstraction to allow for dependency injection
type demoParser interface {
	ParseToEnd() error
	RegisterNetMessageHandler(handler any) dp.HandlerIdentifier
	RegisterEventHandler(handler any) dp.HandlerIdentifier
	GameState() demoinfocs.GameState
}
