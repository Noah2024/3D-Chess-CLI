package color

//"\033[31mThis is red text\033[0m"

type COLOR string

const (
	//Foreground
	CLR    = "\033[0m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Green  = "\033[32m"
	Blue   = "\033[34m"
	Purple = "\033[35m"

	// Background
	BgBlack  = "\033[40m"
	BgRed    = "\033[41m"
	BgYellow = "\033[43m"
	BgGreen  = "\033[42m"
	BgBlue   = "\033[44m"
	BgPurple = "\033[45m"
	BgWhite  = "\033[47m"
)

func ColorText(str string, color COLOR) string {
	return string(color) + str + CLR
}
