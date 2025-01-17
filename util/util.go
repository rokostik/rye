// util.go
package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/refaktor/rye/env"
)

func PrintHeader() {
	fmt.Println("=-===============-===-===-=============-=")   // Output: -3
	fmt.Println(" _/|\\\\_-~*>%,_  Rejy ZERO  _,%<*~-_//|\\_") // Output: -3
	fmt.Println("=-===============-===-===-=============-=")   // Output: -3
}

func IndexOfAt(s, sep string, n int) int {
	idx := strings.Index(s[n:], sep)
	if idx > -1 {
		idx += n
	}
	return idx
}

func IndexOfSlice(ps *env.ProgramState, slice []env.Object, value env.Object) int {
	for i, v := range slice {
		if v.Equal(value) {
			return i
		}
	}
	return -1 // not found
}

func IsTruthy(o env.Object) bool {
	switch oo := o.(type) {
	case env.Integer:
		return oo.Value > 0
	case env.String:
		return len(oo.Value) > 0
	default:
		return false
	}
}

func Dict2Context(ps *env.ProgramState, s1 env.Dict) env.RyeCtx {
	ctx := env.NewEnv(ps.Ctx)
	for k, v := range s1.Data {
		word := ps.Idx.IndexWord(k)
		switch v1 := v.(type) {
		case env.Integer:
			ctx.Set(word, v1)
		case env.String:
			ctx.Set(word, v1)
		case string:
			ctx.Set(word, *env.NewString(v1))
		}
	}
	return *ctx
}

func StringToFieldsWithQuoted(str string, sepa string, quote string) env.Block {
	quoted := false
	spl := strings.FieldsFunc(str, func(r rune) bool {
		if string(r) == quote {
			quoted = !quoted
		}
		return !quoted && string(r) == sepa
	})
	lst := make([]env.Object, len(spl))
	re := regexp.MustCompile("[0-9]+")
	for i := 0; i < len(spl); i++ {
		//fmt.Println(spl[i])
		// TODO -- detect numbers and turn them to integers or floats, later we can also detect other types of values
		// val, _ := loader.LoadString(spl[i], false)
		// numeric, _ := regexp.MatchString("[0-9]+", spl[i])
		numeric := re.MatchString(spl[i])
		// fmt.Println(numeric)
		pass := false
		if numeric {
			num, err := strconv.Atoi(spl[i])
			if err == nil {
				lst[i] = *env.NewInteger(int64(num))
				pass = true
			}
		}
		if !pass {
			clean := regexp.MustCompile(`^"(.*)"$`).ReplaceAllString(spl[i], `$1`)
			val := *env.NewString(clean)
			lst[i] = val
		}
	}
	return *env.NewBlock(*env.NewTSeries(lst))
}

func FormatJson(val env.Object, e env.Idxs) string {
	// TODO -- this is currently made just for block of strings and integers
	var r strings.Builder
	switch b := val.(type) {
	case env.Block:
		r.WriteString("[ ")
		for i := 0; i < b.Series.Len(); i += 1 {
			if b.Series.Get(i) != nil {
				if i > 0 {
					r.WriteString(", ")
				}
				r.WriteString(b.Series.Get(i).Print(e))
			}
		}
		r.WriteString(" ]")
	}
	return r.String()
}

func FormatCsv(val env.Object, e env.Idxs) string {
	// TODO -- this is currently made just for block of strings and integers
	var r strings.Builder
	switch b := val.(type) {
	case env.Block:
		for i := 0; i < b.Series.Len(); i += 1 {
			if b.Series.Get(i) != nil {
				if i > 0 {
					r.WriteString(",")
				}
				r.WriteString(b.Series.Get(i).Print(e))
			}
		}
	}
	return r.String()
}

func FormatSsv(val env.Object, e env.Idxs) string {
	// TODO -- this is currently made just for block of strings and integers
	var r strings.Builder
	switch b := val.(type) {
	case env.Block:
		for i := 0; i < b.Series.Len(); i += 1 {
			if b.Series.Get(i) != nil {
				if i > 0 {
					r.WriteString(" ")
				}
				o := b.Series.Get(i)
				switch ob := o.(type) {
				case env.String:
					r.WriteString(ob.Value)
				default:
					r.WriteString(ob.Print(e))
				}
			}
		}
	}
	return r.String()
}

func SplitEveryString(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string = make([]string, 0, (len(s)-1)/chunkSize+1)
	currentLen := 0
	currentStart := 0
	for i := range s {
		if currentLen == chunkSize {
			chunks = append(chunks, s[currentStart:i])
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, s[currentStart:])
	return chunks
}

func SplitEveryList(s []env.Object, chunkSize int) [][]env.Object {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return append(make([][]env.Object, 0), s)
	}
	var chunks [][]env.Object = make([][]env.Object, 0, (len(s)-1)/chunkSize+1)
	var chunk []env.Object = make([]env.Object, 0, chunkSize)
	currentLen := 0
	// currentStart := 0
	for i := range s {
		chunk = append(chunk, s[i])
		currentLen++
		if currentLen == chunkSize {
			chunks = append(chunks, chunk)
			currentLen = 0
			//		currentStart = i + 1
			chunk = make([]env.Object, 0, chunkSize)
		}
	}
	if len(chunk) > 0 {
		chunks = append(chunks, chunk)
	}
	return chunks
}

func IntersectStrings(a string, b string) string {
	res := ""
	for _, ch := range a {
		if strings.Contains(b, string(ch)) && !strings.Contains(res, string(ch)) {
			res = res + string(ch)
		}
	}

	return res
}

func IntersectLists(ps *env.ProgramState, a []env.Object, b []env.Object) []env.Object {
	set := make([]env.Object, 0)
	for _, v := range a {
		if ContainsVal(ps, b, v) && !ContainsVal(ps, set, v) {
			set = append(set, v)
		}
	}
	return set
}

func UnionOfLists(ps *env.ProgramState, a []env.Object, b []env.Object) []env.Object {
	elementMap := make(map[env.Object]bool)
	for _, element := range a {
		elementMap[element] = true
	}
	for _, element := range b {
		elementMap[element] = true
	}
	mergedSlice := make([]env.Object, 0)
	for element := range elementMap {
		mergedSlice = append(mergedSlice, element)
	}
	return mergedSlice
}

func SplitMulti(s string, seps string) []string {
	splitter := func(r rune) bool {
		return strings.ContainsRune(seps, r)
	}
	return strings.FieldsFunc(s, splitter)
}

func ContainsVal(ps *env.ProgramState, b []env.Object, val env.Object) bool {
	for _, a := range b {
		if a.Equal(val) {
			return true
		}
	}
	return false
}

func RemoveDuplicate(ps *env.ProgramState, slice []env.Object) []env.Object {
	allKeys := make(map[env.Object]bool)
	list := []env.Object{}
	for _, item := range slice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

/*
func RemoveDuplicate[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
*/

/* func Transpose(slice []env.Object) []env.Object {
	yl := len(slice)
	var xl int
	switch blk := slice[0].(type) {
	case env.Block:
		xl = len(blk.Series.S)
	}
	if xl == 0 {
		return nil
	}
	result := make([]env.Object, xl)
	//for i := range result { // TODOD .... finish this or comment it out next time ... left as it is
	// result[i] = env.NewBlock(env.NewTSeries()) // make([]env.Object, yl)
	//}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = slice[j][i]
		}
	}
	return result
        }*/
