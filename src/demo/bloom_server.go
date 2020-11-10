package demo

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/pierrec/lz4"
	"github.com/steakknife/bloomfilter"
	hash2 "hash"
	"hash/fnv"
	"io"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
	"util"
)

var (
	//_maxElements         = uint64(500 * 10000)
	//_probCollide         = 0.0000001

	key          = ""
	keys         = ""
	bucket       = ""
	bucketBuffer = ""
	forceAdd     = uint64(0)
	percent      = 0.00001
	elements     = uint64(1 * 10000 * 10000)

	backDumpPrefix       = "BF"
	bloomfilterMemoryMap = map[string]int64{}
	bloomfilterMap       = map[string]*bloomfilter.Filter{}

	logLevel  = "debug"
	logConfig = `
<seelog type="asynctimer" asyncinterval="1000000" minlevel="` + logLevel + `" maxlevel="error">
    <outputs formatid="main">
        <console/>
        <splitter formatid="format1">
            <file path="./logs/bloom_server.log"/>
        </splitter>
    </outputs>
    <formats>
        <format id="main" format="%Date(2006-1-02/3:04:05.0000) [%LEVEL] %File(%Line) - %Msg%n"/>
        <format id="format1" format="%Date(2006-1-02/3:04:05.0000) [%LEVEL] %File(%Line) - %Msg%n"/>
    </formats>
</seelog>
`
	logger, _ = seelog.LoggerFromConfigAsString(logConfig)

	serverPort = 9002
	serverName = "bloomfilter server "

	saveTimer = time.NewTimer(time.Second * 2)
)

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type Zip struct {
	filename string
	filter   *bloomfilter.Filter
}

func (z *Zip) write() (n int64, err error) {
	num := int64(z.filter.N())
	start := time.Now().UnixNano()
	logger.Infof("BACKUP save bucket %s start %d", z.filename, start)
	w, err := os.Create(z.filename)
	if err != nil {
		logger.Errorf("BACKUP open file %s error %s", z.filename, err)
		return num, nil
	}
	defer func() {
		err = w.Close()
	}()

	//rawW := lz4.NewWriter(w)
	rawW := gzip.NewWriter(w)
	defer func() {
		err = rawW.Close()
	}()
	content, err := z.filter.MarshalBinary()
	if nil != err {
		logger.Errorf("BACKUP filter binary error %s", err)
		return -1, err
	}
	logger.Infof("BACKUP filter size %d", len(content))
	intN, err := rawW.Write(content)
	logger.Infof("BACKUP filter bucket %s size %d, use time %d ms", z.filename, intN, (time.Now().UnixNano()-start)/1000/1000)

	return num, err
}

func Main() {
	runtime.GOMAXPROCS(2)
	seelog.ReplaceLogger(logger)
	defer seelog.Flush()

	http.HandleFunc("/", Index)
	http.HandleFunc("/config", SetConfig)
	http.HandleFunc("/add", Add)
	http.HandleFunc("/batch_add", BatchAdd)
	http.HandleFunc("/exists", Exists)
	http.HandleFunc("/batch_exists", BatchExists)
	http.HandleFunc("/clear", Clear)
	http.HandleFunc("/save", Save)
	http.HandleFunc("/load", Load)
	http.HandleFunc("/memory", Memory)
	http.HandleFunc("/buckets", Buckets)
	http.HandleFunc("/key_count", KeyCount)
	http.HandleFunc("/dump_email_zip", dumpEmailZip)
	http.HandleFunc("/import_email", importEmailData)
	http.HandleFunc("/get_email_count", getEmailCount)

	go func() {
		loadOldData()

		startSaveTask()
	}()

	logger.Info(serverName + " start success " + strconv.Itoa(serverPort))
	logger.Info(http.ListenAndServe(":"+strconv.Itoa(serverPort), nil))
}

