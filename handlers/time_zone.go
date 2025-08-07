package handlers

import "time"

func ToWIB(t time.Time) string {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	return t.In(loc).Format("2006-01-02 15:04:05")
}
