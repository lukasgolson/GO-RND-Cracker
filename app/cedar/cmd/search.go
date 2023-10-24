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

		if err != nil {
			panic(err)
		}

		directory, err := cmd.Flags().GetStringArray("directory")
		if err != nil {
			panic(err)
		}

		concurrentTrees, err := cmd.Flags().GetInt("concurrent")
		if err != nil {
			panic(err)
		}

		stride, err := cmd.Flags().GetInt("stride")
		if err != nil {
			panic(err)
		}

		prefetch, err := cmd.Flags().GetBool("prefetch")
		if err != nil {
			panic(err)
		}

		searchDistance, err := cmd.Flags().GetUint32("distance")
		if err != nil {
			panic(err)
		}

		var deliminator string

		if isCSV {
			deliminator = ", "
		} else {
			deliminator = " "
		}

		_, err = os.Stat(inputFile)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("input file '%s' does not exist", inputFile)
				return
			}

		}

		err = application.Search(inputFile, deliminator, directory, concurrentTrees, stride, searchDistance, prefetch)
		if err != nil {
			panic(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringP("input", "i", "numbers.txt", "The input file containing the sequence to search for")
	searchCmd.Flags().BoolP("csv", "e", false, "The input file is CSV deliminated")
	searchCmd.Flags().StringArrayP("directory", "d", []string{"data"}, "The directories to search for seed graphs in")
	searchCmd.Flags().IntP("concurrent", "m", 1, "The number of concurrent trees to search")
	searchCmd.Flags().IntP("stride", "s", 16, "The stride length to use when searching")
	searchCmd.Flags().BoolP("prefetch", "p", false, "Prefetch the current tree into memory before searching")
	searchCmd.Flags().Uint32P("distance", "t", 16, "The maximum distance to search for")

}
