package color

//"\033[31mThis is red text\033[0m"

type COLOR string

const (
	CLR    = "\033[0m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Green  = "\033[32m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
)

func ColorText(str string, color COLOR) string {
	return string(color) + str + CLR
}
