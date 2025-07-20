/*
MIT License

Copyright (c) 2023-2024 The Trzsz SSH Authors.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package tssh

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func shouldSaveHost(args *sshArgs) bool {
	if args.NoSave {
		return false
	}
	return userConfig.autoSaveHost
}

func saveHostToConfig(args *sshArgs, session *sshClientSession) error {
	if !shouldSaveHost(args) {
		return nil
	}

	user, host, port := parseDestination(args.Destination)
	if host == "" {
		return fmt.Errorf("invalid destination: %s", args.Destination)
	}

	existingHost := findHostInConfig(host)
	if existingHost != "" {
		return updateHostInConfig(existingHost, user, host, port, args.Group)
	}

	return addHostToConfig(user, host, port, args.Group)
}

func findHostInConfig(host string) string {
	loadSshConfig()

	for _, h := range userConfig.allHosts {
		if h.Host == host {
			return h.Alias
		}
	}

	return ""
}

func updateHostInConfig(alias, user, host, port, group string) error {
	configPath := getConfigPath()
	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	hostLine := -1
	hostEndLine := len(lines)
	inHostBlock := false

	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "Host ") {
			if inHostBlock {
				hostEndLine = i
				break
			}
			hostNames := strings.Fields(trimmedLine)[1:]
			for _, name := range hostNames {
				if name == alias {
					hostLine = i
					inHostBlock = true
					break
				}
			}
		}
	}

	if hostLine == -1 {
		return fmt.Errorf("host %s not found in config file", alias)
	}

	var newLines []string
	newLines = append(newLines, lines[:hostLine]...)
	newLines = append(newLines, fmt.Sprintf("Host %s", alias))
	newLines = append(newLines, fmt.Sprintf("    HostName %s", host))
	if user != "" {
		newLines = append(newLines, fmt.Sprintf("    User %s", user))
	}
	if port != "" {
		newLines = append(newLines, fmt.Sprintf("    Port %s", port))
	}
	if group != "" {
		newLines = append(newLines, fmt.Sprintf("    #!! GroupLabels %s", group))
	}
	newLines = append(newLines, lines[hostEndLine:]...)

	err = os.WriteFile(configPath, []byte(strings.Join(newLines, "\n")), 0600)
	if err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	userConfig.loadHosts = sync.Once{}

	fmt.Fprintf(os.Stderr, "\r\nHost %s updated in %s\r\n", alias, configPath)
	return nil
}

func addHostToConfig(user, host, port, group string) error {
	configPath := getConfigPath()

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	file, err := os.OpenFile(configPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString("\n")
	writer.WriteString(fmt.Sprintf("Host %s\n", host))
	writer.WriteString(fmt.Sprintf("    HostName %s\n", host))
	if user != "" {
		writer.WriteString(fmt.Sprintf("    User %s\n", user))
	}
	if port != "" {
		writer.WriteString(fmt.Sprintf("    Port %s\n", port))
	}
	if group != "" {
		writer.WriteString(fmt.Sprintf("    #!! GroupLabels %s\n", group))
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	userConfig.loadHosts = sync.Once{}

	fmt.Fprintf(os.Stderr, "\r\nHost %s added to %s\r\n", host, configPath)
	return nil
}

func getConfigPath() string {
	if userConfig.configPath != "" {
		return userConfig.configPath
	}
	return filepath.Join(userHomeDir, ".ssh", "config")
}

func loadSshConfig() {
	getAllHosts()
}
