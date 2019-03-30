package utils

import "time"

func UnixtimeMilli() int64 {
	return time.Now().UnixNano() / 1000000
}
func UnixtimeSec() int64 {
	return time.Now().Unix()
}
