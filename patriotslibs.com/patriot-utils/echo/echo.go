package echo

import (
	"fmt"
	"os"
)

type CursorANSI string

const (
	clearScreen      CursorANSI = "\033[2J"
	setAtStart       CursorANSI = "\033[H" // sets the cursor on the begining of the screen
	clearCurrentLine CursorANSI = "\033[K" // sets the cursor on the begining of the screen
)

type StylesANSI string

const (
	ItalicMode StylesANSI = "\033[3m"
)

type ColorANSI string

const (
	RESET     ColorANSI = "\x1b[0m"
	RedFG     ColorANSI = "\x1b[38;5;202m"
	RedBG     ColorANSI = "\x1b[48;5;202m"
	YellowFG  ColorANSI = "\x1b[38;5;220m"
	YellowBG  ColorANSI = "\x1b[48;5;220m"
	BlueFG    ColorANSI = "\x1b[38;5;27m"
	BlueBG    ColorANSI = "\x1b[48;5;27m"
	CyanFG    ColorANSI = "\x1b[38;5;123m"
	CyanBG    ColorANSI = "\x1b[48;5;123m"
	PinkFG    ColorANSI = "\x1b[38;5;206m"
	PinkBG    ColorANSI = "\x1b[48;5;206m"
	SkyBlueFG ColorANSI = "\x1b[38;5;81m"
	SkyBlueBG ColorANSI = "\x1b[48;5;81m"
	OrangeFG  ColorANSI = "\x1b[38;5;208m"
	OrangeBG  ColorANSI = "\x1b[48;5;208m"
	GreenFG   ColorANSI = "\x1b[38;5;190m"
	GreenBG   ColorANSI = "\x1b[48;5;190m"
	PurpleFG  ColorANSI = "\x1b[38;5;93m"
	PurpleBG  ColorANSI = "\x1b[48;5;93m"
	WhiteFG   ColorANSI = "\x1b[38;5;15m"
	WhiteBG   ColorANSI = "\x1b[48;5;15m"
)

func Echo(color ColorANSI, data ...interface{}) {
	fmt.Print(color)
	for h := 0; h < len(data); h++ {
		fmt.Print(data[h])
	}
	fmt.Println(RESET)
}

func EchoWarn(msg string) {
	fmt.Printf("%sWarning: %s%s\n", YellowFG, msg, RESET)
}

func EchoErr(err error) {
	fmt.Println("\033[5;37;41m ERROR "+string(RESET), string(RedFG)+err.Error())
}

func EchoDebug(msg string) {
	if os.Getenv("EDEBUG") != "" {
		fmt.Printf("%sDEBUG: %s%s\n", SkyBlueFG, msg, RESET)
	}
}
