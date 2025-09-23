/*
Copyright 2025 Finvi, Ontario Systems

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

package internal

import (
	"os"
	"os/user"
	"path"
	"strings"
)

// ExpandPath gives back an expanded, clean path
// relative paths are expanded and ~/ replaced with the home directory
func ExpandPath(filePath string) string {
	filePath = path.Clean(filePath)
	if strings.HasPrefix(filePath, "~") {
		rest := filePath[2:]
		filePath = path.Join(getHomeFolder(), rest)
	}
	return filePath
}

func getHomeFolder() string {
	current, err := user.Current()
	if err != nil {
		if home := os.Getenv("HOME"); home != "" {
			return home
		}
		return os.Getenv("USERPROFILE")
	}
	return current.HomeDir
}
