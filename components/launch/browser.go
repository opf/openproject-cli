package launch

import (
	"github.com/opf/openproject-cli/components/printer"
	"net/url"
	"os/exec"
	"runtime"
)

func Browser(url url.URL) error {
	var command *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		command = exec.Command("cmd", "/c", "start", url.String())
	case "darwin":
		command = exec.Command("open", url.String())
	case "linux":
		command = exec.Command("xdg-open", url.String())
	default:
		printer.ErrorText("operating system not supported")
	}

	return command.Start()
}
