package demo

import (
	"encoding/json"
	"fmt"
	"github.com/steakknife/bloomfilter"
	"hash/fnv"
	"log"
	"net/http"
	"os"
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
	code int8
	msg  string
	data string
}

func Main() {
	SPLIT = rune(',')
	http.HandleFunc("/config", Config)
	http.HandleFunc("/add", Add)
	http.HandleFunc("/batch_add", BatchAdd)
	http.HandleFunc("/exists", Exists)
	http.HandleFunc("/clear", Clear)
	http.HandleFunc("/save", Save)
	http.HandleFunc("/load", Load)
	http.HandleFunc("/memory", Memory)
	http.HandleFunc("/key_count", KeyCount)

	log.Println(SERVER_NAME + " start success " + strconv.Itoa(SERVER_PORT))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(SERVER_PORT), nil))
}

func getErrorResult(msg string) Result {
	return Result{code: http.StatusInternalServerError, msg: msg, data: nil}
}

func getResult(data map[string]string) Result {
	jsonData, _ := json.Marshal(data)
	return Result{code: http.StatusOK, msg: "success", data: string(jsonData)}
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
	if nil != query[paramKey] {
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

	if "" == bucket || len(bucket) <= 0 {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(fmt.Sprintf("params bucket must set!")))
		return false, params
	}
	if bloomfilterMap[bucket] == nil && !strings.Contains(request.RequestURI, "config") {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(fmt.Sprintf("bucket %s not exists, please use /config init first", bucket)))
		return false, params
	}
	return true, params
}

func initBucket() {
	if len(bucket) > 0 {
		bloomfilterMap[bucket], err = bloomfilter.NewOptimal(elements, percent)
	}
}

func BatchAdd(response http.ResponseWriter, request *http.Request) {
	check, params := initParams(response, request)
	if !check {
		return
	}
	if !params["keys"] {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(fmt.Sprintf("keys must set!")))
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
	response.Write([]byte(strconv.Itoa(count)))
}

func Add(response http.ResponseWriter, request *http.Request) {
	check, params := initParams(response, request)
	if !check {
		return
	}
	if !params["key"] {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(fmt.Sprintf("key must set!")))
		return
	}

	hash := fnv.New64()
	hash.Write([]byte(key))
	bloomfilterMap[bucket].Add(hash)
	response.Write([]byte("success"))
}

func Exists(response http.ResponseWriter, request *http.Request) {
	check, params := initParams(response, request)
	if !check {
		return
	}
	if !params["key"] {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(fmt.Sprintf("key must set!")))
		return
	}

	hash := fnv.New64()
	hash.Write([]byte(key))
	exists := bloomfilterMap[bucket].Contains(hash)
	if exists {
		response.Write([]byte("1"))
	} else {
		response.Write([]byte("0"))
	}
}

func Clear(response http.ResponseWriter, request *http.Request) {
	check, _ := initParams(response, request)
	if !check {
		return
	}
	bloomfilterMap[bucket] = nil
	response.Write([]byte("success"))
}

func Memory(response http.ResponseWriter, request *http.Request) {
	check, _ := initParams(response, request)
	if !check {
		return
	}

	response.Write([]byte(fmt.Sprintf("%d", bloomfilterMap[bucket].M())))
}

func KeyCount(response http.ResponseWriter, request *http.Request) {
	check, _ := initParams(response, request)
	if !check {
		return
	}
	response.Write([]byte(fmt.Sprintf("%d", bloomfilterMap[bucket].K())))
}

func Config(response http.ResponseWriter, request *http.Request) {
	check, params := initParams(response, request)
	if !check {
		return
	}
	if bloomfilterMap[bucket] != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(fmt.Sprintf("bucket %s exists, please reset first", bucket)))
		return
	}
	if !params["percent"] || !params["elements"] {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(fmt.Sprintf("percent(%f) and elements(%d) must set!", percent, elements)))
		return
	}
	bloomfilterMap[bucket], err = bloomfilter.NewOptimal(elements, percent)
	response.Write([]byte("success"))
}

func Load(response http.ResponseWriter, request *http.Request) {
	check, _ := initParams(response, request)
	if !check {
		return
	}
	filename := fmt.Sprintf("BF.%s.bucket", bucket)

	r, err := os.Open(filename)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(fmt.Sprintf("back file %s not exists", filename)))
		return
	}
	defer func() {
		err = r.Close()
	}()

	bloomfilterMap[bucket].ReadFrom(r)
	response.Write([]byte(fmt.Sprintf("load from file %s success", filename)))
}

func Save(response http.ResponseWriter, request *http.Request) {
	check, _ := initParams(response, request)
	if !check {
		return
	}

	filename := fmt.Sprintf("BF.%s.bucket", bucket)
	bloomfilterMap[bucket].WriteFile(filename)
	response.Write([]byte(fmt.Sprintf("save to file %s success", filename)))
}
