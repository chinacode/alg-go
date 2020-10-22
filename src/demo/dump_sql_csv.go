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
	"strings"
	"time"
)

type EmailPrefix struct {
	id            uint
	use_status    int8
	email_prefix  string
	success_count uint
	fail_count    uint
}

type UserName struct {
	username string
}
type Email struct {
	email_name    string
	email_prefix  string
	email_name2   sql.NullString
	email_prefix2 sql.NullString
}

var (
	DEBUG         = false
	dotCount      = 0
	noDotCount    = 0
	dateTemplate  = "2006-01-02"          //常规类型
	timeTemplate  = "2006-01-02 15:04:05" //常规类型
	statPrefixMap = make(map[string]int64)
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
		statPrefixMap[prefix.email_prefix] = 0
	}

	if DEBUG {
		log.Println(prefixMap)
	}
	return prefixMap
}

func generateEmail(emailName string, emailPrefix string, prefixMap map[uint]EmailPrefix, apiData [][]string, scriptData [][]string) ([][]string, [][]string, int) {
	emailCount := 0
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
		emailCount++
		statPrefixMap[prefix.email_prefix] = statPrefixMap[prefix.email_prefix] + 1
	}
	return apiData, scriptData, emailCount
}

func writeCsv(rows *sql.Rows, prefixMap map[uint]EmailPrefix, apiData [][]string, scriptData [][]string) ([][]string, [][]string) {
	for rows.Next() {
		var email Email
		err := rows.Scan(&email.email_name, &email.email_prefix, &email.email_name2, &email.email_prefix2)
		if err != nil {
			log.Fatalf("email prefix scan fail ", err)
		}

		emailCount := 0
		//log.Println(email.email_name)
		apiData, scriptData, emailCount = generateEmail(email.email_name, email.email_prefix, prefixMap, apiData, scriptData)
		noDotCount = noDotCount + emailCount

		apiData, scriptData, emailCount = generateEmail(email.email_name2.String, email.email_prefix2.String, prefixMap, apiData, scriptData)
		dotCount = dotCount + emailCount
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
	log.Printf("dump start %s", time.Now().String())
	dbUser := flag.String("user", user, "database user")
	dbPassword := flag.String("password", password, "database password")
	dbHost := flag.String("hostname", host, "database host")
	dbPort := flag.String("port", port, "database port")

	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", *dbUser, *dbPassword, *dbHost, *dbPort, dbName)
	if DEBUG {
		log.Println(dbUrl)
	}
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
		if DEBUG {
			log.Println(sql)
		}
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

	if DEBUG {
		for k, v := range statPrefixMap {
			log.Printf("%s , %d", k, v)
		}
	}

	log.Printf("statistics time (%s~%s) total:%d, api:%d, script:%d, noDot:%d, dot:%d", start, end, total, ac, sc, noDotCount, dotCount)
	log.Printf("dump used time %d ms", (time.Now().UnixNano()-startTime)/1000/1000)
}

func stringSpilt(c rune) bool {
	if c == '-' {
		return true
	} else {
		return false
	}
}

func isLetterAndNumber(username string) bool {
	for _, v := range username {
		if v < 45 || v > 122 {
			return false
		}
	}
	return true
}

func containLetterAndNumber(username string) bool {
	letter := false
	number := false
	for _, v := range username {
		if v >= 45 && v <= 57 {
			number = true
		} else if v >= 65 && v <= 122 {
			letter = true
		}
		if letter && number {
			return true
		}
	}
	return false
}

func isNumber(username string) bool {
	for _, v := range username {
		if v < 45 || v > 57 {
			return false
		}
	}
	return true
}

func isLetter(username string) bool {
	for _, v := range username {
		if v < 65 || v > 122 {
			return false
		}
	}
	return true
}

func getEmailName(username string) string {
	//log.Println(username)
	if !isLetterAndNumber(username) {
		return ""
	}
	//if "a-aziz-hussin-3a794394" == username {
	//	print(username)
	//}
	splitList := strings.FieldsFunc(username, stringSpilt)
	if len(splitList) == 1 {
		if len(splitList[0]) > 23 {
			return ""
		}
		return strings.Join(splitList, "")
	}
	if len(splitList) == 2 {
		return strings.Join(splitList, "")
	}
	//delete last string random by linkedin
	if len(splitList) >= 3 && containLetterAndNumber(splitList[len(splitList)-1]) {
		splitList = splitList[:len(splitList)-1]
	}
	//delete too short string
	for index, string := range splitList {
		if len(string) <= 2 {
			splitList[index] = ""
		}
	}
	var letterList []string
	for _, string := range splitList {
		if isLetter(string) && len(string) > 2 && len(string) < 20 {
			letterList = append(letterList, string)
		}
	}
	if len(letterList) <= 2 {
		return strings.Join(letterList, "")
	}

	return letterList[0] + letterList[1]
}

func dumpUnValidEmail(host string, port string, user string, password string, dbName string, limit string) {
	startTime := time.Now().UnixNano()
	log.Printf("dump start %s", time.Now().String())
	dbUser := flag.String("user", user, "database user")
	dbPassword := flag.String("password", password, "database password")
	dbHost := flag.String("hostname", host, "database host")
	dbPort := flag.String("port", port, "database port")

	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", *dbUser, *dbPassword, *dbHost, *dbPort, dbName)
	if DEBUG {
		log.Println(dbUrl)
	}
	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		log.Fatalf("Could not connect to server: %s\n", err)
	}
	defer db.Close()

	var allData [][]string
	prefixMap := getPrefix(db)

	var prefixList []string
	for _, v := range prefixMap {
		if v.use_status <= 0 || v.email_prefix == "username_count" {
			continue
		}
		prefixList = append(prefixList, v.email_prefix)
	}

	for index := 0; index < 108; index++ {
		sql := fmt.Sprintf("SELECT username FROM likedin_usernames_%d WHERE finished = 1 AND email_name = '' limit %s", index, limit)
		if DEBUG {
			log.Println(sql)
		}
		rows, _ := db.Query(sql)
		if nil == rows {
			continue
		}
		for rows.Next() {
			var username UserName
			err := rows.Scan(&username.username)
			if err != nil {
				log.Fatalf("email prefix scan fail %s", err)
			}
			emailName := getEmailName(username.username)
			if emailName == "" || len(emailName) <= 3 {
				continue
			}
			//log.Printf("%s %s", username, emailName)
			for _, prefix := range prefixList {
				email := fmt.Sprintf("%s@%s", emailName, prefix)
				allData = append(allData, []string{email})
			}
		}
	}

	ac := len(allData)
	name := fmt.Sprintf("dump_email_un_valid_(%s).csv", limit)
	write(name, allData)
	log.Printf("name: %s, count: %d", name, ac)
	log.Printf("dump used time %d ms", (time.Now().UnixNano()-startTime)/1000/1000)
}

