package login

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

func getTerminalSize() (int, int, error) {
	cmd := exec.Command("mode", "con")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ls: error reading console width %s", err.Error())
	}
	re := regexp.MustCompile(`\d+`)
	rs := re.FindAllString(out.String(), -1)
	i, err := strconv.Atoi(rs[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ls: error transfering string to int %s", err.Error())
	}
	ch <- i
}
