package util

import "time"

func ToPtr[T any](x T) *T {
	return &x
}

func Today() string {
	return time.Now().Format("2006-01-02")
}