func loadOldData() {
	pwd, _ := os.Getwd()
	fileInfoList, err := ioutil.ReadDir(pwd)
	if err != nil {
		logger.Error(err)
	}
	fmt.Println(len(fileInfoList))
	for i := range fileInfoList {
		filename := fileInfoList[i].Name()
		if strings.HasPrefix(filename, backDumpPrefix) {
			start := time.Now().UnixNano()
			logger.Infof("LOAD bucket %s start %d", filename, start)

			r, err := os.Open(filename)
			if err != nil {
				logger.Errorf("LOAD file error %s ", err)
			}
			defer func() {
				err = r.Close()
			}()

			SPLIT = rune('.')
			list := strings.FieldsFunc(filename, stringSpilt)
			if len(list) != 3 {
				logger.Infof("LOAD filename %s jump", filename)
				continue
			}

			_bucket := list[1]
			bloomfilterMap[_bucket], err = bloomfilter.NewOptimal(elements, percent)
			_, err = bloomfilterMap[_bucket].ReadFrom(r)
			if nil != err {
				logger.Errorf("LOAD file data error %s", err)
			}
			//初始化
			bloomfilterMemoryMap[_bucket] = int64(bloomfilterMap[_bucket].N())
			logger.Infof("LOAD file %s, bucket %s , count %d, num %d, memory %d, use time %d ms",
				filename, _bucket,
				bloomfilterMap[_bucket].K(),
				bloomfilterMap[_bucket].N(),
				bloomfilterMap[_bucket].M(),
				(time.Now().UnixNano()-start)/1000/1000,
			)
			runtime.GC()
		}
	}
}

func backupFilter(filename string, filter *bloomfilter.Filter) int64 {
	num := int64(filter.N())
	start := time.Now().UnixNano()
	logger.Infof("BACKUP save bucket %s start %d", filename, start)
	w, err := os.Create(filename)
	if err != nil {
		logger.Errorf("BACKUP open file %s error %s", filename, err)
		return num
	}
	defer func() {
		err = w.Close()
	}()
	//rawW := gzip.NewWriter(w)
	//defer func() {
	//	err = rawW.Close()
	//}()
	rawW := lz4.NewWriter(w)
	defer func() {
		err = rawW.Close()
	}()
	//filter.WriteFile(filename)
	content, err := filter.MarshalBinary()
	if nil != err {
		logger.Errorf("BACKUP filter binary error %s", err)
	}
	logger.Infof("BACKUP filter size %d", len(content))
	intN, err := rawW.Write(content)

	logger.Infof("BACKUP filter bucket %s size %d, use time %d ms", filename, intN, (time.Now().UnixNano()-start)/1000/1000)
	return num
}

func startSaveTask() {
	go func() {
		for range saveTimer.C {
			//logger.Infof("timer process ")
			for bucket, filter := range bloomfilterMap {
				num := int64(filter.N())
				filename := fmt.Sprintf("%s.%s.bucket", backDumpPrefix, bucket)

				currentMemory := bloomfilterMemoryMap[bucket]
				if currentMemory != num && currentMemory >= 0 {
					// for repeat save
					bloomfilterMemoryMap[bucket] = -1

					////this way will lock then bloom server stop the world
					//num, err = filter.WriteFile(filename)
					//if nil != err {
					//	logger.Errorf("save file %s fail", filename)
					//}

					////this way use lz4 zip algorithm speed is the best, size will bigger 50%
					//zip := Zip{filename: filename, filter: filter}
					//num, _ = zip.Write()

					num = backupFilter(filename, filter)
					bloomfilterMemoryMap[bucket] = num
				}
			}
			runtime.GC()
			saveTimer.Reset(time.Second * 5 * 60)
		}
	}()
}

func download(fileName string, response http.ResponseWriter) {
	file, _ := os.Open(fileName)
	defer file.Close()
	fileHeader := make([]byte, 512)
	file.Read(fileHeader)
	fileStat, _ := file.Stat()
	response.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	response.Header().Set("Content-Type", http.DetectContentType(fileHeader))
	response.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))
	file.Seek(0, 0)
	io.Copy(response, file)
}

func getParams(request *http.Request, key string) string {
	query := request.URL.Query()
	value, exists := query[key]
	if !exists {
		return ""
	}
	return value[0]
}

