/*
   Copyright [2018] [Samuel Stevens]

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package main

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
)

type host struct {
	Name     string
	HostName string
	User     string
}

func readFile(file string) string {
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("%q does not exist", file)
	}
	return fmt.Sprintf("%s", contents)
}

func readHosts(lines []string) []host {
	hosts := make([]host, 0)
	current := host{}

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) <= 1 {
			continue
		}
		nameLower := strings.ToLower(fields[0])
		if nameLower == "host" {
			if current.Name != "" {
				hosts = append(hosts, current)
			}
			currentUser, err := user.Current()
			currentUserName := ""
			if err == nil {
				currentUserName = currentUser.Username
			}
			current = host{Name: fields[1], HostName: fields[1], User: currentUserName}
		}
		if nameLower == "hostname" {
			current.HostName = fields[1]
		}
		if nameLower == "user" {
			current.User = fields[1]
		}
	}
	if current.Name != "" {
		hosts = append(hosts, current)
	}

	return hosts
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatalf("Usage: %s [configfile]", os.Args[0])
	}
	fileName := args[0]
	lines := strings.Split(readFile(fileName), "\n")
	for i, line := range lines {
		lines[i] = strings.Trim(line, "\r\t ")
	}
	hosts := readHosts(lines)

	searcher := func(input string, index int) bool {
		host := hosts[index]
		name := strings.Replace(strings.ToLower(host.Name+host.User+host.HostName), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   ">{{ .User | cyan }}@({{ .HostName | red }})",
		Inactive: " {{ .User | cyan }}@({{ .HostName | red }})",
		Selected: "{{ .User | cyan }}@({{ .HostName | red }})",
	}

	prompt := promptui.Select{
		Label:     "Select Host/Option",
		Items:     hosts,
		Templates: templates,
		Size:      4,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose number %d: %s\n", i+1, hosts[i].Name)
}
