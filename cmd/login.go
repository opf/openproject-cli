package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/opf/openproject-cli/components/common"
	"github.com/opf/openproject-cli/components/configuration"
	openerrors "github.com/opf/openproject-cli/components/errors"
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
	RunE: login,
}

const (
	urlInputError      = "There was a problem parsing the input. Please try again and put in a valid URL."
	missingSchemeError = "URL scheme is missing, please define a complete URL."
	noOpInstanceError  = "URL does not point to a valid OpenProject instance."
	tokenInputError    = "There was a problem parsing the token input. Please try again."
)

func login(cmd *cobra.Command, _ []string) error {
	profile, err := resolveLoginProfile(cmd)
	if err != nil {
		printer.Error(err)
		return openerrors.ErrHandled
	}

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
		printer.Input(fmt.Sprintf("OpenProject API Token (Visit %s/my/access_tokens to generate one): ", hostUrl))
		ok, t := requestApiToken()
		if !ok {
			printer.ErrorText(tokenInputError)
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

	return storeLoginData(profile, hostUrl, token)
}

// resolveLoginProfile determines the profile name for the login command.
//
//   - If --profile was explicitly passed: validate the value immediately,
//     display it, and use it without prompting.
//   - If OP_CLI_PROFILE is set (but --profile was not): display the value
//     and use it without prompting.
//   - Otherwise: prompt the user interactively.
func resolveLoginProfile(cmd *cobra.Command) (string, error) {
	if cmd.Root().PersistentFlags().Changed("profile") {
		if err := configuration.ValidateProfileName(profileName); err != nil {
			return "", err
		}
		printer.Info(fmt.Sprintf("Profile: %s", profileName))
		return profileName, nil
	}

	if env := os.Getenv(configuration.EnvProfile); env != "" {
		if err := configuration.ValidateProfileName(env); err != nil {
			return "", err
		}
		printer.Info(fmt.Sprintf("Profile: %s", env))
		return env, nil
	}

	return promptProfileName(configuration.DefaultProfile)
}

// promptProfileName shows an interactive prompt, re-prompting on invalid input
// and offering the sanitized form as the next default.
func promptProfileName(defaultName string) (string, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		printer.Input(fmt.Sprintf("Profile name? [%s] ", defaultName))
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		input = common.SanitizeLineBreaks(strings.TrimSpace(input))
		if input == "" {
			return defaultName, nil
		}

		if err := configuration.ValidateProfileName(input); err != nil {
			sanitized := configuration.SanitizeProfileName(input)
			printer.ErrorText(
				"Invalid profile name. Only letters, numbers, - and _ are allowed (no leading/trailing hyphens).",
			)
			defaultName = sanitized
			continue
		}

		return input, nil
	}
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

	statusCode, header, body, err := requests.Probe(paths.Root())
	if err != nil {
		printer.Debug(Verbose, fmt.Sprintf("Error probing OpenProject API: %+v", err))
		return false
	}

	// Public instance: standard check on the root resource
	if statusCode == http.StatusOK {
		c := parser.Parse[dtos.ConfigDto](body)
		return c.Type == "Root" && len(c.InstanceName) > 0
	}

	// Auth-required instance: detect OpenProject via the Link header added before authentication
	if statusCode == http.StatusUnauthorized {
		linkHeader := header.Get("Link")
		if strings.Contains(linkHeader, "/api/v3/openapi.json") {
			return true
		}
		// Fallback: check error body for OpenProject-specific error identifier
		return strings.Contains(string(body), "openproject-org")
	}

	return false
}

func requestApiToken() (ok bool, token string) {
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')
	if err != nil {
		return false, ""
	}

	return true, input
}

func confirmOverwrite(profile string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	printer.Input(fmt.Sprintf("Profile %q already exists, overwrite? [y/N] ", profile))
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	answer := strings.ToLower(strings.TrimSpace(common.SanitizeLineBreaks(input)))
	return answer == "y" || answer == "yes", nil
}

func storeLoginData(profile string, host *url.URL, token string) error {
	profiles, err := configuration.AllProfiles()
	if err != nil {
		printer.Error(err)
		return openerrors.ErrHandled
	}

	for _, p := range profiles {
		if p.Name == profile {
			ok, err := confirmOverwrite(profile)
			if err != nil {
				printer.Error(err)
				return openerrors.ErrHandled
			}
			if !ok {
				printer.Info("Login cancelled.")
				return nil
			}
			break
		}
	}

	if err := configuration.WriteConfigForProfile(profile, host.String(), token); err != nil {
		printer.Error(err)
		return openerrors.ErrHandled
	}
	printer.Done()
	return nil
}
