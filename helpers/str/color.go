package str

import "fmt"

func Red(str interface{}) string {
	return fmt.Sprintf("\n\033[36;31m%v\033[0m", str)
}

func Redf(format string, args ...interface{}) string {
	return Red(fmt.Sprintf(format, args...))
}

func Yellow(str interface{}) string {
	return fmt.Sprintf("\n\033[33m%v\033[0m", str)
}

func Yellowf(format string, args ...interface{}) string {
	return Yellow(fmt.Sprintf(format, args...))
}

func Magenta(str interface{}) string {
	return fmt.Sprintf("\n\033[35m%v\033[0m", str)
}

func Magentaf(format string, args ...interface{}) string {
	return Magenta(fmt.Sprintf(format, args...))
}
