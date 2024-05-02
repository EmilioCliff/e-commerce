package api

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	token "github.com/EmilioCliff/e-commerce/db/token"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	AuthorizationHeader = "Authorization"
	AuthorizationType   = "Bearer"
	PayloadKey          = "payload"
)

func authMiddleware(maker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader(AuthorizationHeader)
		if header == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("no authorization header passed")))
			return
		}

		args := strings.Fields(header)
		if len(args) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("missing the bearer or token arguments")))
			return
		}

		if args[0] != AuthorizationType {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("authorization type not supported")))
			return
		}

		payload, footer, err := maker.VerifyToken(args[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("invalid access token")))
			return
		}
		if footer != token.Footer {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("does not recognize the token footer")))
			return
		}

		ctx.Set(PayloadKey, payload)
		// ctx.MustGet(PayloadKey)
		ctx.Next()
	}
}

func loggerMiddleware() gin.HandlerFunc {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		duration := time.Since(start)

		var errors []error
		for _, err := range ctx.Errors {
			errors = append(errors, err)
		}
		logger := log.Info()
		if len(ctx.Errors) > 0 {
			logger = log.Error().Errs("errors", errors)
		}

		logger.
			Str("method", ctx.Request.Method).
			Str("path", ctx.Request.RequestURI).
			Int("status_code", ctx.Writer.Status()).
			Str("status_text", http.StatusText(ctx.Writer.Status())).
			Dur("duration", duration)
	}
}
