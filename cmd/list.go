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
	"path/filepath"

	diskfs "github.com/diskfs/go-diskfs"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
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
		path, err := os.Getwd()

		ListPartitions(path + "/rpi.img")

		//if err != nil {
		//	log.Fatal(err)
		//}
		//files, err := ioutil.ReadDir(path)
		//
		//if err != nil {
		//	log.Fatal(err)
		//}
		//
		//for _, f := range files {
		//	if filepath.Ext(f.Name()) == ".img" {
		//		fmt.Println(f.Name())
		//	}
		//}

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
