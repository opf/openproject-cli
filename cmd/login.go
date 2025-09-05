package cmd

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/components/configuration"
	"github.com/opf/openproject-cli/components/parser"
	"github.com/opf/openproject-cli/components/paths"
	"github.com/opf/openproject-cli/components/printer"
	"github.com/opf/openproject-cli/components/requests"
	"github.com/opf/openproject-cli/components/resources/users"
	"github.com/opf/openproject-cli/dtos"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticates the user against an OpenProject instance",
	Long: `Enables the login flow, which enables the user to use
this tool for a specific OpenProject instance. The login
needs the host URL of the OpenProject instance and a
generated API token.`,
	Run: login,
}

const (
	urlInputError      = "There was a problem parsing the input. Please try again and put in a valid URL."
	missingSchemeError = "URL scheme is missing, please define a complete URL."
	noOpInstanceError  = "URL does not point to a valid OpenProject instance."
	tokenInputError    = "There was a problem parsing the token input. Please try again."
)

func login(_ *cobra.Command, _ []string) {
	var hostUrl *url.URL
	var token string

	for {
		printer.Debug(Verbose, "Parsing host URL ...")
		printer.Input("OpenProject host URL: ")

		ok, msg, host := parseHostUrl()
		if !ok {
			printer.ErrorText(msg)
			continue
		}

		printer.Debug(Verbose, "Initializing requests client ...")
		requests.Init(host, "", Verbose)
		ok = checkOpenProjectApi()
		if !ok {
			printer.ErrorText(noOpInstanceError)
			continue
		}

		hostUrl = host
		break
	}

	for {
		fmt.Printf("OpenProject API Token (Visit %s/my/access_token to generate one): ", hostUrl)
		ok, t := requestApiToken()
		if !ok {
			fmt.Println(tokenInputError)
			continue
		}

		token = common.SanitizeLineBreaks(t)

		requests.Init(hostUrl, token, Verbose)
		user, err := users.Me()
		if err != nil {
			printer.Error(err)
			continue
		}

		if user.Name == "Anonymous" {
			printer.ErrorText("no authenticate given")
			continue
		}

		break
	}

	storeLoginData(hostUrl, token)
}

func parseHostUrl() (ok bool, errMessage string, host *url.URL) {
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		printer.Debug(Verbose, fmt.Sprintf("Error reading string input: %+v", err))
		return false, urlInputError, nil
	}

	printer.Debug(Verbose, fmt.Sprintf("Parsed input %q.", input))
	printer.Debug(Verbose, "Sanitizing input ...")

	input = common.SanitizeLineBreaks(input)
	input = strings.TrimSuffix(input, "/")

	printer.Debug(Verbose, fmt.Sprintf("Sanitized input '%s'.", input))
	printer.Debug(Verbose, "Parsing input as url ...")

	parsed, err := url.Parse(input)
	if err != nil {
		printer.Debug(Verbose, fmt.Sprintf("Error parsing url: %+v", err))
		return false, urlInputError, nil
	}

	printer.Debug(Verbose, fmt.Sprintf("Parsed url '%s'.", parsed))
	printer.Debug(Verbose, "Checking for http host and scheme ...")

	if parsed.Scheme == "" || parsed.Host == "" {
		return false, missingSchemeError, nil
	}

	printer.Debug(Verbose, "Parsing input successful, continuing with next steps.")
	return true, "", parsed
}

func checkOpenProjectApi() bool {
	printer.Debug(Verbose, "Fetching API root to check for instance configuration ...")

	response, err := requests.Get(paths.Root(), nil)
	if err != nil {
		return false
	}

	c := parser.Parse[dtos.ConfigDto](response)

	return c.Type == "Root" && len(c.InstanceName) > 0
}

func requestApiToken() (ok bool, token string) {
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		return false, ""
	}

	return true, input
}

func storeLoginData(host *url.URL, token string) {
	err := configuration.WriteConfigFile(host.String(), token)
	if err != nil {
		printer.Error(err)
	}
}
