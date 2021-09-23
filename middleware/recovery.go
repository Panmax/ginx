package middleware

import "github.com/panmax/ginx"

func Recovery() ginx.ControllerHandler {
	return func(c *ginx.Context) error {
		defer func() {
			if err := recover(); err != nil {
				c.Json(500, err)
			}
		}()
		return c.Next()
	}
}
