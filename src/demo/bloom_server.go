package demo

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/steakknife/bloomfilter"
	hash2 "hash"
	"hash/fnv"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	//_maxElements         = uint64(500 * 10000)
	//_probCollide         = 0.0000001

	key      = ""
	keys     = ""
	bucket   = ""
	percent  = 0.00001
	elements = uint64(1 * 10000 * 10000)

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

func Main() {
	runtime.GOMAXPROCS(2)
	seelog.ReplaceLogger(logger)
	defer seelog.Flush()

	http.HandleFunc("/", Index)
	http.HandleFunc("/config", Config)
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
	rawW := gzip.NewWriter(w)
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

					//go func() {
					num = backupFilter(filename, filter)
					bloomfilterMemoryMap[bucket] = num
					//}()
				}
			}
			//runtime.GC()
			saveTimer.Reset(time.Second * 60)
		}
	}()
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
		elements, _ = strconv.ParseUint(query["elements"][0], 10, 64)
		params[paramKey] = true
	} else {
		params[paramKey] = false
	}

	logger.Debugf(">>>request url %s , body %s", request.RequestURI, bodyData)
	response.Header().Add("server", "bloom server")

	if !strings.Contains(request.RequestURI, "config") &&
		!strings.Contains(request.RequestURI, "buckets") {
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

func Config(response http.ResponseWriter, request *http.Request) {
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
