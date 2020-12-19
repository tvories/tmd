/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"fmt"
	"io"
	"os"

	//"path/filepath"

	diskfs "github.com/diskfs/go-diskfs"
	"github.com/spf13/cobra"
	//"io/ioutil"
	"github.com/gookit/color"
	"log"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if source != "" && dest != "" {
			fmt.Println("source and dest are defined. moving on")
			fmt.Println("source: " + source + "\ndest: " + dest)

			// Now we copy
			wd, err := os.Getwd()
			if err != nil {
				log.Panic(err)
			}
			//TODO: Figure out how to do this per OS type
			var src = wd + "/" + source
			var dst = wd + "/" + dest
			fmt.Println("value of source path: " + src)
			fmt.Println("value of dest path: " + dst)
			err = CopyFile(src, dst)
			if err != nil {
				color.Error.Println("Copy failed!\n")
				log.Panic(err)
			} else {
				color.Success.Println("Successfully copied!  Path: " + dst)
			}

		} else {
			color.Error.Println("You need both --source and --dest defined")
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func ReadFilesystem(p string) {
	disk, err := diskfs.Open(p)
	if err != nil {
		log.Panic(err)
	}

	fs, err := disk.GetFilesystem(0) // assuming the whole disk, so partition = 0
	if err != nil {
		log.Panic(err)
	}
	files, err := fs.ReadDir("/") // this should list everything at the root
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(files)
}

func ListPartitions(imgFile string) {
	disk, err := diskfs.Open(imgFile)
	if err != nil {
		log.Panic(err)
	}

	partitions, err := disk.GetPartitionTable()
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(partitions)
}

func ReadContents() {
	path, err := os.Getwd() // gets current working directory

	if err != nil {
		log.Panic(err)
	}

	disk, err := diskfs.Open(path + "/rpi.img")
	if err != nil {
		log.Panic(err)
	}

	fs, err := disk.GetFilesystem(1)
	if err != nil {
		log.Panic(err)
	}

	files, err := fs.ReadDir("/")
	if err != nil {
		log.Panic(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}
}

func copy(src, dst string, BUFFERSIZE int64) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		log.Panic(err)
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		log.Panic(err)
	}
	defer source.Close()

	_, err = os.Stat(dst)
	if err == nil {
		return fmt.Errorf("file %s already exists", dst)
	}

	destination, err := os.Create(dst)
	if err != nil {
		log.Panic(err)
	}
	defer destination.Close()

	buf := make([]byte, BUFFERSIZE)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return err
}

// TODO: Add existing file check
// TODO: Add overwrite option
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
