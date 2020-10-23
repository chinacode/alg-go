package demo

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // or the driver of your choice
	"github.com/shopspring/decimal"
	"github.com/steakknife/bloomfilter"
	"hash/fnv"
	"io"
	"log"
	"os"
	"strconv"
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
	id       int64
	username string
}

type Email struct {
	email_name    string
	email_prefix  string
	email_name2   sql.NullString
	email_prefix2 sql.NullString
}

type EmailData struct {
	email_name    string
	email_prefix  []uint
	email_name2   string
	email_prefix2 []uint
	unValid       bool
}

var (
	DEBUG        = false
	dotCount     = 0
	noDotCount   = 0
	SPLIT        = rune('-')
	dateTemplate = "2006-01-02"          //常规类型
	timeTemplate = "2006-01-02 15:04:05" //常规类型

	maxElements   = uint64(500 * 10000)
	probCollide   = 0.0000001
	statPrefixMap = make(map[string]int64)

	repeatCount        = 0
	bloomInstance, err = bloomfilter.NewOptimal(maxElements, probCollide) //check repeat email
)

func write(fileName string, data [][]string) {
	os.Remove(fileName)
	//isNew := false
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	//if isNew {
	f.WriteString("\xEF\xBB\xBF") // 写入一个UTF-8 BOM
	//}
	w := csv.NewWriter(f) //创建一个新的写入文件流
	w.WriteAll(data)
	w.Flush()
}

func readCsv(filename string) [][]string {
	var lines [][]string
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Error:", err)
		return lines
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			continue
			//fmt.Println("Error:", err)
			//return lines
		}
		lines = append(lines, record)
	}
	return lines
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
		hash := fnv.New64()
		hash.Write([]byte(email))
		if bloomInstance.Contains(hash) {
			repeatCount++
			continue
		}
		emails := []string{email}
		if prefix.use_status == 1 {
			apiData = append(apiData, emails)
		} else if prefix.use_status == 2 {
			scriptData = append(scriptData, emails)
		}
		emailCount++
		statPrefixMap[prefix.email_prefix] = statPrefixMap[prefix.email_prefix] + 1

		bloomInstance.Add(hash)
	}
	return apiData, scriptData, emailCount
}

