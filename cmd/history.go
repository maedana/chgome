// Copyright © 2017 maedana <maeda.na@gmail.com>
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
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// historyCmd represents the history command
var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "標準出力にGoogle Chromeの閲覧履歴を表示",
	Run: func(cmd *cobra.Command, args []string) {

		// chromeが動いてるとHistoryがロックされていて読めないのでコピーする
		b, err := ioutil.ReadFile(viper.GetString("history_db_path"))
		if err != nil {
			panic(err)
		}
		var useDbPath = os.Getenv("HOME") + "/.chgome/History"
		err = ioutil.WriteFile(useDbPath, b, 0744)
		if err != nil {
			panic(err)
		}

		db, err := sql.Open("sqlite3", useDbPath)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// visit_timeが1601年からのマイクロ秒となっているのでunixtime(1701年からの秒)に変換している点に注意
		rows, err := db.Query("select title, urls.url as url, (visits.visit_time - 11676312000000000)/1000/1000 as unixtime from visits inner join urls on visits.url = urls.id order by unixtime desc")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var url string
			var title string
			var unixTime int64
			err = rows.Scan(&url, &title, &unixTime)
			if err != nil {
				log.Fatal(err)
			}
			var visitAt string
			visitAt = time.Unix(unixTime, 0).String()
			fmt.Println(visitAt + "|" + title + "|" + url)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	RootCmd.AddCommand(historyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// historyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// historyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