func getErrorResult(msg string) *Result {
	result := Result{Code: http.StatusInternalServerError, Msg: msg, Data: ""}
	return &result
}

func getResult(data interface{}) *Result {
	var result Result
	if reflect.TypeOf(data).Name() != "string" {
		jsonData, _ := json.Marshal(data)
		result = Result{Code: http.StatusOK, Msg: "success", Data: string(jsonData)}
	} else {
		result = Result{Code: http.StatusOK, Msg: "success", Data: fmt.Sprint(data)}
	}
	return &result
}

func responseSuccess(response http.ResponseWriter, data interface{}) {
	result := getResult(data)
	json, err := json.Marshal(&result)
	if nil != err {
		logger.Error(err)
	}
	response.WriteHeader(result.Code)
	response.Write(json)

	logger.Debugf("<<<response %s", json)
}

func responseError(response http.ResponseWriter, msg string) {
	result := getErrorResult(msg)
	json, _ := json.Marshal(&result)
	response.WriteHeader(result.Code)
	response.Write(json)

	logger.Debugf("<<<response %s", json)
}

func initParams(response http.ResponseWriter, request *http.Request) (bool, map[string]bool) {
	key = ""
	keys = ""
	bucket = ""
	bucketBuffer = ""
	forceAdd = uint64(0)
	percent = 0.00001
	elements = uint64(1 * 10000 * 10000)

	params := make(map[string]bool)
	query := request.URL.Query()

	paramKey := "key"
	if nil != query[paramKey] {
		key = query[paramKey][0]
		params[paramKey] = true
	} else {
		params[paramKey] = false
	}

	paramKey = "keys"
	bodyData := ""
	if nil != request.Body {
		bodyData, _ := ioutil.ReadAll(request.Body)
		keys = string(bodyData)
		params[paramKey] = true
	}
	if nil != query[paramKey] {
		keys = query[paramKey][0]
		params[paramKey] = true
	}
	if "" == keys {
		params[paramKey] = false
	}

	paramKey = "bucket"
	if nil != query[paramKey] {
		bucket = query[paramKey][0]
		bucketBuffer = bucket + ".BUFFER"
		params[paramKey] = true
	} else {
		params[paramKey] = false
	}

	paramKey = "percent"
	if nil != query[paramKey] {
		percent, _ = strconv.ParseFloat(query[paramKey][0], 10)
		params[paramKey] = true
	} else {
		params[paramKey] = false
	}

	paramKey = "elements"
	if nil != query[paramKey] {
		elements, _ = strconv.ParseUint(query[paramKey][0], 10, 64)
		params[paramKey] = true
	} else {
		params[paramKey] = false
	}

	paramKey = "forceAdd"
	if nil != query[paramKey] {
		forceAdd, _ = strconv.ParseUint(query[paramKey][0], 10, 64)
		params[paramKey] = true
	} else {
		params[paramKey] = false
	}

	logger.Debugf(">>>request url %s , body %s", request.RequestURI, bodyData)
	response.Header().Add("server", "bloom server")

	if !strings.Contains(request.RequestURI, "config") &&
		!strings.Contains(request.RequestURI, "buckets") &&
		!strings.Contains(request.RequestURI, "dump") &&
		!strings.Contains(request.RequestURI, "import") {
		if "" == bucket || len(bucket) <= 0 {
			responseError(response, fmt.Sprintf("params bucket must set!"))
			return false, params
		}
		if bloomfilterMap[bucket] == nil {
			//responseError(response, fmt.Sprintf("bucket %s not exists, please use /config init first", bucket))
			//return false, params
			bloomfilterMap[bucket], err = bloomfilter.NewOptimal(elements, percent)
		}
	}

	return true, params
}

func Index(response http.ResponseWriter, request *http.Request) {
	index := `
"/"
"/config"
"/add"
"/batch_add"
"/exists"
"/batch_exists"
"/clear"
"/save"
"/load"
"/memory"
"/buckets"
"/key_count"
	`
	response.Write([]byte(index))
}

