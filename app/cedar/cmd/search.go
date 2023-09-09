/*
Copyright © 2023 Lukas G. Olson <lukasolson64@gmail.com>

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
	"awesomeProject/internal/application"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for the seed that best matches the input sequence.",
	Long:  `Provided with a text file containing a sequence of numbers, Cedar will search for the seed that best matches the input sequence.`,
	Run: func(cmd *cobra.Command, args []string) {

		inputFile, err := cmd.Flags().GetString("input")

		if err != nil {
			panic(err)
		}

		isCSV, err := cmd.Flags().GetBool("csv")

		var deliminator string

		if isCSV {
			deliminator = ", "
		} else {
			deliminator = " "
		}

		cores, _ := cmd.Flags().GetInt("cores")

		_, err = os.Stat(inputFile)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("input file '%s' does not exist", inputFile)
				return
			}

		}

		err = application.Search(inputFile, deliminator, cores)
		if err != nil {
			panic(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	searchCmd.Flags().StringP("input", "i", "numbers.txt", "The input file containing the sequence to search for")
	searchCmd.Flags().BoolP("csv", "e", false, "The input file is CSV deliminated")
}
