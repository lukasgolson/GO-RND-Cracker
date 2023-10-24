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
	"github.com/spf13/cobra"
	"math"
	"runtime"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Cedar by generating the required lookup graphs.",
	Long: `Initialize Cedar by generating the required lookup graphs.

Note that this command will take a long time to complete, and will use a lot of disk space (over 250 GB).`,
	Run: func(cmd *cobra.Command, args []string) {

		coreCount, err := cmd.Flags().GetInt("cores")

		if err != nil {
			panic(err)
		}

		fileCount, err := cmd.Flags().GetInt("files")

		if err != nil {
			panic(err)
		}

		seedCount, err := cmd.Flags().GetInt64("seedCount")
		if err != nil {
			panic(err)
		}

		directory, err := cmd.Flags().GetStringArray("directory")
		if err != nil {
			panic(err)
		}

		sequenceHigh, err := cmd.Flags().GetInt("sequenceHigh")
		if err != nil {
			panic(err)
		}

		sequenceOffset, err := cmd.Flags().GetInt("sequenceOffset")
		if err != nil {
			panic(err)
		}

		err = application.Initialize(coreCount, fileCount, seedCount, sequenceHigh, sequenceOffset, directory)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initCmd.Flags().IntP("files", "f", runtime.NumCPU(), "Number of file partitions to split the seed space into")
	initCmd.Flags().Int64P("seedCount", "s", math.MaxInt32, "Upper bound on the number of seeds to generate")

	initCmd.Flags().IntP("sequenceHigh", "h", 100, "Upper bound on the sequence elements")
	initCmd.Flags().IntP("sequenceOffset", "o", 1, "Offset to add to the sequence elements")

	initCmd.Flags().StringArrayP("directory", "d", []string{"data"}, "The directories to store the lookup graphs in")
}
