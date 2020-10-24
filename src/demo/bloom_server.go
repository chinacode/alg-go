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
)

var (
	//_maxElements         = uint64(500 * 10000)
	//_probCollide         = 0.0000001

	key            = ""
	keys           = ""
	bucket         = ""
	percent        = 0.0000001
	elements       = uint64(500 * 10000)
	bloomfilterMap = map[string]*bloomfilter.Filter{}

	SERVER_PORT = 9002
	SERVER_NAME = "bloomfilter server "
)

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Main() {
	SPLIT = rune(',')

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

	log.Println(SERVER_NAME + " start success " + strconv.Itoa(SERVER_PORT))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(SERVER_PORT), nil))
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
	} else if nil != query[paramKey] || nil != request.Body {
		keys = query[paramKey][0]
		params[paramKey] = true
	} else {
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
	filename := fmt.Sprintf("BF.%s.bucket", bucket)
	r, err := os.Open(filename)
	if err != nil {
		responseError(response, fmt.Sprintf("back file %s not exists", filename))
		return
	}
	defer func() {
		err = r.Close()
	}()

	bloomfilterMap[bucket].ReadFrom(r)
	responseSuccess(response, fmt.Sprintf("load from file %s success", filename))
}

func Save(response http.ResponseWriter, request *http.Request) {
	check, _ := initParams(response, request)
	if !check {
		return
	}

	filename := fmt.Sprintf("BF.%s.bucket", bucket)
	bloomfilterMap[bucket].WriteFile(filename)
	responseSuccess(response, fmt.Sprintf("save to file %s success", filename))
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
	responseSuccess(response, bloomfilterMap[bucket].K())
}
