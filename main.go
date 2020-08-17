package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"docker.io/go-docker"
	"github.com/c-bata/go-prompt"
)

var dockerClient *docker.Client

var portMappingSuggestions []prompt.Suggest

var suggestedImages []prompt.Suggest

var commandExpression = regexp.MustCompile(`(?P<command>exec|stop|start|run|service create|service inspect|service logs|service ls|service ps|service rollback|service scale|service update|service|pull|attach|build|commit|cp|create|events|export|history|images|import|info|inspect|kill|load|login|logs|ps|push|restart|rm|rmi|save|search|stack|stats|update|version)\s{1}`)

func getRegexGroups(text string) map[string]string {
	if !commandExpression.Match([]byte(text)) {
		return nil
	}

	match := commandExpression.FindStringSubmatch(text)
	result := make(map[string]string)
	for i, name := range commandExpression.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	return result
}

// func completer(d prompt.Document) []prompt.Suggest {
// 	word := d.GetWordAfterCursor()

// 	group := getRegexGroups(d.Text)
// 	if group != nil {
// 		command := group["command"]

// 		if command == "exec" || command == "stop" || command == "port" {

// 		}
// 	}
// }

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "users", Description: "Store the username and age"},
		{Text: "articles", Description: "Store the article text posted by user"},
		{Text: "comments", Description: "Store the text commented to articles"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func main() {
	dockerClient, _ = docker.NewEnvClient()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err := dockerClient.Ping(ctx); err != nil {
		fmt.Println("Couldn't check docker status please make sure docker is running.")
		fmt.Println(err)
		return
	}
	for {
		dockerCommand := prompt.Input(">>> docker ",
			completer,
			prompt.OptionTitle("docker prompt"),
			prompt.OptionSelectedDescriptionTextColor(prompt.Turquoise),
			prompt.OptionInputTextColor(prompt.Fuchsia),
			prompt.OptionPrefixBackgroundColor(prompt.Cyan))

		splittedDockerCommands := strings.Split(dockerCommand, " ")
		if splittedDockerCommands[0] == "exit" {
			os.Exit(0)
		}

		var ps *exec.Cmd

		if splittedDockerCommands[0] == "clear" {
			ps = exec.Command("clear")
		} else {
			ps = exec.Command("docker", splittedDockerCommands...)
		}

		res, err := ps.Output()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(res))

		portMappingSuggestions = []prompt.Suggest{}
		suggestedImages = []prompt.Suggest{}
	}
}
