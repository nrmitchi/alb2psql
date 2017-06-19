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
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)
/*

CREATE TABLE alb_logs (
    Http varchar(5),
    RequestTime TIMESTAMP WITH TIME ZONE,
    ELBName varchar(100),
    RequestIP_Port varchar(22),
    BackendIP_Port varchar(22),
    RequestProcessingTime FLOAT,
    BackendProcessingTime FLOAT,
    ClientResponseTime FLOAT,
    ELBResponseCode varchar(3),
    BackendResponseCode varchar(3),
    ReceivedBytes BIGINT,
    SentBytes BIGINT,
    HttpRequest varchar(5083),
    UserAgent varchar(500),
    SSL_Cipher varchar(40),
    SSL_Protocol varchar(40),
    TargetGroup varchar(100),
    RequestId varchar(40)
);
*/

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the alb2psql database",
	//Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Building Database")
		c := exec.Command("/usr/local/bin/createdb", viper.GetString("dbName"))

		c.Stdout = os.Stdout // &out
		c.Stderr = os.Stderr // &stdeff

		err := c.Run()
		if err != nil {
			log.Println(err.Error())
			return
		}

		c = exec.Command("/usr/local/bin/psql", viper.GetString("dbName"),  "-c", fmt.Sprintf(`
		CREATE TABLE alb_logs (
		    Http varchar(5),
		    RequestTime TIMESTAMP WITH TIME ZONE,
		    ELBName varchar(100),
		    RequestIP_Port varchar(22),
		    BackendIP_Port varchar(22),
		    RequestProcessingTime FLOAT,
		    BackendProcessingTime FLOAT,
		    ClientResponseTime FLOAT,
		    ELBResponseCode varchar(3),
		    BackendResponseCode varchar(3),
		    ReceivedBytes BIGINT,
		    SentBytes BIGINT,
		    HttpRequest varchar(5083),
		    UserAgent varchar(500),
		    SSL_Cipher varchar(40),
		    SSL_Protocol varchar(40),
		    TargetGroup varchar(100),
		    RequestId varchar(40)
		);
		`))

		c.Stdout = os.Stdout // &out
		c.Stderr = os.Stderr // &stdeff

		err = c.Run()
		if err != nil {
			log.Println(err.Error())
			return
		}
		log.Println("Database created: alb-logs")
		log.Println("Table created: alb_logs")
		log.Println("  Warning: No index created")
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
