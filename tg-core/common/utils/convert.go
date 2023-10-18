package utils

import "strconv"

func Str2Int(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func Str2IntMust(s string, defaultValue int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return i
}

func Str2Int32(s string) (int32, error) {
	i64, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(i64), nil
}

func Str2Int32Must(s string, defaultValue int32) int32 {
	i64, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return defaultValue
	}
	return int32(i64)
}

func Str2Int64(s string) (int64, error) {
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return i64, nil
}

func Str2Int64Must(s string, defaultValue int64) int64 {
	i64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return defaultValue
	}
	return i64
}

func Str2Float32(s string) (float32, error) {
	f32, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}
	return float32(f32), nil
}

func Str2Float32Must(s string, defaultValue float32) float32 {
	f32, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return defaultValue
	}
	return float32(f32)
}

func Str2Float64(s string) (float64, error) {
	f64, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return float64(f64), nil
}

func Str2Float64Must(s string, defaultValue float64) float64 {
	f64, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultValue
	}
	return float64(f64)
}