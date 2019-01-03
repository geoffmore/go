package node

import (
	"fmt"
	"os/exec"
	"strings"
)

func stringInSlice(a string, list []string) bool {
	// https://stackoverflow.com/questions/15323767/does-go-have-if-x-in-construct-similar-to-python
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func runCommand(osCmdRaw string) []byte {
	// https://stackoverflow.com/questions/19238143/does-golang-support-variadic-function
	// https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
	// https://stackoverflow.com/questions/32721066/pass-string-to-a-function-that-expects-a-variadic-parameter

	osCmdSplit := strings.Split(osCmdRaw, " ")
	osCmd := osCmdSplit[0]
	osCmdArgs := osCmdSplit[1:]

	osCmdStruct := exec.Command(osCmd, osCmdArgs...)
	stdOutErr, _ := osCmdStruct.CombinedOutput()
	// Throwing away err here because STDERR will tell us what went wrong with the command
	return stdOutErr
}

func genCommand(prefix, base, suffix string) string {
	return fmt.Sprintf("%s %s %s", prefix, base, suffix)
}

func someFunc(origVal *string, present bool) {
	if !present {
		*origVal = "DNE"
	} else if *origVal == "" {
		*origVal = "BLANK"
	}
}
