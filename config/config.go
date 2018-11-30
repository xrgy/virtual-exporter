package config

import "strconv"

func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

func ValToGB(str string, unit string) (float64,error) {
	val,err:= strconv.ParseFloat(str,64)
	if err!=nil{
		return 0,err
	}
	if unit=="TB" {
		return val*1024,nil
	}
	if unit=="MB" {
		return val/1024,nil
	}
	return val,nil

}