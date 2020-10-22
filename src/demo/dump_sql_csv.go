package demo

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // or the driver of your choice
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
	email_name    string
	email_prefix  string
	email_name2   sql.NullString
	email_prefix2 sql.NullString
}

var (
	timeTemplate = "2006-01-02 15:04:05" //常规类型
	dateTemplate = "2006-01-02"          //常规类型
)

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

func generateEmail(emailName string, emailPrefix string, prefixMap map[uint]EmailPrefix, apiData [][]string, scriptData [][]string) ([][]string, [][]string) {
	var prefixList []uint
	json.Unmarshal([]byte(emailPrefix), &prefixList)
	for _, id := range prefixList {
		prefix := prefixMap[id]
		email := emailName + "@" + prefix.email_prefix
		emails := []string{email}
		if prefix.use_status == 1 {
			apiData = append(apiData, emails)
		} else if prefix.use_status == 2 {
			scriptData = append(scriptData, emails)
		}
	}
	return apiData, scriptData
}

func writeCsv(rows *sql.Rows, prefixMap map[uint]EmailPrefix, apiData [][]string, scriptData [][]string) ([][]string, [][]string) {
	for rows.Next() {
		var email Email
		err := rows.Scan(&email.email_name, &email.email_prefix, &email.email_name2, &email.email_prefix2)
		if err != nil {
			log.Fatalf("email prefix scan fail ", err)
		}

		//log.Println(email.email_name)
		apiData, scriptData = generateEmail(email.email_name, email.email_prefix, prefixMap, apiData, scriptData)
		apiData, scriptData = generateEmail(email.email_name2.String, email.email_prefix2.String, prefixMap, apiData, scriptData)
	}
	return apiData, scriptData
}

func dumpEmail(host string, port string, user string, password string, dbName string, start string, end string) {
	var _start time.Time
	var _end time.Time
	if "" != start {
		var err error
		_start, err = time.ParseInLocation(timeTemplate, start, time.Local)
		if nil != err {
			_start, err = time.ParseInLocation(dateTemplate, start, time.Local)
		}
	}
	if "" != end {
		var err error
		_end, err = time.ParseInLocation(timeTemplate, end, time.Local)
		if nil != err {
			_end, err = time.ParseInLocation(dateTemplate, end, time.Local)
		}
	}

	startTime := time.Now().UnixNano()
	log.Printf("dump start %d", startTime)
	dbUser := flag.String("user", user, "database user")
	dbPassword := flag.String("password", password, "database password")
	dbHost := flag.String("hostname", host, "database host")
	dbPort := flag.String("port", port, "database port")

	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", *dbUser, *dbPassword, *dbHost, *dbPort, dbName)
	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		log.Fatalf("Could not connect to server: %s\n", err)
	}
	defer db.Close()

	var apiData [][]string
	var scriptData [][]string

	prefixMap := getPrefix(db)
	for index := 0; index < 108; index++ {
		//sql := fmt.Sprintf("SELECT email_name,email_prefix,IFNULL(email_name2,'') email_name2,IFNULL(email_prefix2,'[]') email_prefix2 FROM likedin_usernames_%d WHERE 1 = 1 ", index)
		sql := fmt.Sprintf("SELECT email_name,email_prefix,email_name2,email_prefix2 FROM likedin_usernames_%d WHERE 1 = 1 ", index)
		if _start.Unix() > 0 {
			sql = sql + fmt.Sprintf(" AND update_time >= %d", _start.Unix())
		}
		if _end.Unix() > 0 {
			sql = sql + fmt.Sprintf(" AND update_time < %d", _end.Unix())
		}
		sql = sql + " AND email_name != ''"
		//log.Println(sql)
		rows, _ := db.Query(sql)
		if nil == rows {
			continue
		}
		apiData, scriptData = writeCsv(rows, prefixMap, apiData, scriptData)
	}

	ac := len(apiData)
	sc := len(scriptData)
	total := ac + sc
	name := fmt.Sprintf("dump_email_api_(%s~%s).csv", start, end)
	if start+end == "" {
		name = "dump_email_api_(all).csv"
	}
	write(name, apiData)
	log.Printf("name: %s, count: %d", name, ac)

	name = fmt.Sprintf("dump_email_script_(%s~%s).csv", start, end)
	if start+end == "" {
		name = "dump_email_script_(all).csv"
	}
	write(name, scriptData)
	log.Printf("name: %s, count: %d", name, sc)
	log.Printf("statistics time (%s~%s) total:%d, api:%d, script:%d", start, end, total, ac, sc)
	log.Printf("dump used time %d ms", (time.Now().UnixNano()-startTime)/1000/1000)
}

func Dump() {
	args := os.Args
	if len(args) != 8 && len(args) != 6 {
		log.Println(args)
		log.Println("please enter host port user password dbName fromTime[optional] toTime[optional]")
		return
	}
	if len(args) == 6 {
		args = append(args, "")
		args = append(args, "")
	}
	dumpEmail(args[1], args[2], args[3], args[4], args[5], args[6], args[7])
	//for test
	//demo.DumpEmail("192.168.1.200", 3306, "root", "Paramida@2019", "brandu_crawl", "", "")
}