func generateEmailList(rows *sql.Rows, prefixMap map[uint]EmailPrefix, apiData [][]string, scriptData [][]string) ([][]string, [][]string) {
	for rows.Next() {
		var email Email
		err := rows.Scan(&email.email_name, &email.email_prefix, &email.email_name2, &email.email_prefix2)
		if err != nil {
			log.Fatalf("email prefix scan fail %s", err)
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
		apiData, scriptData = generateEmailList(rows, prefixMap, apiData, scriptData)
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

	log.Printf("statistics time (%s~%s) total:%d, api:%d, script:%d, noDot:%d, dot:%d, repeat %d", start, end, total, ac, sc, noDotCount, dotCount, repeatCount)
	log.Printf("dump used time %d ms", (time.Now().UnixNano()-startTime)/1000/1000)
}

func stringSpilt(c rune) bool {
	if c == SPLIT {
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

func countNumberLetter(username string) (int8, int8) {
	var n int8
	var l int8
	for _, v := range username {
		if v >= 45 && v <= 57 {
			n++
		} else if v >= 65 && v <= 122 {
			l++
		}
	}
	return n, l
}

func getEmailNames(nameList []string) []string {
	if len(nameList) <= 0 {
		return []string{}
	}
	if len(nameList) == 1 {
		return []string{nameList[0]}
	}
	if len(nameList) == 2 {
		return []string{strings.Join(nameList, ""), strings.Join(nameList, ".")}
	}
	nameList = nameList[:2]
	return []string{strings.Join(nameList, ""), strings.Join(nameList, ".")}
}

func isMixing(username string) bool {
	n, l := countNumberLetter(username)
	value, _ := decimal.NewFromFloat(float64(n)).Div(decimal.NewFromFloat(float64(l + n))).Float64()
	percent := value > 0.20 && value < 0.5

	last := 0
	change := 0
	for _, v := range username {
		tmpLast := 0
		if v >= 45 && v <= 57 { //number check
			tmpLast = 1
		} else if v >= 65 && v <= 122 { //letter check
			tmpLast = 2
		}
		if last == 0 {
			last = tmpLast
		}
		if last != tmpLast {
			change++
			//println(string(v))
			last = tmpLast
		}
	}
	//println(username, percent, change)
	return percent || change >= 3
}

func getEmailName(username string) []string {
	//log.Println(username)
	var emails []string
	if !isLetterAndNumber(username) || len(username) <= 3 {
		return emails
	}
	splitList := strings.FieldsFunc(username, stringSpilt)
	if len(splitList) == 1 {
		if len(splitList[0]) > 23 || isMixing(splitList[0]) {
			return emails
		}
		return getEmailNames(splitList)
	}
	//delete last string random by linkedin
	if len(splitList) >= 3 && containLetterAndNumber(splitList[len(splitList)-1]) {
		splitList = splitList[:len(splitList)-1]
	}
	//delete too short string
	var letterList []string
	for _, string := range splitList {
		if len(string) <= 2 || isMixing(string) || len(string) > 23 {
			continue
		}
		letterList = append(letterList, string)
	}
	if len(letterList) <= 2 {
		return getEmailNames(letterList)
	}

	return getEmailNames(letterList)
}

func dumpUnValidEmail(host string, port string, user string, password string, dbName string, status string, limit string) {
	SPLIT = rune('-')
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

	//for email checker
	var emailData [][]string
	//for import data
	var allData [][]string
	prefixMap := getPrefix(db)

	var prefixList []string
	for _, v := range prefixMap {
		if v.use_status <= 0 || v.email_prefix == "username_count" {
			continue
		}
		if "1" == status && v.use_status != 1 {
			continue
		}
		if "2" == status && v.use_status != 2 {
			continue
		}
		prefixList = append(prefixList, v.email_prefix)
	}

	for index := 0; index < 108; index++ {
		sql := fmt.Sprintf("SELECT id,username FROM likedin_usernames_%d WHERE finished = 1 AND email_name = '' limit %s", index, limit)
		if DEBUG {
			log.Println(sql)
		}
		rows, _ := db.Query(sql)
		if nil == rows {
			continue
		}
		for rows.Next() {
			var username UserName
			err := rows.Scan(&username.id, &username.username)
			if err != nil {
				log.Fatalf("email prefix scan fail %s", err)
			}
			emailNames := getEmailName(username.username)
			if len(emailNames) <= 0 {
				continue
			}
			indexString := strconv.Itoa(index)
			idString := strconv.FormatInt(username.id, 10)
			tmpEmails := []string{indexString, idString}
			//log.Printf("%s %s", username, emailName)
			for _, emailName := range emailNames {
				tmpEmails = append(tmpEmails, emailName)
				for _, prefix := range prefixList {
					email := fmt.Sprintf("%s@%s", emailName, prefix)

					hash := fnv.New64()
					hash.Write([]byte(email))
					if bloomInstance.Contains(hash) {
						repeatCount++
						continue
					}

					emailData = append(emailData, []string{email})
					bloomInstance.Add(hash)
				}
			}
			allData = append(allData, tmpEmails)
		}
	}

	name := fmt.Sprintf("dump_email_checker_(%s).csv", limit)
	write(name, emailData)
	name = fmt.Sprintf("dump_email_import_(%s).csv", limit)
	write(name, allData)
	log.Printf("name: %s, user count: %d, email count:%d, repeat %d", name, len(allData), len(emailData), repeatCount)
	log.Printf("dump used time %d ms", (time.Now().UnixNano()-startTime)/1000/1000)
}

func importEmail(host string, port string, user string, password string, dbName string, indexFile string, importFile string) {
	SPLIT = rune('@')
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

	prefixMap := getPrefix(db)
	prefixIdMap := map[string]uint{}
	for k, v := range prefixMap {
		prefixIdMap[v.email_prefix] = k
	}

	updateTime := time.Now().Unix()
	indexList := readCsv(indexFile)
	namesMap := make(map[string][]string)
	for _, indexEmail := range indexList {
		key := strings.TrimSpace(indexEmail[2])
		if nil != namesMap[key] {
			log.Println(indexEmail)
		}
		id := indexEmail[1]
		index := indexEmail[0]
		if !isNumber(index) {
			index = "0"
		}
		namesMap[key] = []string{index, id}
	}

	failCount := 0
	successCount := 0
	emailName := ""
	var emailData EmailData
	tx, _ := db.Begin()
	importList := readCsv(importFile)
	for _, importEmail := range importList {
		emails := strings.FieldsFunc(importEmail[0], stringSpilt)
		if "" == emailName {
			emailName = emails[0]
		}

		prefix := prefixIdMap[emails[1]]
		if strings.Contains(emails[0], ".") {
			emailData.email_name2 = emails[0]
			emailData.email_prefix2 = append(emailData.email_prefix2, prefix)
			//replace .
			emails[0] = strings.Replace(emails[0], ".", "", len(emails[0]))
		} else {
			emailData.email_name = emails[0]
			emailData.email_prefix = append(emailData.email_prefix, prefix)
		}
		if emailName != emails[0] {
			//log.Println(emailName)
			//emailName = "hambali"
			//log.Println(emailName, namesMap["hambali"])
			namesIndex := namesMap[emailName]
			if nil == namesIndex {
				emailData = EmailData{}
				emailName = emails[0]
				continue
			}
			id := namesIndex[1]
			tableIndex := namesIndex[0]
			prefixJson, _ := json.Marshal(emailData.email_prefix)
			prefix2Json, _ := json.Marshal(emailData.email_prefix2)
			prefix2JsonStr := string(prefix2Json)
			if prefix2JsonStr == "null" {
				prefix2JsonStr = "[]"
			}

			sql := fmt.Sprintf(
				"update likedin_usernames_%s set finished = 1,email_name = '%s', email_prefix = '%s',email_name2 = '%s', email_prefix2 = '%s',update_time = %d where id = %s limit 1",
				tableIndex, emailData.email_name, string(prefixJson), emailData.email_name2, prefix2JsonStr, updateTime, id)
			if DEBUG {
				log.Println(sql)
			}

			//_, err := db.Exec(sql)
			_, err := tx.Exec(sql)
			if nil != err {
				log.Fatalln(err)
			}

			//count, err := ret.RowsAffected()
			successCount++

			//reset
			emailData = EmailData{}
			emailName = emails[0]
			//reset names Map for finish fail
			namesMap[emailName] = nil
		}
	}
	tx.Commit()

	tx, _ = db.Begin()
	for _, v := range namesMap {
		if nil == v {
			//ok data
			continue
		}
		id := v[1]
		tableIndex := v[0]
		sql := fmt.Sprintf("update likedin_usernames_%s set finished = 1,update_time = %d where id = %s limit 1", tableIndex, updateTime, id)
		if DEBUG {
			log.Println(sql)
		}
		//_, err := db.Exec(sql)
		_, err := tx.Exec(sql)
		if nil != err {
			log.Fatalln(err)
		}
		//count, err := ret.RowsAffected()
		failCount++
	}
	tx.Commit()

	log.Printf("statstics success %d, fail %d ", successCount, failCount)
	log.Printf("import used time %d ms", (time.Now().UnixNano()-startTime)/1000/1000)
}

func Dump() {
	readme := "please enter method[1:dump valid email,2:dump un valid email,3:import valid email ] \n" +
		"method[1] host port user password dbName fromTime[optional] toTime[optional] debug[optional default 0]\n" +
		"method[2] host port user password dbName status[0:all] limit debug[optional default 0]\n" +
		"method[3] host port user password dbName indexFile importFile debug[optional default 0]"
	args := os.Args
	if len(args) != 8 && len(args) != 9 && len(args) != 10 && len(args) != 7 {
		log.Println(args)
		log.Println(readme)
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
		if len(args) == 9 {
			args = append(args, "0")
		}
		DEBUG = args[9] == "1"
		//getEmailName("-0123456789+ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
		//getEmailName("владимир-конельский-144497158")
		//getEmailName("4d-ageng-anom-1a975518b")
		//println(isMixing("5a50a6162"))
		//println(isMixing("411884187"))
		//println(isMixing("a18b61117"))
		dumpUnValidEmail(args[2], args[3], args[4], args[5], args[6], args[7], args[8])
	} else if args[1] == "3" {
		if len(args) == 9 {
			args = append(args, "0")
		}
		DEBUG = args[9] == "1"
		//getEmailName("-0123456789+ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
		//getEmailName("владимир-конельский-144497158")
		//getEmailName("4d-ageng-anom-1a975518b")
		importEmail(args[2], args[3], args[4], args[5], args[6], args[7], args[8])
		//println(isMixing("5a50a6162"))
	} else {
		log.Println(readme)
	}
	//for test
	//demo.DumpEmail("192.168.1.200", 3306, "root", "Paramida@2019", "brandu_crawl", "", "")
}
