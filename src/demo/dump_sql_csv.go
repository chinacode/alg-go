package demo

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // or the driver of your choice
	"github.com/shopspring/decimal"
	"github.com/steakknife/bloomfilter"
	"hash/fnv"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
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
	tableIndex    string
	id            string
	email_name    string
	email_prefix  []uint
	email_name2   string
	email_prefix2 []uint
}

type CheckData struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
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
	repeatList         []string
	bloomInstance, err = bloomfilter.NewOptimal(maxElements, probCollide) //check repeat email
)

func Write(fileName string, data [][]string) {
	WriteCsv(fileName, data, true)
}

func WriteCsv(fileName string, data [][]string, addBom bool) {
	os.Remove(fileName)
	//isNew := false
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	//if isNew {
	// 写入一个UTF-8 BOM
	if addBom {
		f.WriteString("\xEF\xBB\xBF")
	}
	//}
	w := csv.NewWriter(f) //创建一个新的写入文件流
	w.WriteAll(data)
	w.Flush()
}

func ReadCsv(filename string) [][]string {
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

func getDB(mysql MysqlServer) *sql.DB {
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", mysql.user, mysql.password, mysql.host, mysql.port, mysql.database)
	if DEBUG {
		log.Println(dbUrl)
	}
	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		log.Fatalf("Could not connect to server: %s\n", err)
	}
	defer db.Close()
	return db
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
			if DEBUG && len(repeatList) <= 10 {
				repeatList = append(repeatList, email)
			}
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

func DumpEmailSimple(mysql MysqlServer, start string, end string, finished int) (string, int) {
	DEBUG = false
	return dumpEmail(mysql.host, strconv.Itoa(mysql.port), mysql.user, mysql.password, mysql.database, start, end, finished)
}

func dumpEmail(host string, port string, user string, password string, dbName string, start string, end string, finished int) (string, int) {
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
	//dbUser := flag.String("user", user, "database user")
	//dbPassword := flag.String("password", password, "database password")
	//dbHost := flag.String("hostname", host, "database host")
	//dbPort := flag.String("port", port, "database port")
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbName)
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
		//sql := fmt.Sprintf("SELECT email_name,email_prefix,IFNULL(email_name2,'') email_name2,IFNULL(email_prefix2,'[]') email_prefix2 FROM linkedin_usernames_%d WHERE 1 = 1 ", index)
		sql := fmt.Sprintf("SELECT email_name,email_prefix,email_name2,email_prefix2 FROM linkedin_usernames_%d WHERE 1 = 1 ", index)
		if _start.Unix() > 0 {
			sql = sql + fmt.Sprintf(" AND update_time >= %d", _start.Unix())
		}
		if _end.Unix() > 0 {
			sql = sql + fmt.Sprintf(" AND update_time < %d", _end.Unix())
		}
		sql = sql + " AND email_name != ''"
		if finished >= 0 {
			sql = sql + fmt.Sprintf(" AND finished = %d", finished)
		}
		sql = sql + " ORDER BY id "
		if DEBUG {
			logger.Infof(sql)
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
	Write(name, apiData)
	log.Printf("name: %s, count: %d", name, ac)

	//name = fmt.Sprintf("dump_email_script_(%s~%s).csv", start, end)
	//if start+end == "" {
	//	name = "dump_email_script_(all).csv"
	//}
	//Write(name, scriptData)
	log.Printf("name: %s, count: %d", name, sc)

	if DEBUG {
		for k, v := range statPrefixMap {
			log.Printf("%s , %d", k, v)
		}
	}

	if len(repeatList) > 0 {
		log.Printf("repeat list samples %s", repeatList)
	}
	log.Printf("statistics time (%s~%s) total:%d, api:%d, script:%d, noDot:%d, dot:%d, repeat %d", start, end, total, ac, sc, noDotCount, dotCount, repeatCount)
	log.Printf("dump used time %d ms", (time.Now().UnixNano()-startTime)/1000/1000)

	return name, len(apiData)
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

func getEmailNames(nameList []string, returnPart bool) []string {
	if len(nameList) > 3 {
		nameList = nameList[:3]
	}
	//return username depart
	if returnPart {
		return nameList
	}
	if len(nameList) <= 0 {
		return []string{}
	}
	if len(nameList) == 1 {
		return []string{nameList[0]}
	}
	if len(nameList) == 2 {
		return []string{strings.Join(nameList, ""), strings.Join(nameList, ".")}
	}
	return []string{strings.Join(nameList, ""), strings.Join(nameList, ".")}
}

func StringSum(username string) uint64 {
	sum := uint64(0)
	for _, v := range username {
		sum += uint64(v)
	}
	return sum
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

func permute(nums []string) [][]string {
	if len(nums) == 1 {
		return [][]string{nums}
	}
	var ret [][]string
	for i := 0; i < len(nums); i++ {
		buf := make([]string, len(nums)-1)
		copy(buf, nums[0:i])
		copy(buf[i:], nums[i+1:])
		r := permute(buf)
		for _, v := range r {
			ret = append(ret, append(v, nums[i]))
		}
	}
	return ret
}

func GetEmailName(username string, returnPart bool) []string {
	//log.Println(username)
	var emails []string
	if !isLetterAndNumber(username) || len(username) <= 3 {
		return emails
	}
	SPLIT = rune('-')
	splitList := strings.FieldsFunc(username, stringSpilt)
	if len(splitList) == 1 {
		if len(splitList[0]) > 23 || len(splitList[0]) < 5 || isMixing(splitList[0]) {
			return emails
		}
		return getEmailNames(splitList, returnPart)
	}
	//delete last string random by linkedin
	if len(splitList) >= 3 && containLetterAndNumber(splitList[len(splitList)-1]) {
		splitList = splitList[:len(splitList)-1]
	}
	//delete too short string
	var letterList []string
	for _, string := range splitList {
		if len(string) <= 1 || isMixing(string) || len(string) > 23 {
			continue
		}
		letterList = append(letterList, string)
	}
	lastPart := letterList[len(letterList)-1]
	if isNumber(lastPart) && len(lastPart) > 5 {
		letterList = letterList[:len(letterList)-1]
	}
	if len(letterList) <= 2 {
		return getEmailNames(letterList, returnPart)
	}

	return getEmailNames(letterList, returnPart)
}

func GetEmailEffectiveRankList(username string) []string {
	departs := GetEmailName(username, true)
	if len(departs) <= 1 {
		return departs
	}

	endPrefixList := []string{"gmail.com", "yahoo.com", "outlook.com", "hotmail.com", "live.com"}
	emailBucket := "email_ok"
	emailOKBloomFilter := bloomfilterMap[emailBucket]
	emailBucket = "email_fail"
	emailFailBloomFilter := bloomfilterMap[emailBucket]
	if nil == emailFailBloomFilter || nil == emailOKBloomFilter {
		return getEmailNames(departs, false)
	}

	getNotExistsCount := func(emailPrefixList []string) int {
		notExists := 0
		for _, emailPrefix := range emailPrefixList {
			for _, endPrefix := range endPrefixList {
				email := emailPrefix + "@" + endPrefix
				hash := fnv.New64()
				hash.Write([]byte(email))
				if !emailOKBloomFilter.Contains(hash) && !emailFailBloomFilter.Contains(hash) {
					notExists++
				}
				hash.Reset()
			}
		}
		return notExists
	}

	prefixMap := map[int][]string{}
	possibilities := permute(departs)
	for _, parts := range possibilities {
		email := strings.Join(parts, "")
		dotEmail := strings.Join(parts, ".")
		count := getNotExistsCount([]string{email, dotEmail})
		prefixMap[count] = parts
	}

	var keys []int
	//if len(prefixMap) > 1 {
	//	println(prefixMap)
	//}
	for key := range prefixMap {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	selectParts := prefixMap[keys[len(keys)-1]]
	return getEmailNames(selectParts, false)
}

func dumpUnValidEmailApi(mysql MysqlServer, status string, finished string, limit string) ([][]string, [][]string) {
	return dumpUnValidEmail(mysql.host, strconv.Itoa(mysql.port), mysql.user, mysql.password, mysql.database, status, finished, limit, false)
}

func dumpUnValidEmail(host string, port string, user string, password string, dbName string, status string, finished string, limit string, writeFile bool) ([][]string, [][]string) {
	bloomInstance, err = bloomfilter.NewOptimal(maxElements, probCollide)
	startTime := time.Now().UnixNano()
	log.Printf("dump start %s", time.Now().String())
	//dbUser := flag.String("user", user, "database user")
	//dbPassword := flag.String("password", password, "database password")
	//dbHost := flag.String("hostname", host, "database host")
	//dbPort := flag.String("port", port, "database port")
	//dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", *dbUser, *dbPassword, *dbHost, *dbPort, dbName)

	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbName)
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

	discardList := make(map[int][]string)
	whereSql := " "
	if finished != "0" {
		whereSql = " AND update_time > 0"
	}
	for index := 0; index < 108; index++ {
		sql := fmt.Sprintf("SELECT id,username FROM linkedin_usernames_%d WHERE finished = %s AND email_name = '' %s limit %s", index, finished, whereSql, limit)
		if DEBUG {
			logger.Infof(sql)
		}
		rows, _ := db.Query(sql)
		if nil == rows {
			continue
		}
		var discardTableList []string
		for rows.Next() {
			var username UserName
			err := rows.Scan(&username.id, &username.username)
			if err != nil {
				log.Fatalf("email prefix scan fail %s", err)
			}
			emailNames := GetEmailName(username.username, false)
			if len(emailNames) <= 0 {
				discardTableList = append(discardTableList, strconv.FormatInt(username.id, 10))
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

			//add empty dot email
			if len(emailNames) == 1 {
				tmpEmails = append(tmpEmails, "")
			}
			tmpEmails = append(tmpEmails, username.username)
			allData = append(allData, tmpEmails)
		}

		discardList[index] = discardTableList
		discardTableList = nil
	}

	//check bloom server
	var newEmailData [][]string
	allIndex := 0
	pageSize := 1000
	_page, _ := decimal.NewFromInt(int64(len(emailData))).Div(decimal.NewFromInt(int64(pageSize))).Float64()
	pageCount := int(math.Ceil(_page))
	for page := 0; page < pageCount; page++ {
		var pageList [][]string
		start := page * pageSize
		end := (page + 1) * pageSize
		if end > len(emailData) {
			end = len(emailData)
		}
		pageList = emailData[start:end]
		logger.Infof("start %d, end %d", start, end)
		var emailList []string
		for _, email := range pageList {
			emailList = append(emailList, email[0])
		}

		var checkList = bloomRequest("http://"+config.bloom.host+":"+strconv.Itoa(config.bloom.port)+"/batch_exists?bucket=email_ok", emailList)
		for _, index := range checkList {
			if index == 1 {
				allIndex++
				continue
			}
			newEmailData = append(newEmailData, emailData[allIndex])
			allIndex++
		}
		emailList = nil
		checkList = nil
	}

	//update can't generate email user id
	tx, _ := db.Begin()
	for index, discards := range discardList {
		if len(discards) <= 0 {
			continue
		}
		sql := fmt.Sprintf("UPDATE linkedin_usernames_%d set update_time = 0 where id IN (%s) AND finished = 1", index, strings.Join(discards, ","))
		//log.Println(sql)
		ret, err := db.Exec(sql)
		if nil != err {
		}
		rowsAffected, _ := ret.RowsAffected()
		logger.Infof("table index %d, ret %s, length %d ", index, rowsAffected, len(discards))
	}
	tx.Commit()
	discardList = nil

	name := fmt.Sprintf("dump_email_checker_(%s).csv", limit)
	if writeFile {
		Write(name, newEmailData)

		name = fmt.Sprintf("dump_email_import_(%s).csv", limit)
		Write(name, allData)
	}

	logger.Infof("name: %s, user-count: %d, email-count:%d, un-exists-email-count:%d, repeat %d", name, len(allData), len(emailData), len(newEmailData), repeatCount)
	logger.Infof("dump used time %d ms", (time.Now().UnixNano()-startTime)/1000/1000)
	return allData, newEmailData
}

func GetEmailCount(mysql MysqlServer, finished int, countryId int, gender int, start int64, end int64) (int, int) {
	startTime := time.Now().UnixNano()
	log.Printf("dump start %s", time.Now().String())

	emailCount := 0
	usernameCount := 0
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", mysql.user, mysql.password, mysql.host, mysql.port, mysql.database)
	if DEBUG {
		log.Println(dbUrl)
	}
	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		log.Fatalf("Could not connect to server: %s\n", err)
	}
	defer db.Close()

	for index := 0; index < 108; index++ {
		sql := fmt.Sprintf("SELECT email_name,email_prefix,email_name2,email_prefix2 FROM linkedin_usernames_%d WHERE 1 = 1 ", index)
		if finished == -1 {
			sql += " AND finished > 0 "
		} else if finished == 0 {
			sql += fmt.Sprintf(" AND finished = %d ", finished)
		} else {
			sql += fmt.Sprintf(" AND email_name != '' AND finished = %d ", finished)
		}
		if finished != 0 {
			sql += fmt.Sprintf(" AND update_time > %d AND update_time < %d", start, end)
		}
		if countryId > 0 {
			sql += fmt.Sprintf(" AND country_id = %d", countryId)
		}
		if gender > 0 {
			sql += fmt.Sprintf(" AND gender = %d", gender)
		}

		//if DEBUG {
		//	log.Println(sql)
		//}
		rows, _ := db.Query(sql)
		if nil == rows {
			continue
		}

		for rows.Next() {
			usernameCount++

			var email Email
			err := rows.Scan(&email.email_name, &email.email_prefix, &email.email_name2, &email.email_prefix2)
			if err != nil {
				log.Fatalf("email prefix scan fail %s", err)
			}

			var prefixList []uint
			if "" != email.email_prefix {
				json.Unmarshal([]byte(email.email_prefix), &prefixList)
				emailCount += len(prefixList)
			}
			if "" != email.email_prefix2.String {
				json.Unmarshal([]byte(email.email_prefix2.String), &prefixList)
				emailCount += len(prefixList)
			}
		}
	}

	log.Printf("dump used time %d ms", (time.Now().UnixNano()-startTime)/1000/1000)
	return usernameCount, emailCount
}

func bloomRequest(url string, emailList []string) []int8 {
	var checkList []int8
	//use memory bloom
	memoryBloomCheck := func(emailList []string) []int8 {
		bucket = "email_ok"
		bucketBloomFilter := bloomfilterMap[bucket]
		logger.Errorf("bloom filter is nil bucket is {}", bucket)

		for _, email := range emailList {
			if nil == bucketBloomFilter {
				checkList = append(checkList, 0)
				continue
			}
			hash := fnv.New64()
			hash.Write([]byte(email))
			exists := bucketBloomFilter.Contains(hash)
			if exists {
				checkList = append(checkList, 1)
			} else {
				checkList = append(checkList, 0)
			}
			hash.Reset()
		}
		return checkList
	}

	//remote checker
	requestBloomCheck := func(url string, emailList []string) []int8 {
		response, _ := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(strings.Join(emailList, ",")))
		jsonData, _ := ioutil.ReadAll(response.Body)
		var checkData CheckData
		err := json.Unmarshal(jsonData, &checkData)
		if nil != err {
			log.Printf("decode check json fail %s", err)
		}

		if len(checkData.Data) <= 11 {
			count, _ := strconv.ParseInt(checkData.Data, 10, 64)
			checkList = append(checkList, int8(count))
			return checkList
		}

		err = json.Unmarshal([]byte(checkData.Data), &checkList)
		if nil != err {
			log.Printf("decode check json fail %s", err)
		}
		return checkList
	}

	//add use url
	if strings.Contains(url, "add") {
		checkList = requestBloomCheck(url, emailList)
	} else {
		checkList = memoryBloomCheck(emailList)
	}

	return checkList
}

func importEmailApi(mysql MysqlServer, indexFile string, importFile string, forceImport bool) (int, int, int, int, int) {
	//return 0, 0, 0, 0, 0
	return importEmail(mysql.host, strconv.Itoa(mysql.port), mysql.user, mysql.password, mysql.database, indexFile, importFile, forceImport)
}

func importEmail(host string, port string, user string, password string, dbName string, indexFile string, importFile string, forceImport bool) (int, int, int, int, int) {
	bloomInstance, err = bloomfilter.NewOptimal(maxElements, probCollide)
	startTime := time.Now().UnixNano()
	logger.Infof("import start %s", time.Now().String())
	//dbUser := flag.String("user", user, "database user")
	//dbPassword := flag.String("password", password, "database password")
	//dbHost := flag.String("hostname", host, "database host")
	//dbPort := flag.String("port", port, "database port")

	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbName)
	if DEBUG {
		log.Println(dbUrl)
	}
	db, err := sql.Open("mysql", dbUrl)

	if err != nil {
		logger.Errorf("Could not connect to server: %s\n", err)
	}
	defer db.Close()

	prefixMap := getPrefix(db)
	prefixIdMap := map[string]uint{}
	for k, v := range prefixMap {
		prefixIdMap[v.email_prefix] = k
	}

	namesRepeatCount := 0
	updateTime := time.Now().Unix()
	indexList := ReadCsv(indexFile)
	namesMap := make(map[string][]string)
	for _, indexEmail := range indexList {
		if len(indexEmail) < 3 {
			logger.Infof("jump email %s", indexEmail)
		}
		key := strings.TrimSpace(indexEmail[2])
		if nil != namesMap[key] {
			//log.Println(indexEmail)
			namesRepeatCount++
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
	emailCount := 0
	emailRepeatCount := 0
	importList := ReadCsv(importFile)

	if nil == importList || len(importList) == 0 {
		log.Fatalf("csv file is empty. %s !!!", importFile)
		return successCount, failCount, namesRepeatCount, emailCount, emailRepeatCount
	}

	executeSuccess := func() {
		log.Println("depart data and check data")
		updateMap := make(map[string]*EmailData)
		pageSize := 1000
		_page, _ := decimal.NewFromInt(int64(len(importList))).Div(decimal.NewFromInt(int64(pageSize))).Float64()
		pageCount := int(math.Ceil(_page))
		for page := 0; page < pageCount; page++ {
			var pageList [][]string
			start := page * pageSize
			end := (page + 1) * pageSize
			if end > len(importList) {
				end = len(importList)
			}
			pageList = importList[start:end]

			log.Printf("start %d, end %d", start, end)
			var emailList []string
			for _, email := range pageList {
				emailList = append(emailList, email[0])
			}

			var checkList []int8
			if !forceImport {
				checkList = bloomRequest("http://"+config.bloom.host+":"+strconv.Itoa(config.bloom.port)+"/batch_exists?bucket=email_ok", emailList)
			}
			for index, email := range emailList {
				if !forceImport && checkList[index] == 1 {
					emailRepeatCount++
					continue
				}
				emailCount++
				if page == 0 && index == 0 {
					email = email[3:]
				}
				SPLIT = rune('@')
				emails := strings.FieldsFunc(email, stringSpilt)

				emailName := emails[0]
				prefix := prefixIdMap[emails[1]]
				containDot := false
				if strings.Contains(emailName, ".") {
					containDot = true
					emailName = strings.Replace(emailName, ".", "", len(emailName))
				}

				emailData, exists := updateMap[emailName]
				if !exists {
					emailData = &EmailData{}
					namesIndex := namesMap[emailName]
					if nil == namesIndex {
						continue
					}
					emailData.tableIndex = namesIndex[0]
					emailData.id = namesIndex[1]
					updateMap[emailName] = emailData
				}

				if containDot {
					emailData.email_name2 = emails[0] //use dot origin string
					emailData.email_prefix2 = append(emailData.email_prefix2, prefix)
				} else {
					emailData.email_name = emailName
					emailData.email_prefix = append(emailData.email_prefix, prefix)
				}
			}
			//set bloom filter server
			bloomRequest("http://"+config.bloom.host+":"+strconv.Itoa(config.bloom.port)+"/batch_add?bucket=email_ok", emailList)
		}

		log.Println("import data to database ")
		rowIndex := 0
		tx, _ := db.Begin()
		for emailName, emailData := range updateMap {
			id := emailData.id
			tableIndex := emailData.tableIndex
			prefixJson, _ := json.Marshal(emailData.email_prefix)
			prefix2Json, _ := json.Marshal(emailData.email_prefix2)
			prefix2JsonStr := string(prefix2Json)
			if prefix2JsonStr == "null" {
				prefix2JsonStr = "[]"
			}

			sql := fmt.Sprintf(
				"update linkedin_usernames_%s set finished = 2,email_name = '%s', email_prefix = '%s',email_name2 = '%s', email_prefix2 = '%s',update_time = %d where id = %s limit 1",
				tableIndex, emailData.email_name, string(prefixJson), emailData.email_name2, prefix2JsonStr, updateTime, id)
			//if DEBUG {
			//	log.Println(sql)
			//}

			//_, err := db.Exec(sql)
			_, err := tx.Exec(sql)
			if nil != err {
				log.Fatalln(err)
			}

			//count, err := ret.RowsAffected()
			successCount++

			//reset
			emailData = nil
			//reset names Map for finish fail
			namesMap[emailName] = nil

			if rowIndex%2000 == 0 {
				tx.Commit()
				log.Printf("import success email part commit index %d.", rowIndex)
				tx, _ = db.Begin()
			}
			rowIndex++
		}
		tx.Commit()
	}

	executeFail := func() {
		//set un finish data
		tx, _ := db.Begin()
		for index, email := range importList {
			emailCount++
			SPLIT = rune('@')
			emails := strings.FieldsFunc(email[0], stringSpilt)

			emailName := emails[0]
			if strings.Contains(emailName, ".") {
				emailName = strings.Replace(emailName, ".", "", len(emailName))
			}
			failCount++

			namesIndex, exists := namesMap[emailName]
			if !exists {
				continue
			}
			//delete key for one key execute once
			delete(namesMap, emailName)
			id := namesIndex[1]
			tableIndex := namesIndex[0]
			//for update fail and if is success then jump it
			sql := fmt.Sprintf("update linkedin_usernames_%s set finished = 2,update_time = %d where id = %s AND email_name = '' limit 1", tableIndex, updateTime, id)
			//if DEBUG {
			//	log.Println(sql)
			//}
			_, err := tx.Exec(sql)
			if nil != err {
				log.Fatalln(err)
			}

			if index%5000 == 0 {
				tx.Commit()
				log.Printf("import fail email part commit index %d.", index)
				tx, _ = db.Begin()
			}
		}
		tx.Commit()
	}

	if strings.Contains(importFile, "_fail") {
		executeFail()
	} else {
		executeSuccess()
	}

	log.Printf("statstics success %d, fail %d, names repeat: %d, email count: %d, email repeat count: %d", successCount, failCount, namesRepeatCount, emailCount, emailRepeatCount)
	log.Printf("import used time %d ms", (time.Now().UnixNano()-startTime)/1000/1000)
	return successCount, failCount, namesRepeatCount, emailCount, emailRepeatCount
}

func Dump() {
	readme := "please enter method[1:dump valid email,2:dump un valid email,3:import valid email ] \n" +
		"method[1] host port user password dbName fromTime[optional] toTime[optional] debug[optional default 0]\n" +
		"method[2] host port user password dbName status[0:all] limit debug[optional default 0]\n" +
		"method[3] host port user password dbName indexFile importFile debug[optional default 0]"

	//GetEmailCount(config.mysql, 2, 1604246400,1604332800 )
	//return

	DumpEmailSimple(config.mysql, "2020-11-18", "2020-11-24", -1)
	return

	args := os.Args
	//println(stringSum("be-creative") % 108)
	//return
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
		dumpEmail(args[2], args[3], args[4], args[5], args[6], args[7], args[8], -1)
	} else if args[1] == "2" {
		if len(args) == 9 {
			args = append(args, "0")
		}
		DEBUG = args[9] == "1"
		//GetEmailName("-0123456789+ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
		//GetEmailName("владимир-конельский-144497158")
		//GetEmailName("4d-ageng-anom-1a975518b")
		//println(isMixing("5a50a6162"))
		//println(isMixing("411884187"))
		//println(isMixing("a18b61117"))
		dumpUnValidEmail(args[2], args[3], args[4], args[5], args[6], args[7], "1", args[8], true)
	} else if args[1] == "3" {
		if len(args) == 9 {
			args = append(args, "0")
		}
		DEBUG = args[9] == "1"
		//GetEmailName("-0123456789+ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
		//GetEmailName("владимир-конельский-144497158")
		//GetEmailName("4d-ageng-anom-1a975518b")
		importEmail(args[2], args[3], args[4], args[5], args[6], args[7], args[8], false)
		//println(isMixing("5a50a6162"))
	} else {
		log.Println(readme)
	}

}
