package utils

import (
	"context"
	"time"
)

// Context - Context ve Timeout u döndürür
func Ctx(waitSecond time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), waitSecond*time.Second)
}
