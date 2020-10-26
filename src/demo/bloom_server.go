package demo

import (
	"encoding/json"
	"fmt"
	"github.com/steakknife/bloomfilter"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
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
	percent  = 0.0000001
	elements = uint64(1000 * 10000)

	backDumpPrefix      = "BF"
	bloomfilterCountMap = map[string]uint64{}
	bloomfilterMap      = map[string]*bloomfilter.Filter{}

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
	http.HandleFunc("/config", Config)
	http.HandleFunc("/add", Add)
	http.HandleFunc("/batch_add", BatchAdd)
	http.HandleFunc("/exists", Exists)
	http.HandleFunc("/batch_exists", BatchExists)
	http.HandleFunc("/clear", Clear)
	http.HandleFunc("/save", Save)
	http.HandleFunc("/load", Load)
	http.HandleFunc("/memory", Memory)
	http.HandleFunc("/key_count", KeyCount)

	startSaveTask()
	go loadOldData()
	log.Println(serverName + " start success " + strconv.Itoa(serverPort))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(serverPort), nil))
}

func loadOldData() {
	pwd, _ := os.Getwd()
	fileInfoList, err := ioutil.ReadDir(pwd)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(fileInfoList))
	for i := range fileInfoList {
		filename := fileInfoList[i].Name()
		if strings.HasPrefix(filename, backDumpPrefix) {
			r, err := os.Open(filename)
			if err != nil {
				log.Fatalf("laod file error %s ", err)
			}
			defer func() {
				err = r.Close()
			}()

			SPLIT = rune('.')
			list := strings.FieldsFunc(filename, stringSpilt)
			if len(list) != 3 {
				log.Printf("filename %s jump", filename)
				continue
			}

			_bucket := list[1]
			bloomfilterMap[_bucket], err = bloomfilter.NewOptimal(elements, percent)
			_, err = bloomfilterMap[_bucket].ReadFrom(r)
			if nil != err {
				log.Fatalf("load file data error %s", err)
			}
			log.Printf("load file %s, bucket %s , count %d, memory %d", filename, _bucket, bloomfilterMap[_bucket].N(), bloomfilterMap[_bucket].M())
		}
	}

}

func startSaveTask() {
	go func() {
		for range saveTimer.C {
			go func() {
				for bucket, filter := range bloomfilterMap {
					KeyCount := filter.N()
					filename := fmt.Sprintf("%s.%s.bucket", backDumpPrefix, bucket)
					if bloomfilterCountMap[bucket] != KeyCount {
						filter.WriteFile(filename)
						bloomfilterCountMap[bucket] = KeyCount
						log.Printf(" save bucket %s", filename)
					}
				}
			}()
			saveTimer.Reset(time.Second * 2)
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
		log.Fatal(err)
	}
	response.WriteHeader(result.Code)
	response.Write(json)

	log.Printf("<<<response %s", json)
}

func responseError(response http.ResponseWriter, msg string) {
	result := getErrorResult(msg)
	json, _ := json.Marshal(&result)
	response.WriteHeader(result.Code)
	response.Write(json)

	log.Printf("<<<response %s", json)
}

func initParams(response http.ResponseWriter, request *http.Request) (bool, map[string]bool) {
	key = ""
	keys = ""
	bucket = ""
	percent = 0.0000001
	elements = uint64(500 * 10000)
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

	log.Printf(">>>request url %s , body %s", request.RequestURI, bodyData)
	response.Header().Add("server", "bloom server")
	if "" == bucket || len(bucket) <= 0 {
		responseError(response, fmt.Sprintf("params bucket must set!"))
		return false, params
	}
	if bloomfilterMap[bucket] == nil && !strings.Contains(request.RequestURI, "config") {
		responseError(response, fmt.Sprintf("bucket %s not exists, please use /config init first", bucket))
		return false, params
	}
	return true, params
}

func Config(response http.ResponseWriter, request *http.Request) {
	check, params := initParams(response, request)
	if !check {
		return
	}
	if bloomfilterMap[bucket] != nil {
		responseError(response, fmt.Sprintf("bucket %s exists, please reset first", bucket))
		return
	}
	if !params["percent"] || !params["elements"] {
		responseError(response, fmt.Sprintf("percent(%f) and elements(%d) must set!", percent, elements))
		return
	}
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

	filename := fmt.Sprintf("%s.%s.bucket", backDumpPrefix, bucket)
	bloomfilterMap[bucket].WriteFile(filename)
	responseSuccess(response, fmt.Sprintf("save to file %s success, count %d", filename, bloomfilterMap[bucket].N()))
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
	list := strings.FieldsFunc(keys, stringSpilt)
	for _, keyString := range list {
		if len(keyString) <= 0 {
			continue
		}
		hash := fnv.New64()
		hash.Write([]byte(keyString))
		bloomfilterMap[bucket].Add(hash)
		count++
	}
	responseSuccess(response, count)
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
	responseSuccess(response, "success")
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
	list := strings.FieldsFunc(keys, stringSpilt)
	for _, keyString := range list {
		if len(keyString) <= 0 {
			continue
		}
		hash := fnv.New64()
		hash.Write([]byte(keyString))
		exists := bloomfilterMap[bucket].Contains(hash)
		if exists {
			batchList = append(batchList, 1)
		} else {
			batchList = append(batchList, 0)
		}
	}
	responseSuccess(response, batchList)
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
	responseSuccess(response, bloomfilterMap[bucket].N())
}
