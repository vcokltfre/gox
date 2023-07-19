package gox

import "github.com/labstack/echo/v4"

type Gox struct {
	*echo.Echo
}

func New() *Gox {
	gox := &Gox{
		Echo: echo.New(),
	}

	return gox
}
