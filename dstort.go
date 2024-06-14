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

func main() {
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

	for _, file := range files {
		filename := filepath.Base(file)
		os.Rename(file, directories[rand.Intn(len(directories))]+"/"+filename)
	}

	for i := 1; i < len(directories); i++ {
		dst := directories[rand.Intn(len(directories))] + "/" + filepath.Base(directories[i])
		err := os.Rename(directories[i], dst)
		if err == nil {
			directories[i] = dst
		}
	}
}
