/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"io/ioutil"
	"os"
	"os/exec"
	"text/template"

	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		folder, _ := cmd.Flags().GetString("folder")
		release, _ := cmd.Flags().GetString("release")
		download(folder, release)
	},
}

func init() {
	downloadCmd.Flags().StringP("folder", "f", "", "Folder to download artifacts")
	downloadCmd.MarkFlagRequired("folder")
	downloadCmd.Flags().StringP("release", "r", "", "OpenShift release version")
	downloadCmd.MarkFlagRequired("folder")
	rootCmd.AddCommand(downloadCmd)
}

type ImageSet struct {
	Release string
}

func download(folder, release string) {
	tmpDir, err := ioutil.TempDir("", "fp-cli-")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	t, err := template.ParseFiles("imageset.tmpl")
	if err != nil {
		panic(err)
	}

	f, err := os.Create(tmpDir + "/imageset.yaml")
	if err != nil {
		panic(err)
	}

	d := ImageSet{release}
	err = t.Execute(f, d)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("./oc-mirror", "-c", "imageset.yaml", "file://mirror", "--ignore-history", "--dry-run")
	_, err = cmd.Output()
	if err != nil {
		panic(err)
	}

	artifactsCmd := "cat mirror/oc-mirror-workspace/mapping.txt | cut -d \"=\" -f1 > artifacts.txt"
	cmd = exec.Command("bash", "-c", artifactsCmd)
	_, err = cmd.Output()
	if err != nil {
		panic(err)
	}

}