func Dump() {
	args := os.Args
	if len(args) != 8 && len(args) != 9 && len(args) != 10 && len(args) != 7 {
		log.Println(args)
		log.Println("please enter method[1:valid,2:unValid] \n" +
			"method[1] host port user password dbName fromTime[optional] toTime[optional] debug[optional default 0]\n" +
			"method[2] host port user password dbName limit debug[optional default 0]")
		return
	}
	if args[1] == "1" {
		if len(args) == 9 {
			args = append(args, "0")
		}
		if len(args) == 7 {
			args = append(args, "")
			args = append(args, "")
			args = append(args, "0")
		}
		DEBUG = args[9] == "1"
		dumpEmail(args[2], args[3], args[4], args[5], args[6], args[7], args[8])
	} else if args[1] == "2" {
		if len(args) == 8 {
			args = append(args, "0")
		}
		DEBUG = args[8] == "1"
		//getEmailName("-0123456789+ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
		//getEmailName("владимир-конельский-144497158")
		//getEmailName("4d-ageng-anom-1a975518b")
		dumpUnValidEmail(args[2], args[3], args[4], args[5], args[6], args[7])

	}
	//for test
	//demo.DumpEmail("192.168.1.200", 3306, "root", "Paramida@2019", "brandu_crawl", "", "")
}