func SetConfig(response http.ResponseWriter, request *http.Request) {
	check, _ := initParams(response, request)
	if !check {
		return
	}
	if bloomfilterMap[bucket] != nil {
		responseError(response, fmt.Sprintf("bucket %s exists, please reset first", bucket))
		return
	}
	//if !params["percent"] || !params["elements"] {
	//	responseError(response, fmt.Sprintf("percent(%f) and elements(%d) must set!", percent, elements))
	//	return
	//}
	bloomfilterMap[bucket], err = bloomfilter.NewOptimal(elements, percent)
	responseSuccess(response, fmt.Sprintf("success"))
}

func Load(response http.ResponseWriter, request *http.Request) {
	check, _ := initParams(response, request)
	if !check {
		return
	}
	filename := fmt.Sprintf("%s.%s.bucket", backDumpPrefix, bucket)
	r, err := os.Open(filename)
	if err != nil {
		responseError(response, fmt.Sprintf("back file %s not exists", filename))
		return
	}
	defer func() {
		err = r.Close()
	}()

	count, err := bloomfilterMap[bucket].ReadFrom(r)
	responseSuccess(response, fmt.Sprintf("load from file %s success count %d", filename, count))
}

func Save(response http.ResponseWriter, request *http.Request) {
	check, _ := initParams(response, request)
	if !check {
		return
	}

	start := time.Now().UnixNano()
	filename := fmt.Sprintf("%s.%s.bucket", backDumpPrefix, bucket)
	bloomfilterMap[bucket].WriteFile(filename)
	useTime := (time.Now().UnixNano() - start) / 1000 / 1000

	responseSuccess(response, fmt.Sprintf("save to file %s success, count %d, use time %d ms", filename, bloomfilterMap[bucket].N(), useTime))
}

func BatchAdd(response http.ResponseWriter, request *http.Request) {
	check, params := initParams(response, request)
	if !check {
		return
	}
	if !params["keys"] {
		responseError(response, fmt.Sprintf("keys must set!"))
		return
	}

	count := 0
	SPLIT = rune(',')
	var hash hash2.Hash64
	list := strings.FieldsFunc(keys, stringSpilt)
	for _, keyString := range list {
		if len(keyString) <= 0 {
			continue
		}
		hash = fnv.New64()
		_, err := hash.Write([]byte(keyString))
		if nil != err {
			logger.Errorf("add key %s error %s", keyString, err)
		}
		bloomfilterMap[bucket].Add(hash)
		count++
	}
	hash.Reset()
	responseSuccess(response, count)
	list = nil
}

func BatchExists(response http.ResponseWriter, request *http.Request) {
	check, params := initParams(response, request)
	if !check {
		return
	}
	if !params["keys"] {
		responseError(response, fmt.Sprintf("keys must set!"))
		return
	}

	var batchList []int8
	SPLIT = rune(',')
	var hash hash2.Hash64
	list := strings.FieldsFunc(keys, stringSpilt)
	for _, keyString := range list {
		if len(keyString) <= 0 {
			continue
		}
		hash = fnv.New64()
		hash.Write([]byte(keyString))
		exists := bloomfilterMap[bucket].Contains(hash)
		if exists {
			batchList = append(batchList, 1)
		} else {
			batchList = append(batchList, 0)
			//exists and add key
			if forceAdd == 1 {
				bloomfilterMap[bucket].Add(hash)
			}
		}
	}
	hash.Reset()
	responseSuccess(response, batchList)
	list = nil
	batchList = nil
}

func Add(response http.ResponseWriter, request *http.Request) {
	check, params := initParams(response, request)
	if !check {
		return
	}
	if !params["key"] {
		responseError(response, fmt.Sprintf("key must set!"))
		return
	}

	hash := fnv.New64()
	hash.Write([]byte(key))
	bloomfilterMap[bucket].Add(hash)
	hash.Reset()
	responseSuccess(response, "success")
}

func Exists(response http.ResponseWriter, request *http.Request) {
	check, params := initParams(response, request)
	if !check {
		return
	}
	if !params["key"] {
		responseError(response, fmt.Sprintf("key must set!"))
		return
	}

	ret := -1
	hash := fnv.New64()
	hash.Write([]byte(key))
	exists := bloomfilterMap[bucket].Contains(hash)
	if exists {
		ret = 1
	} else {
		ret = 0
	}
	hash.Reset()
	responseSuccess(response, ret)
}

