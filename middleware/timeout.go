package middleware

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/panmax/ginx"
)

func Timeout(d time.Duration) ginx.ControllerHandler {
	return func(c *ginx.Context) error {
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		durationCtx, cancel := context.WithTimeout(c.BaseContext(), d)
		defer cancel()

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			c.Next()
			finish <- struct{}{}
		}()

		select {
		case p := <-panicChan:
			log.Println(p)
			c.Json(500, "internal error")
		case <-finish:
			fmt.Println("finish")
		case <-durationCtx.Done():
			c.SetHasTimeout()
			c.Json(500, "timeout")
		}
		return nil
	}
}
