package loghelper

import (
	"log"

	strhelper "github.com/rayyone/go-core/helpers/str"
)

func PrintRed(str string) {
	log.Println(strhelper.Red(str))
}

func PrintRedf(format string, args ...interface{}) {
	log.Println(strhelper.Redf(format, args...))
}

func PrintYellow(str string) {
	log.Println(strhelper.Yellow(str))
}

func PrintYellowf(format string, args ...interface{}) {
	log.Println(strhelper.Yellowf(format, args...))
}

func PrintMagenta(str string) {
	log.Println(strhelper.Magenta(str))
}

func PrintMagentaf(format string, args ...interface{}) {
	log.Println(strhelper.Magentaf(format, args...))
}
