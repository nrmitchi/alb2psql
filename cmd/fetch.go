// Copyright © 2017 Nick Mitchinson <me@nrmitchi.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	// "bytes"
	"log"
	"fmt"
	"time"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	DATE_FORMAT = "2006-01-02"
)

func FormatS3Dir (date time.Time) string {
	bucket := viper.GetString("bucket")
	accountId := "584106525078"
	region := "us-east-1"
	year, month, day := date.Date()

	return fmt.Sprintf("s3://%s/AWSLogs/%s/elasticloadbalancing/%s/%s/%s/%s/", 
		bucket, accountId, region, 
		fmt.Sprintf("%d", year), 
		fmt.Sprintf("%02d", month), 
		fmt.Sprintf("%02d", day))

}
func FetchDailyLogs (dir string, date time.Time) error {
	// Fetch the S3 logs into the temporary directory, unpack, and psql into the db
	// fmt.Println(FormatS3Dir(date))

	// dir = "/var/folders/38/v061dt3n4q9853qs3lhvw08w0000gn/T/490200550"

	cmd := exec.Command("/usr/local/bin/s3cmd", "get", FormatS3Dir(date), "--recursive")

	// Have this execute in our temporary dir
	cmd.Dir = dir

	// var out bytes.Buffer
	// var stderr bytes.Buffer
	cmd.Stdout = os.Stdout // &out
	cmd.Stderr = os.Stderr // &stderr

	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return err
	}

	// Unzip all of the files
	// Do it one-at-a-time because '*' gets stupidly escaped
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		// Now we want to unzip all of them
		cmd = exec.Command("/usr/bin/gunzip", f.Name())

		// Have this execute in our temporary dir
		cmd.Dir = dir

		cmd.Stdout = os.Stdout // &out
		cmd.Stderr = os.Stderr // &stderr

		err = cmd.Run()
		if err != nil {
			log.Println(err)
			return err
		}

	}

	// Load these into psql
	files, _ = ioutil.ReadDir(dir)
	for _, f := range files {
		cmd = exec.Command("/usr/local/bin/psql", viper.GetString("dbName"), "-c", fmt.Sprintf("copy alb_logs from '%s/%s' DELIMITER ' ' QUOTE '\"' CSV;", dir, f.Name()))

		// Have this execute in our temporary dir
		cmd.Dir = dir

		cmd.Stdout = os.Stdout // &out
		cmd.Stderr = os.Stderr // &stdeff

		log.Println(f.Name())
		err := cmd.Run()
		if err != nil {
			log.Println(err.Error())
			return err
		}
	}

	return  nil
}

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Pull ALB logs for the given date or range",
	//Long: ``,
	Run: func(cmd *cobra.Command, args []string) {

		// Todo: Validate that s3cmd exists: https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go

		var (
			date time.Time
			endDate   time.Time
			err 	  error
		)

		fmt.Printf("Date: %v\n", viper.GetString("date"))
		fmt.Printf("startDate: %v\n", viper.GetString("startDate"))
		fmt.Printf("endDate: %v\n", viper.GetString("endDate"))

		fmt.Printf("dbName: %v\n", viper.GetString("dbName"))
		fmt.Printf("Bucket: %v\n", viper.GetString("bucket"))

		// Lets get a temporary directory to clone the logs into from S3

		// Validate date/startDate/endDate conditions

		if viper.GetString("date") != "" {
			date, err = time.Parse(DATE_FORMAT, viper.GetString("date"))

			if err != nil {
				log.Fatal(fmt.Sprintf("`date` must be of the format %s\n", DATE_FORMAT))
			}

			endDate = time.Now().Local().AddDate(0,0,1)
		} else if viper.GetString("startDate") != "" {
			date, err = time.Parse(DATE_FORMAT, viper.GetString("startDate"))

			if err != nil {
				log.Fatal(fmt.Sprintf("`startDate` must be of the format %s\n", DATE_FORMAT))
			}

			if viper.GetString("endDate") != "" {
				endDate, err = time.Parse(DATE_FORMAT, viper.GetString("endDate"))

				if err != nil {
					log.Fatal(fmt.Sprintf("`endDate` must be of the format %s\n", DATE_FORMAT))
				}
			} else {
				endDate = time.Now().Local().AddDate(0,0,1)
			}

		} else {
			log.Fatal("Either `date` or `startDate` must be provided")
		}

		// Load each date in, and then check if we're at the end
		for {

			fmt.Printf("Loading %v\n", date)

			// Get a temp directory for this day
			dir, err := ioutil.TempDir("", "")
			if err != nil {
				panic(err)
			}
			fmt.Println(dir)
			defer os.RemoveAll(dir) // clean up

			// Fetch the logs
			err = FetchDailyLogs(dir, date)
			if err != nil {
				panic(err)
			}

			// Add 1 day
			date = date.AddDate(0,0,1)

			if date.After(endDate) {
				fmt.Printf("Logs Loaded into %s\n", viper.GetString("dbName"))
				return
			}

		}
		fmt.Printf("date: %v\n", date)
		fmt.Printf("endDate: %v\n", endDate)
	},
}

func init() {
	RootCmd.AddCommand(fetchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fetchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

	fetchCmd.Flags().StringP("date", "d", "", "A single date to fetch")
	fetchCmd.Flags().StringP("startDate", "s", "", "The start of a date range to fetch")
	fetchCmd.Flags().StringP("endDate", "e", "", "The end of a date range to fetch")
	viper.BindPFlag("date", fetchCmd.Flags().Lookup("date"))
	viper.BindPFlag("startDate", fetchCmd.Flags().Lookup("startDate"))
	viper.BindPFlag("endDate", fetchCmd.Flags().Lookup("endDate"))
}