func Clear(response http.ResponseWriter, request *http.Request) {
	check, _ := initParams(response, request)
	if !check {
		return
	}
	bloomfilterMap[bucket] = nil
	responseSuccess(response, "success")
}

func Memory(response http.ResponseWriter, request *http.Request) {
	check, _ := initParams(response, request)
	if !check {
		return
	}

	count := bloomfilterMap[bucket].M()
	responseSuccess(response, fmt.Sprintf(" %d mb, %d kb , %d", count/1024/1024, count/1024, count))
}

func KeyCount(response http.ResponseWriter, request *http.Request) {
	check, _ := initParams(response, request)
	if !check {
		return
	}
	filter := *bloomfilterMap[bucket]
	responseSuccess(response, fmt.Sprintf("hit count: %d, key count: %d, memory: %d", filter.N(), filter.K(), filter.M()))
}

func Buckets(response http.ResponseWriter, request *http.Request) {
	check, _ := initParams(response, request)
	if !check {
		return
	}
	var buckets []string
	for k, _ := range bloomfilterMap {
		buckets = append(buckets, k)
	}
	responseSuccess(response, strings.Join(buckets, ","))
}

func dumpEmailZip(response http.ResponseWriter, request *http.Request) {
	limit := getParams(request, "limit")
	if "" == limit {
		limit = "100"
	}
	status := getParams(request, "status")
	if "" == status {
		status = "1"
	}
	depart := int64(3)
	_depart := getParams(request, "depart")
	if "" != _depart {
		depart, _ = strconv.ParseInt(_depart, 10, 34)
	}
	forceCover := getParams(request, "forceCover")
	if "" == forceCover {
		forceCover = "0"
	}
	timePrefix := time.Now().Format("0102")

	//delete old file
	deleteOldFile := func() {
		timeInt, _ := strconv.ParseInt(timePrefix, 10, 64)
		for i := timeInt - 7; i > timeInt-21; i-- {
			oldFileName := fmt.Sprintf("%d_index.csv", i)
			if i < 1000 {
				oldFileName = fmt.Sprintf("0%d_index.csv", i)
			}
			err := os.Remove(oldFileName)
			if nil != err {
				continue
			}
			logger.Infof("delete old file %s", oldFileName)
		}
	}
	deleteOldFile()

	indexFilename := fmt.Sprintf("%s_index.csv", timePrefix)
	if util.FileExist(indexFilename) && "0" == forceCover {
		responseError(response, fmt.Sprintf("file %s has generate success, if want cover it please add params foreCover=1!", indexFilename))
		return
	}

	logger.Infof("start dump un valid email.")
	allData, emailData := dumpUnValidEmailApi(config.mysql, status, limit)
	logger.Infof("finish dump un valid email.")

	endPrefix := ""
	partIndex := 1
	emailCount := len(emailData)
	departPerCount := emailCount / int(depart)
	departList := make([]int, depart)
	for i := departPerCount - 1; i < emailCount; i++ {
		SPLIT = rune('@')
		_emailSplit := strings.FieldsFunc(emailData[i][0], stringSpilt)
		emailName := _emailSplit[0]
		if strings.Contains(emailName, ".") {
			emailName = strings.Replace(emailName, ".", "", len(emailName))
		}

		//logger.Infof("%s %s", emailName, emailData[i][0])
		if "" == endPrefix {
			endPrefix = emailName
		}

		if endPrefix != emailName {
			departList[partIndex-1] = i
			partIndex++
			endPrefix = ""
			i = departPerCount * partIndex

			if partIndex == int(depart) {
				departList[partIndex-1] = emailCount
				break
			}
		}
	}

	logger.Infof("depart list %s, username size %d, email size %d", departList, len(allData), len(emailData))

	compressZip := func() {
		Write(indexFilename, allData)
		csvFiles := make([]*os.File, depart)
		//time.Now().Format("2006-01-02 15:04:05")
		lastEmailIndex := 0
		partFileList := make([]string, depart)
		for index, emailIndex := range departList {
			tmpPartList := emailData[lastEmailIndex:emailIndex]

			logger.Infof("start write part file %d .", index)
			lastEmailIndex = emailIndex
			partFilename := fmt.Sprintf("%s0%d.csv", timePrefix, index+1)
			Write(partFilename, tmpPartList)

			logger.Infof("finish write part file %s .", partFilename)

			file, err := os.Open(partFilename)
			if nil != err {
				logger.Errorf("write file fail %s", partFilename)
			}
			csvFiles[index] = file
			logger.Infof("finish read part file %s .", partFilename)

			partFileList = append(partFileList, partFilename)
		}

		logger.Info("start compress file.")
		zipFileName := fmt.Sprintf("%s(%d).tar.gz", timePrefix, len(emailData))
		os.Remove(zipFileName)
		util.Compress(csvFiles, zipFileName)

		for _, fileName := range partFileList {
			os.Remove(fileName)
		}
		download(zipFileName, response)
		os.Remove(zipFileName)
	}

	compressZip()

	allData = nil
	emailData = nil
	departList = nil
}

