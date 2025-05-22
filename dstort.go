// dstort - "organize" a directory randomly
// Copyright (C) 2024  Tomás Gutiérrez L. (0x00)

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const (
	version      = "0.1.0"
	license_text = `Copyright (C) 2024 Tomás Gutiérrez L. (0x00)
License GPLv3: GNU GPL version 3 <https://gnu.org/licenses/gpl.html>.
This is free software: you are free to change and redistribute it.
This program comes with ABSOLUTELY NO WARRANTY.
`
)

var directories []string
var files []string

func list_files_and_directories(dir string) {
	file_dir, err := os.Open(dir)
	if err != nil {
		fmt.Printf("Error opening directory %s\n", dir)
		return
	}
	defer file_dir.Close()

	files_dir, err := file_dir.Readdir(-1)
	if err != nil {
		fmt.Printf("Error reading directory %s\n", dir)
		return
	}

	for _, file := range files_dir {
		if file.IsDir() {
			directories = append(directories, filepath.Join(dir, file.Name()))
			list_files_and_directories(filepath.Join(dir, file.Name()))
		} else {
			files = append(files, filepath.Join(dir, file.Name()))
		}
	}

}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomName(originalName string) string {
	ext := filepath.Ext(originalName)

	// Determine a random length for the base name (e.g., 8 to 12 characters)
	// rand.Intn(max - min + 1) + min
	nameLength := rand.Intn(12-8+1) + 8

	randomBaseName := make([]byte, nameLength)
	for i := range randomBaseName {
		randomBaseName[i] = charset[rand.Intn(len(charset))]
	}

	return string(randomBaseName) + ext
}

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	version_flag := flag.Bool("v", false, "Print the version number.")
	help_flag := flag.Bool("h", false, "Print this help message.")

	flag.Parse()

	if *version_flag {
		fmt.Printf("dstort %s\n%s", version, license_text)
		os.Exit(0)
	}

	if *help_flag {
		flag.PrintDefaults()
		os.Exit(0)
	}

	arguments := flag.Args()
	if len(arguments) != 1 {
		fmt.Println("Only 1 directory at a time")
		os.Exit(1)
	}

	dir := arguments[0]
	directories = append(directories, dir)
	list_files_and_directories(dir)

	// Move files to random directories
	for i, file := range files {
		filename := filepath.Base(file)
		newPath := filepath.Join(directories[rand.Intn(len(directories))], filename)
		err := os.Rename(file, newPath)
		if err == nil {
			files[i] = newPath // Update path in slice for subsequent renaming
		} else {
			fmt.Printf("Error moving file %s to %s: %v\n", file, newPath, err)
		}
	}

	// File Renaming Loop
	for i, oldFilePath := range files {
		dirPath := filepath.Dir(oldFilePath)
		originalFileName := filepath.Base(oldFilePath)
		newBaseName := generateRandomName(originalFileName)
		newFilePath := filepath.Join(dirPath, newBaseName)

		// Collision Check
		if _, statErr := os.Stat(newFilePath); !os.IsNotExist(statErr) {
			// Attempt to resolve collision by appending a random number
			fmt.Printf("File rename collision: %s already exists. Trying to append random number.\n", newFilePath)
			newBaseName = generateRandomName(originalFileName[:len(originalFileName)-len(filepath.Ext(originalFileName))] + "_" + fmt.Sprintf("%d", rand.Intn(100))) + filepath.Ext(originalFileName)
			newFilePath = filepath.Join(dirPath, newBaseName)
			if _, statErr2 := os.Stat(newFilePath); !os.IsNotExist(statErr2) {
				fmt.Printf("File rename collision still exists for %s. Skipping rename for %s.\n", newFilePath, oldFilePath)
				continue // Skip if new path still exists
			}
		}

		err := os.Rename(oldFilePath, newFilePath)
		if err == nil {
			files[i] = newFilePath
		} else {
			fmt.Printf("Error renaming file %s to %s: %v\n", oldFilePath, newFilePath, err)
		}
	}

	// Directory Moving Loop
	for i := 1; i < len(directories); i++ { // Start from 1 to skip the root
		oldDirPath := directories[i]
		dirName := filepath.Base(oldDirPath)
		
		// Prevent moving a directory into itself or its current location by ensuring destination is different
		// and not a subdirectory of itself (more complex check, simplified here)
		var newParentDir string
		attempts := 0
		for attempts < len(directories) { // Try to find a suitable parent
			potentialParentIndex := rand.Intn(len(directories))
			potentialParentDir := directories[potentialParentIndex]
			// Check if potential parent is not the directory itself or a subdirectory of the directory being moved
			// This is a simplified check. A full check would be `strings.HasPrefix(potentialParentDir, oldDirPath)`.
			if filepath.Clean(oldDirPath) != filepath.Clean(potentialParentDir) && !strings.HasPrefix(filepath.Clean(potentialParentDir), filepath.Clean(oldDirPath)) {
				newParentDir = potentialParentDir
				break
			}
			attempts++
		}

		if newParentDir == "" {
			fmt.Printf("Could not find a suitable new parent for directory %s. Skipping move.\n", oldDirPath)
			continue
		}

		newDirPath := filepath.Join(newParentDir, dirName)
		
		err := os.Rename(oldDirPath, newDirPath)
		if err == nil {
			// Update all affected paths in directories and files slices
			for j, d := range directories {
				if strings.HasPrefix(d, oldDirPath+string(os.PathSeparator)) {
					directories[j] = filepath.Join(newDirPath, strings.TrimPrefix(d, oldDirPath))
				}
			}
			directories[i] = newDirPath // Update the path of the moved directory itself

			for j, f := range files {
				if strings.HasPrefix(f, oldDirPath+string(os.PathSeparator)) {
					files[j] = filepath.Join(newDirPath, strings.TrimPrefix(f, oldDirPath))
				}
			}
		} else {
			fmt.Printf("Error moving directory %s to %s: %v\n", oldDirPath, newDirPath, err)
		}
	}

	// Directory Renaming Loop
	for i := 1; i < len(directories); i++ { // Start from 1 to skip the root
		oldDirPath := directories[i]
		parentDir := filepath.Dir(oldDirPath)
		originalDirName := filepath.Base(oldDirPath)
		newBaseDirName := generateRandomName(originalDirName)
		newDirPath := filepath.Join(parentDir, newBaseDirName)

		// Collision Check
		if _, statErr := os.Stat(newDirPath); !os.IsNotExist(statErr) {
			fmt.Printf("Directory rename collision: %s already exists. Trying to append random number.\n", newDirPath)
			newBaseDirName = generateRandomName(originalDirName + "_" + fmt.Sprintf("%d", rand.Intn(100)))
			newDirPath = filepath.Join(parentDir, newBaseDirName)
			if _, statErr2 := os.Stat(newDirPath); !os.IsNotExist(statErr2) {
				fmt.Printf("Directory rename collision still exists for %s. Skipping rename for %s.\n", newDirPath, oldDirPath)
				continue
			}
		}

		err := os.Rename(oldDirPath, newDirPath)
		if err == nil {
			// Update all affected paths in directories and files slices
			for j, d := range directories {
				if strings.HasPrefix(d, oldDirPath+string(os.PathSeparator)) {
					directories[j] = filepath.Join(newDirPath, strings.TrimPrefix(d, oldDirPath))
				}
			}
			directories[i] = newDirPath // Update the path of the renamed directory itself

			for j, f := range files {
				if strings.HasPrefix(f, oldDirPath+string(os.PathSeparator)) {
					files[j] = filepath.Join(newDirPath, strings.TrimPrefix(f, oldDirPath))
				}
			}
		} else {
			fmt.Printf("Error renaming directory %s to %s: %v\n", oldDirPath, newDirPath, err)
		}
	}
}
