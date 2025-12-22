package api_key

import (
	"mjrc/core/chix"
)

const Name = "api_key"

func Middleware(apiKey string) *chix.Middleware {
	hdlr := &handler{apiKey}
	return chix.NewMiddleware(Name, hdlr.validateApiKey)
}
