package gox

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Gox struct {
	*echo.Echo

	DB     *gorm.DB
	models []any

	goxLogger *logrus.Logger
}

func New() *Gox {
	gox := &Gox{
		Echo:   echo.New(),
		models: []any{},
	}

	gox.HideBanner = true

	gox.goxLogger = logrus.New()
	gox.goxLogger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	gox.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}] ${protocol} ${status} ${method} ${uri} ${latency_human} ${error}\n",
	}))

	return gox
}

type GoxHandlerFunc func(echo.Context, *Gox) error

func (g *Gox) WithGox(handler GoxHandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handler(c, g)
	}
}
