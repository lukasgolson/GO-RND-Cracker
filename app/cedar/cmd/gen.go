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
	"math/rand"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "From a given seed, generates a sequence of numbers.",
	Run: func(cmd *cobra.Command, args []string) {

		seed, err := cmd.Flags().GetInt("seed")

		if err != nil {
			panic(err)
		}

		length, err := cmd.Flags().GetInt("length")

		if err != nil {
			panic(err)
		}

		csv, err := cmd.Flags().GetBool("csv")

		if err != nil {
			panic(err)
		}

		high, err := cmd.Flags().GetInt("sequenceHigh")
		if err != nil {
			panic(err)
		}

		offset, err := cmd.Flags().GetInt("sequenceOffset")
		if err != nil {
			panic(err)
		}

		sequence := application.GenerateRandomSequence(int64(seed), int64(length), high, offset, rand.New(rand.NewSource(0)))

		if csv {
			csvString, err := application.FormatByteArrayAsCSV(sequence)
			if err != nil {
				panic(err)
			}
			fmt.Println(csvString)
			return
		} else {
			fmt.Println(application.FormatByteArrayAsNumbers(sequence, 16))
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(genCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// genCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// genCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	genCmd.Flags().IntP("seed", "s", 0, "The seed to use when generating the sequence")
	genCmd.Flags().IntP("length", "l", 32, "The length of the sequence to generate")
	genCmd.Flags().BoolP("csv", "e", false, "Format the output as a CSV record")

	genCmd.Flags().IntP("sequenceHigh", "h", 100, "Upper bound on the sequence elements")
	genCmd.Flags().IntP("sequenceOffset", "o", 1, "Offset to add to the sequence elements")
}