func importEmailData(response http.ResponseWriter, request *http.Request) {
	request.ParseMultipartForm(1024 * 1024)
	csvFile, fileHeader, err := request.FormFile("csv")
	if nil == fileHeader {
		responseError(response, "csv file not select.")
		return
	}
	if err != nil {
		responseError(response, err.Error())
		return
	}

	defer csvFile.Close()
	fileName := fileHeader.Filename
	if fileHeader.Size <= 0 {
		responseError(response, "the csv file is empty.")
		return
	}

	copyUploadFile := func() {
		os.Remove(fileName)
		cur, err := os.Create(fileName)
		defer cur.Close()
		if err != nil {
			logger.Error(err)
		}
		io.Copy(cur, csvFile)
	}

	copyUploadFile()

	SPLIT = rune('.')
	filePrefix := strings.FieldsFunc(fileName, stringSpilt)
	timePrefix := filePrefix[0]
	if strings.Contains(timePrefix, "_fail") {
		SPLIT = rune('_')
		timePrefix = strings.FieldsFunc(fileName, stringSpilt)[0]
	}
	if len(timePrefix) == 6 {
		timePrefix = timePrefix[:4]
	}

	indexFile := fmt.Sprintf("%s_index.csv", timePrefix)
	_, err = os.Open(indexFile)
	if nil != err {
		responseError(response, fmt.Sprintf("the index file not exists %d, please check you csv name.", indexFile))
		return
	}

	logger.Infof("index file %s", indexFile)

	logger.Infof("start import email.")
	successCount, failCount, namesRepeatCount, emailCount, emailRepeatCount := importEmailApi(config.mysql, indexFile, fileName)
	logger.Infof("finish import email.")

	os.Remove(fileName)

	jsonMap := make(map[string]string)
	jsonMap["filename"] = fileName
	jsonMap["successCount"] = strconv.Itoa(successCount)
	jsonMap["failCount"] = strconv.Itoa(failCount)
	jsonMap["namesRepeatCount"] = strconv.Itoa(namesRepeatCount)
	jsonMap["emailCount"] = strconv.Itoa(emailCount)
	jsonMap["emailRepeatCount"] = strconv.Itoa(emailRepeatCount)
	responseSuccess(response, jsonMap)
}

func getEmailCount(response http.ResponseWriter, request *http.Request) {
	//check, _ := initParams(response, request)
	//if !check {
	//	return
	//}
	query := request.URL.Query()

	if nil == query["finished"] || nil == query["start"] || nil == query["end"] {
		responseError(response, "finished start(2020-11-01) end(2020-11-01) must set")
		return
	}

	finished, _ := strconv.ParseInt(query["finished"][0], 10, 64)
	start := query["start"][0]
	end := query["end"][0]

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

	count := GetEmailCount(config.mysql, int(finished), _start.Unix(), _end.Unix())
	responseSuccess(response, count)
}
