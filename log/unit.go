package log

import "time"

type timeUnit string

func (t timeUnit) Format() string {
	switch t {
	case "minute":
		return ".%Y%m%d%H%M"
	case "hour":
		return ".%Y%m%d%H"
	case "day":
		return ".%Y%m%d"
	case "month":
		return ".%Y%m"
	case "year":
		return ".%Y"
	default:
		return ".%Y%m%d"
	}
}

func (t timeUnit) RotationGap() time.Duration {
	switch t {
	case "minute":
		return time.Minute
	case "hour":
		return time.Hour
	case "day":
		return time.Hour * 24
	case "month":
		return time.Hour * 24 * 30
	case "year":
		return time.Hour * 24 * 365
	default:
		return time.Hour * 24
	}
}
