package demo

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // or the driver of your choice
	_ "github.com/joho/sqltocsv"
	"log"
	"os"
	"time"
)

type EmailPrefix struct {
	id            uint
	use_status    int8
	email_prefix  string
	success_count uint
	fail_count    uint
}

type Email struct {
	email_name   string
	email_prefix string
}

func getPrefix(db *sql.DB) map[uint]EmailPrefix {
	//var prefixMap map[int]string
	rows, _ := db.Query("select * from linkedin_email_prefix")
	prefixMap := make(map[uint]EmailPrefix)
	for rows.Next() {
		var prefix EmailPrefix
		err := rows.Scan(&prefix.id, &prefix.use_status, &prefix.email_prefix, &prefix.success_count, &prefix.fail_count)
		if err != nil {
			log.Fatal(err)
		}
		//log.Println(prefix)
		prefixMap[prefix.id] = prefix
	}

	log.Println(prefixMap)
	return prefixMap
}
func write(fileName string, data [][]string) {
	isNew := false
	f, err := os.Open(fileName)
	if err != nil {
		isNew = true
		f, err = os.Create(fileName)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	if isNew {
		f.WriteString("\xEF\xBB\xBF") // 写入一个UTF-8 BOM
	}
	w := csv.NewWriter(f) //创建一个新的写入文件流
	w.WriteAll(data)
	w.Flush()
}

func writeCsv(rows *sql.Rows, prefixMap map[uint]EmailPrefix, apiData [][]string, scriptData [][]string) ([][]string, [][]string) {
	for rows.Next() {
		var email Email
		err := rows.Scan(&email.email_name, &email.email_prefix)
		if err != nil {
			log.Fatal(err)
		}

		//log.Println(email.email_name)
		var prefixList []uint
		json.Unmarshal([]byte(email.email_prefix), &prefixList)
		for _, id := range prefixList {
			prefix := prefixMap[id]
			email := email.email_name + "@" + prefix.email_prefix
			emails := []string{email}
			if prefix.use_status == 1 {
				apiData = append(apiData, emails)
			} else if prefix.use_status == 2 {
				scriptData = append(scriptData, emails)
			}
		}
	}
	return apiData, scriptData
}

func DumpSql2Csv() {
	startTime := time.Now().UnixNano()
	log.Printf("dump start %d", startTime)
	dbName := "brandu_crawl"
	dbUser := flag.String("user", "root", "database user")
	dbPassword := flag.String("password", "Paramida@2019", "database password")
	dbHost := flag.String("hostname", "192.168.1.200", "database host")
	dbPort := flag.Int("port", 3306, "database port")

	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", *dbUser, *dbPassword, *dbHost, *dbPort, dbName)
	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not connect to server: %s\n", err)
	}
	defer db.Close()

	var apiData [][]string
	var scriptData [][]string

	prefixMap := getPrefix(db)
	for index := 0; index < 108; index++ {
		sql := fmt.Sprintf("SELECT email_name,email_prefix FROM likedin_usernames_%d WHERE email_name != ''", index)
		//log.Println(sql)
		rows, _ := db.Query(sql)
		apiData, scriptData = writeCsv(rows, prefixMap, apiData, scriptData)
		//err := sqltocsv.WriteFile("./dump_email.csv", rows)
		//if err != nil {
		//	panic(err)
		//}
	}

	name := "dump_email_api.csv"
	write(name, apiData)
	log.Printf("name: %s, count: %d", name, len(apiData))

	name = "dump_email_script.csv"
	write(name, scriptData)
	log.Printf("name: %s, count: %d", name, len(scriptData))
	log.Printf("dump used time %d", (time.Now().UnixNano()-startTime)/1000/1000)
}
