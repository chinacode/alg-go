package demo

import (
	"bufio"
	"fmt"
	"github.com/shopspring/decimal"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

// A data structure to hold a key/value pair.
type Pair struct {
	Key   string
	Value int
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// A function to turn a map into a PairList, then sort and return it.
func sortMapByValue(m map[string]int) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
	}
	sort.Sort(p)
	return p
}

func OssEmailEndPrefix() map[string]int {
	endPrefix := make(map[string]int)
	args := os.Args

	if len(args) < 5 {
		fmt.Printf("params: （path split index limit）")
		return endPrefix
	}
	path := args[1]
	split := args[2]
	index, err := strconv.Atoi(args[3])
	if nil != err {
		fmt.Printf("index error %s", err)
		return endPrefix
	}
	limit, err2 := strconv.Atoi(args[4])
	if nil != err2 {
		fmt.Printf("index error %s", err)
		return endPrefix
	}
	total := 0
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		filename := f.Name()
		if filename[0] != 'x' {
			continue
		}
		fileOperate, err := os.Open(filename)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		//println(filename)
		defer fileOperate.Close()
		br := bufio.NewReader(fileOperate)
		for {
			line, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			lineStr := string(line)
			if !strings.Contains(lineStr, "@") {
				continue
			}
			if len(lineStr) <= 4 {
				continue
			}
			//sql
			end := ""
			if index == -1 {
				end = strings.Split(lineStr, "@")[1]
				end = strings.ToLower(strings.Split(end, "'")[0])
			} else {
				ends := strings.Split(lineStr, split)
				if len(ends) < index+1 {
					continue
				}
				end = ends[index]
				if !strings.Contains(end, "@") {
					continue
				}
				end = strings.ToLower(strings.Split(end, "@")[1])
			}
			if len(end) <= 0 {
				continue
			}
			if _, ok := endPrefix[end]; ok {
				endPrefix[end] = endPrefix[end] + 1
			} else {
				endPrefix[end] = 1
			}
			total++
		}
	}

	jumpCount := 0
	listI := 0
	var keys []int
	endPrefixMap := make(map[int]string)
	for key, value := range endPrefix {
		if value < limit {
			jumpCount++
			continue
		}
		newValueKey := value * 100
		if _, ok := endPrefixMap[newValueKey]; ok {
			newValueKey = newValueKey + listI
			listI++
		}

		endPrefixMap[newValueKey] = key
		keys = append(keys, newValueKey)
	}
	sort.Ints(keys)

	//endPrefix, _ = SortMap(endPrefix)
	for i := len(keys) - 1; i >= 0; i-- {
		valueKey := keys[i]
		value := valueKey / 100
		//if value < 10000 {
		//	continue
		//}
		key := endPrefixMap[valueKey]
		percent := decimal.NewFromFloat(float64(value*100)).DivRound(decimal.NewFromFloat(float64(total)), 4)
		newPercent, _ := percent.Float64()

		fmt.Printf("%s：%d，(%.2f)\n", key, value, newPercent)
	}

	fmt.Printf("total count %d，jump count %d", total, jumpCount)
	return endPrefix
}
