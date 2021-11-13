/*
Copyright © 2021 Drew Stinnett <drew@drewlink.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	_ "github.com/drewstinnett/labdoc/internal/plugins/all"
	"github.com/drewstinnett/labdoc/pkg/labdoc"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate template_file",
	Short: "Generate README.md",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		templateIn, err := ioutil.ReadFile(args[0])
		cobra.CheckErr(err)

		allFunctions := template.FuncMap{}

		// Load All plugin data
		for name, c := range labdoc.Plugins {
			itemFunction, err := c().TemplateFunctions()
			cobra.CheckErr(err)
			for k, v := range itemFunction {
				funcName := fmt.Sprintf("%v%v", name, strings.Title(k))
				allFunctions[funcName] = v
			}
		}

		cobra.CheckErr(err)
		tpl, err := template.New("tpl").Funcs(allFunctions).Parse(string(templateIn))
		cobra.CheckErr(err)
		out := os.Stdout
		err = tpl.Execute(out, nil)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
