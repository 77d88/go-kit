package xarray

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"

	"github.com/77d88/go-kit/basic/xcore"
	"golang.org/x/exp/constraints"
)

// Create a static variable to store the hash table.
// This variable has the same lifetime as the entire program and can be shared by functions that are called more than once.
var (
	memoryHashMap     = make(map[string]map[any]int)
	memoryHashCounter = make(map[string]int)
)

// Contain 检查目标值是否在切片中存在。
func Contain[T comparable](slice []T, target T) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}

	return false
}

// ContainAll 检查目标切片是否包含给定的所有元素。
func ContainAll[T comparable](slice []T, target []T) bool {
	for _, item := range target {
		if !Contain(slice, item) {
			return false
		}
	}
	return true
}

// ContainAny 检查目标切片是否包含给定的任意元素。
func ContainAny[T comparable](slice []T, target []T) bool {
	for _, item := range target {
		if Contain(slice, item) {
			return true
		}
	}
	return false
}

// ContainBy 根据提供的谓词函数判断元素是否存在于切片中。
func ContainBy[T any](slice []T, predicate func(item T) bool) bool {
	for _, item := range slice {
		if predicate(item) {
			return true
		}
	}

	return false
}

// ContainSubSlice 检查主切片是否包含给定的子切片。
func ContainSubSlice[T comparable](slice, subSlice []T) bool {
	if len(subSlice) == 0 {
		return true
	}
	if len(slice) == 0 {
		return false
	}

	elementMap := make(map[T]struct{}, len(slice))
	for _, item := range slice {
		elementMap[item] = struct{}{}
	}

	for _, item := range subSlice {
		if _, ok := elementMap[item]; !ok {
			return false
		}
	}
	return true
}

// Chunk 将切片元素分割成指定大小的组，返回二维切片。
func Chunk[T any](slice []T, size int) [][]T {
	result := [][]T{}

	if len(slice) == 0 || size <= 0 {
		return result
	}

	for _, item := range slice {
		l := len(result)
		if l == 0 || len(result[l-1]) == size {
			result = append(result, []T{})
			l++
		}

		result[l-1] = append(result[l-1], item)
	}

	return result
}

// Compact 移除切片中的所有'假值'元素（如false, nil, 0, ""），并返回新切片。
func Compact[T comparable](slice []T) []T {
	var zero T

	result := make([]T, 0, len(slice))

	for _, v := range slice {
		if v != zero {
			result = append(result, v)
		}
	}
	return result[:len(result):len(result)]
}

// Difference 差集 返回存在于第一个切片中但不在第二个比较切片中的元素组成的切片。
func Difference[T comparable](slice, comparedSlice []T) []T {
	var result []T

	for _, v := range slice {
		if !Contain(comparedSlice, v) {
			result = append(result, v)
		}
	}

	return result
}

// DifferenceBy 根据提供的迭代函数处理元素后，返回两个切片的差集。
func DifferenceBy[T comparable](slice []T, comparedSlice []T, iteratee func(index int, item T) T) []T {
	orginSliceAfterMap := Map(slice, iteratee)
	comparedSliceAfterMap := Map(comparedSlice, iteratee)

	result := make([]T, 0)
	for i, v := range orginSliceAfterMap {
		if !Contain(comparedSliceAfterMap, v) {
			result = append(result, slice[i])
		}
	}

	return result
}

// DifferenceWith 根据提供的比较器函数比较元素，返回差集。
func DifferenceWith[T any](slice []T, comparedSlice []T, comparator func(item1, item2 T) bool) []T {
	result := make([]T, 0)

	getIndex := func(arr []T, item T, comparison func(v1, v2 T) bool) int {
		index := -1
		for i, v := range arr {
			if comparison(item, v) {
				index = i
				break
			}
		}
		return index
	}

	for i, v := range slice {
		index := getIndex(comparedSlice, v, comparator)
		if index == -1 {
			result = append(result, slice[i])
		}
	}

	return result
}

// Equal 判断两个切片是否相等：长度相同且所有元素顺序及值均相等。
func Equal[T comparable](slice1, slice2 []T) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}

	return true
}

// EqualWith 使用比较器函数判断两个切片是否相等。
func EqualWith[T, U any](slice1 []T, slice2 []U, comparator func(T, U) bool) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	for i, v := range slice1 {
		if !comparator(v, slice2[i]) {
			return false
		}
	}

	return true
}

// Every 如果切片中所有元素都满足谓词函数，则返回true。
func Every[T any](slice []T, predicate func(index int, item T) bool) bool {
	for i, v := range slice {
		if !predicate(i, v) {
			return false
		}
	}

	return true
}

// None 如果切片中所有元素都不满足条件，则返回true。
func None[T any](slice []T, predicate func(index int, item T) bool) bool {
	l := 0
	for i, v := range slice {
		if !predicate(i, v) {
			l++
		}
	}

	return l == len(slice)
}

// Some 如果切片中有任何元素满足谓词函数，则返回true。
func Some[T any](slice []T, predicate func(index int, item T) bool) bool {
	for i, v := range slice {
		if predicate(i, v) {
			return true
		}
	}

	return false
}

// Filter 过滤切片中的元素，仅保留满足谓词函数的元素。
func Filter[T any](slice []T, predicate func(index int, item T) bool) []T {
	result := make([]T, 0, len(slice))

	for i, v := range slice {
		if predicate(i, v) {
			result = append(result, v)
		}
	}

	return result
}

// Count 统计给定元素在切片中出现的次数。
func Count[T comparable](slice []T, item T) int {
	count := 0

	for _, v := range slice {
		if item == v {
			count++
		}
	}

	return count
}

// CountBy 根据谓词函数统计切片中匹配元素的数量。
func CountBy[T any](slice []T, predicate func(index int, item T) bool) int {
	count := 0

	for i, v := range slice {
		if predicate(i, v) {
			count++
		}
	}

	return count
}

// GroupBy 根据条件函数将切片元素分为两组并返回。
func GroupBy[T any](slice []T, groupFn func(index int, item T) bool) ([]T, []T) {
	if len(slice) == 0 {
		return make([]T, 0), make([]T, 0)
	}

	groupB := make([]T, 0)
	groupA := make([]T, 0)

	for i, v := range slice {
		ok := groupFn(i, v)
		if ok {
			groupA = append(groupA, v)
		} else {
			groupB = append(groupB, v)
		}
	}

	return groupA, groupB
}

// GroupWith 根据迭代函数的结果对切片元素进行分组，返回一个映射表。

func GroupWith[T any, U comparable](slice []T, iteratee func(item T) U) map[U][]T {
	result := make(map[U][]T)

	for _, v := range slice {
		key := iteratee(v)
		if _, ok := result[key]; !ok {
			result[key] = []T{}
		}
		result[key] = append(result[key], v)
	}

	return result
}

// Find 查找切片中第一个满足谓词函数的元素，直接返回元素和布尔值，无需解引用。
func Find[T any](slice []T, predicate func(index int, item T) bool) (v T, ok bool) {
	index := -1

	for i, v := range slice {
		if predicate(i, v) {
			index = i
			break
		}
	}

	if index == -1 {
		return v, false
	}

	return slice[index], true
}

// FindLast 从切片末尾开始查找最后一个满足谓词函数的元素，直接返回元素和布尔值，无需解引用。
func FindLast[T any](slice []T, predicate func(index int, item T) bool) (v T, ok bool) {
	index := -1

	for i := len(slice) - 1; i >= 0; i-- {
		if predicate(i, slice[i]) {
			index = i
			break
		}
	}

	if index == -1 {
		return v, false
	}

	return slice[index], true
}

// Flatten 展开一层切片。
func Flatten(slice any) any {
	sv := sliceValue(slice)

	var result reflect.Value
	if sv.Type().Elem().Kind() == reflect.Interface {
		result = reflect.MakeSlice(reflect.TypeOf([]interface{}{}), 0, sv.Len())
	} else if sv.Type().Elem().Kind() == reflect.Slice {
		result = reflect.MakeSlice(sv.Type().Elem(), 0, sv.Len())
	} else {
		return result
	}

	for i := 0; i < sv.Len(); i++ {
		item := reflect.ValueOf(sv.Index(i).Interface())
		if item.Kind() == reflect.Slice {
			for j := 0; j < item.Len(); j++ {
				result = reflect.Append(result, item.Index(j))
			}
		} else {
			result = reflect.Append(result, item)
		}
	}

	return result.Interface()
}

// FlattenDeep 递归展开切片。
func FlattenDeep(slice any) any {
	sv := sliceValue(slice)
	st := sliceElemType(sv.Type())

	tmp := reflect.MakeSlice(reflect.SliceOf(st), 0, 0)

	result := flattenRecursive(sv, tmp)

	return result.Interface()
}

func flattenRecursive(value reflect.Value, result reflect.Value) reflect.Value {
	for i := 0; i < value.Len(); i++ {
		item := value.Index(i)
		kind := item.Kind()

		if kind == reflect.Slice {
			result = flattenRecursive(item, result)
		} else {
			result = reflect.Append(result, item)
		}
	}

	return result
}

// ForEach 遍历切片并对每个元素应用函数。
func ForEach[T any](slice []T, iteratee func(index int, item T)) {
	for i := 0; i < len(slice); i++ {
		iteratee(i, slice[i])
	}
}

// ForEachWithBreak 遍历切片并对每个元素应用函数，可中断循环。
func ForEachWithBreak[T any](slice []T, iteratee func(index int, item T) bool) {
	for i := 0; i < len(slice); i++ {
		if !iteratee(i, slice[i]) {
			break
		}
	}
}

// Map 通过迭代函数映射切片中的每个元素到新类型。
func Map[T any, U any](slice []T, iteratee func(index int, item T) U) []U {
	result := make([]U, len(slice), cap(slice))

	for i := 0; i < len(slice); i++ {
		result[i] = iteratee(i, slice[i])
	}

	return result
}

// MapBy 对切片元素进行过滤和映射，返回新切片。
func MapBy[T any, U any](slice []T, iteratee func(index int, item T) (U, bool)) []U {
	result := make([]U, 0)
	for i, v := range slice {
		if a, ok := iteratee(i, v); ok {
			result = append(result, a)
		}
	}
	return result
}

func MapByErr[T any, U any](slice []T, iteratee func(index int, item T) (U, error)) ([]U, error) {
	result := make([]U, 0, len(slice))
	var errs error
	for i, v := range slice {
		if a, err := iteratee(i, v); err == nil {
			result = append(result, a)
		} else {
			errs = errors.Join(errs, err)
		}
	}
	return result, errs
}

// MapUnique 对切片元素进行过滤和映射，返回新切片，并确保结果中的元素是唯一的。
func MapUnique[T any, V comparable](slice []T, extractor func(index int, item T) (V, bool)) []V {
	seen := make(map[V]struct{})
	result := make([]V, 0, len(slice))
	for i, item := range slice {
		if value, ok := extractor(i, item); ok {
			if _, exists := seen[value]; !exists {
				seen[value] = struct{}{}
				result = append(result, value)
			}
		}
	}
	return result
}

// MapUniqueErr 对切片元素进行过滤和映射，返回新切片，并确保结果中的元素是唯一的。如果有错误，则返回错误。
func MapUniqueErr[T any, V comparable](slice []T, extractor func(index int, item T) (V, error)) ([]V, error) {
	seen := make(map[V]struct{})
	result := make([]V, 0, len(slice))
	var errs error
	for i, item := range slice {
		if value, err := extractor(i, item); err == nil {
			if _, exists := seen[value]; !exists {
				seen[value] = struct{}{}
				result = append(result, value)
			}
		} else {
			errs = errors.Join(errs, err)
		}
	}
	return result, errs
}

// FlatMap 扁平化映射切片，即映射后结果再扁平化。
func FlatMap[T any, U any](slice []T, iteratee func(index int, item T) []U) []U {
	result := make([]U, 0, len(slice))

	for i, v := range slice {
		result = append(result, iteratee(i, v)...)
	}

	return result
}

// Reduce 通过迭代函数累积切片元素生成单一值。
func Reduce[T any](slice []T, iteratee func(index int, item1, item2 T) T, initial T) T {
	accumulator := initial

	for i, v := range slice {
		accumulator = iteratee(i, v, accumulator)
	}

	return accumulator
}

// ReduceBy 累积切片元素生成单一值，带初始累积器。
func ReduceBy[T any, U any](slice []T, initial U, reducer func(index int, item T, agg U) U) U {
	accumulator := initial

	for i, v := range slice {
		accumulator = reducer(i, v, accumulator)
	}

	return accumulator
}

// ReduceRight 类似ReduceBy，但反向遍历切片。
func ReduceRight[T any, U any](slice []T, initial U, reducer func(index int, item T, agg U) U) U {
	accumulator := initial

	for i := len(slice) - 1; i >= 0; i-- {
		accumulator = reducer(i, slice[i], accumulator)
	}

	return accumulator
}

// Replace 替换切片中前n个指定元素。
func Replace[T comparable](slice []T, old T, new T, n int) []T {
	result := make([]T, len(slice))
	copy(result, slice)

	for i := range result {
		if result[i] == old && n != 0 {
			result[i] = new
			n--
		}
	}

	return result
}

// ReplaceAll 替换切片中所有指定元素。
func ReplaceAll[T comparable](slice []T, old T, new T) []T {
	return Replace(slice, old, new, -1)
}

// Repeat 生成包含重复元素的新切片。
func Repeat[T any](item T, n int) []T {
	result := make([]T, n)

	for i := range result {
		result[i] = item
	}

	return result
}

// Delete 删除切片中指定元素
func Delete[T comparable](slice []T, element T) []T {
	result := make([]T, 0, len(slice)-1)

	for _, v := range slice {
		if v != element {
			result = append(result, v)
		}
	}

	return result
}

// DeleteAt 删除切片中指定索引的元素
func DeleteAt[T any](slice []T, index int) []T {
	if index >= len(slice) {
		index = len(slice) - 1
	}

	result := make([]T, len(slice)-1)
	copy(result, slice[:index])
	copy(result[index:], slice[index+1:])

	return result
}

// DeleteRange 删除切片中指定范围的元素。
func DeleteRange[T any](slice []T, start, end int) []T {
	result := make([]T, 0, len(slice)-(end-start))

	for i := 0; i < start; i++ {
		result = append(result, slice[i])
	}

	for i := end; i < len(slice); i++ {
		result = append(result, slice[i])
	}

	return result
}

// Drop 从切片开头删除n个元素。
func Drop[T any](slice []T, n int) []T {
	size := len(slice)

	if size <= n {
		return []T{}
	}

	if n <= 0 {
		return slice
	}

	result := make([]T, 0, size-n)

	return append(result, slice[n:]...)
}

// DropRight 从切片末尾删除n个元素
func DropRight[T any](slice []T, n int) []T {
	size := len(slice)

	if size <= n {
		return []T{}
	}

	if n <= 0 {
		return slice
	}

	result := make([]T, 0, size-n)

	return append(result, slice[:size-n]...)
}

// DropWhile 根据条件函数从切片开头删除元素。
func DropWhile[T any](slice []T, predicate func(item T) bool) []T {
	i := 0

	for ; i < len(slice); i++ {
		if !predicate(slice[i]) {
			break
		}
	}

	result := make([]T, 0, len(slice)-i)

	return append(result, slice[i:]...)
}

// DropRightWhile 从切片末尾开始，根据谓词函数删除元素，直到函数返回false。
func DropRightWhile[T any](slice []T, predicate func(item T) bool) []T {
	i := len(slice) - 1

	for ; i >= 0; i-- {
		if !predicate(slice[i]) {
			break
		}
	}

	result := make([]T, 0, i+1)

	return append(result, slice[:i+1]...)
}

// InsertAt 在指定索引处将值或另一个切片插入到切片中。
func InsertAt[T any](slice []T, index int, value any) []T {
	size := len(slice)

	if index < 0 || index > size {
		return slice
	}

	if v, ok := value.(T); ok {
		slice = append(slice[:index], append([]T{v}, slice[index:]...)...)
		return slice
	}

	if v, ok := value.([]T); ok {
		slice = append(slice[:index], append(v, slice[index:]...)...)
		return slice
	}

	return slice
}

// UpdateAt 更新切片中指定索引位置的元素。
func UpdateAt[T any](slice []T, index int, value T) []T {
	size := len(slice)

	if index < 0 || index >= size {
		return slice
	}
	slice = append(slice[:index], append([]T{value}, slice[index+1:]...)...)

	return slice
}

// Unique 移除切片中的重复元素。
func Unique[T comparable](slice []T) []T {
	result := []T{}
	exists := map[T]bool{}
	for _, t := range slice {
		if exists[t] {
			continue
		}
		exists[t] = true
		result = append(result, t)
	}
	return result
}

// UniqueBy 根据键提取函数移除切片中的重复元素
func UniqueBy[T any, K comparable](slice []T, keyFunc func(T) K) []T {
	if len(slice) == 0 {
		return slice
	}

	seen := make(map[K]struct{})
	result := make([]T, 0, len(slice))

	for _, item := range slice {
		key := keyFunc(item)
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

// Union 合并所有给定切片中的唯一元素，保持顺序。
func Union[T comparable](slices ...[]T) []T {
	var result []T
	contain := map[T]struct{}{}

	for _, slice := range slices {
		for _, item := range slice {
			if _, ok := contain[item]; !ok {
				contain[item] = struct{}{}
				result = append(result, item)
			}
		}
	}

	return result
}

// UnionBy 类似Union，但接受一个谓词函数来处理每个元素。
func UnionBy[T any, V comparable](predicate func(item T) V, slices ...[]T) []T {
	result := []T{}
	contain := map[V]struct{}{}

	for _, slice := range slices {
		for _, item := range slice {
			val := predicate(item)
			if _, ok := contain[val]; !ok {
				contain[val] = struct{}{}
				result = append(result, item)
			}
		}
	}

	return result
}

// Merge 将所有给定的切片合并成一个切片。
func Merge[T any](slices ...[]T) []T {
	totalLen := 0
	for _, v := range slices {
		totalLen += len(v)
	}
	result := make([]T, 0, totalLen)

	for _, v := range slices {
		result = append(result, v...)
	}

	return result
}

// Intersection 创建一个切片，包含所有给定切片中的唯一公共元素。交集
func Intersection[T comparable](slices ...[]T) []T {
	if len(slices) == 0 {
		return []T{}
	}
	if len(slices) == 1 {
		return Unique(slices[0])
	}

	reducer := func(sliceA, sliceB []T) []T {
		hashMap := make(map[T]int)
		for _, v := range sliceA {
			hashMap[v] = 1
		}

		out := make([]T, 0)
		for _, val := range sliceB {
			if v, ok := hashMap[val]; v == 1 && ok {
				out = append(out, val)
				hashMap[val]++
			}
		}
		return out
	}

	result := reducer(slices[0], slices[1])

	reduceSlice := make([][]T, 2)
	for i := 2; i < len(slices); i++ {
		reduceSlice[0] = result
		reduceSlice[1] = slices[i]
		result = reducer(reduceSlice[0], reduceSlice[1])
	}

	return result
}

// SymmetricDifference 对切片执行对称差集操作。
func SymmetricDifference[T comparable](slices ...[]T) []T {
	if len(slices) == 0 {
		return []T{}
	}
	if len(slices) == 1 {
		return Unique(slices[0])
	}

	result := make([]T, 0)

	intersectSlice := Intersection(slices...)

	for i := 0; i < len(slices); i++ {
		slice := slices[i]
		for _, v := range slice {
			if !Contain(intersectSlice, v) {
				result = append(result, v)
			}
		}

	}

	return Unique(result)
}

// Reverse 反转切片中元素的顺序。
func Reverse[T any](slice []T) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// Shuffle 随机打乱切片中的元素顺序。
func Shuffle[T any](slice []T) []T {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})

	return slice
}

// Sort 对任何有序类型（数字或字符串）的切片进行排序，默认升序，可通过参数指定降序。
func Sort[T constraints.Ordered](slice []T, desc bool) {
	if desc {
		quickSort(slice, 0, len(slice)-1, "desc")
	} else {
		quickSort(slice, 0, len(slice)-1, "asc")
	}
}

// SortBy 根据提供的比较函数对切片进行排序。
func SortBy[T any](slice []T, less func(a, b T) bool) {
	quickSortBy(slice, 0, len(slice)-1, less)
}

// Without 创建一个新切片，排除给定的所有项。
func Without[T comparable](slice []T, items ...T) []T {
	if len(items) == 0 || len(slice) == 0 {
		return slice
	}

	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if !Contain(items, v) {
			result = append(result, v)
		}
	}

	return result
}

// IndexOf 查找元素在切片中的索引，未找到返回-1。
func IndexOf[T comparable](arr []T, val T) int {
	limit := 10
	// gets the hash value of the array as the key of the hash table.
	key := fmt.Sprintf("%p", arr)
	// determines whether the hash table is empty. If so, the hash table is created.
	if memoryHashMap[key] == nil {
		memoryHashMap[key] = make(map[any]int)
		// iterate through the array, adding the value and index of each element to the hash table.
		for i := len(arr) - 1; i >= 0; i-- {
			memoryHashMap[key][arr[i]] = i
		}
	}
	// update the hash table counter.
	memoryHashCounter[key]++

	// use the hash table to find the specified value. If found, the index is returned.
	if index, ok := memoryHashMap[key][val]; ok {
		// calculate the memory usage of the hash table.
		size := len(memoryHashMap)
		// If the memory usage of the hash table exceeds the memory limit, the hash table with the lowest counter is cleared.
		if size > limit {
			var minKey string
			var minVal int
			for k, v := range memoryHashCounter {
				if k == key {
					continue
				}
				if minVal == 0 || v < minVal {
					minKey = k
					minVal = v
				}
			}
			delete(memoryHashMap, minKey)
			delete(memoryHashCounter, minKey)
		}
		return index
	}
	return -1
}

// LastIndexOf 查找元素在切片中最后一次出现的索引，未找到返回-1。

func LastIndexOf[T comparable](slice []T, item T) int {
	for i := len(slice) - 1; i >= 0; i-- {
		if item == slice[i] {
			return i
		}
	}

	return -1
}

// ToSlicePointer 将变长参数转换为指针切片。
func ToSlicePointer[T any](items ...T) []*T {
	result := make([]*T, len(items))
	for i := range items {
		result[i] = &items[i]
	}

	return result
}

// ToSlice 将变长参数转换为切片。
func ToSlice[T any](items ...T) []T {
	result := make([]T, len(items))
	copy(result, items)

	return result
}

// AppendIfAbsent 如果元素不存在，则追加到切片。
func AppendIfAbsent[T comparable](slice []T, item T) []T {
	if !Contain(slice, item) {
		slice = append(slice, item)
	}
	return slice
}

// AppendFirst 在切片的开头追加元素。
func AppendFirst[T any](slice []T, item T) []T {
	return append([]T{item}, slice...)
}

// SetToDefaultIf 根据谓词函数将切片元素设置为其默认值。
func SetToDefaultIf[T any](slice []T, predicate func(T) bool) ([]T, int) {
	var count int
	for i := 0; i < len(slice); i++ {
		if predicate(slice[i]) {
			var zeroValue T
			slice[i] = zeroValue
			count++
		}
	}
	return slice, count
}

// KeyBy 根据回调函数将切片转换为映射。
func KeyBy[T any, U comparable](slice []T, iteratee func(item T) U) map[U]T {
	result := make(map[U]T, len(slice))

	for _, v := range slice {
		k := iteratee(v)
		result[k] = v
	}

	return result
}

// Join 使用指定分隔符连接切片元素为字符串。
func Join[T any](slice []T, separator string) string {
	str := Map(slice, func(_ int, item T) string {
		return fmt.Sprint(item)
	})

	return strings.Join(str, separator)
}

// Partition 根据多个谓词函数将切片元素划分到不同子切片中。
func Partition[T any](slice []T, predicates ...func(item T) bool) [][]T {
	l := len(predicates)

	result := make([][]T, l+1)

	for _, item := range slice {
		processed := false

		for i, f := range predicates {
			if f == nil {
				panic("predicate function must not be nill")
			}

			if f(item) {
				result[i] = append(result[i], item)
				processed = true
				break
			}
		}

		if !processed {
			result[l] = append(result[l], item)
		}
	}

	return result
}

// Break 根据谓词函数将切片分为两部分。
func Break[T any](values []T, predicate func(T) bool) ([]T, []T) {
	a := make([]T, 0)
	b := make([]T, 0)
	if len(values) == 0 {
		return a, b
	}
	matched := false
	for _, value := range values {

		if !matched && predicate(value) {
			matched = true
		}

		if matched {
			b = append(b, value)
		} else {
			a = append(a, value)
		}
	}
	return a, b
}

// Random 获取切片中的随机元素及其索引，切片为空时返回-1。
func Random[T any](slice []T) (val T, idx int) {
	if len(slice) == 0 {
		return val, -1
	}

	idx = rand.Intn(len(slice))
	return slice[idx], idx
}

// RightPadding 在切片右侧添加指定数量的填充元素。
func RightPadding[T any](slice []T, paddingValue T, paddingLength int) []T {
	if paddingLength == 0 {
		return slice
	}
	for i := 0; i < paddingLength; i++ {
		slice = append(slice, paddingValue)
	}
	return slice
}

// LeftPadding 在切片左侧添加指定数量的填充元素。
func LeftPadding[T any](slice []T, paddingValue T, paddingLength int) []T {
	if paddingLength == 0 {
		return slice
	}

	paddedSlice := make([]T, len(slice)+paddingLength)
	i := 0
	for ; i < paddingLength; i++ {
		paddedSlice[i] = paddingValue
	}
	for j := 0; j < len(slice); j++ {
		paddedSlice[i] = slice[j]
		i++
	}

	return paddedSlice
}

func MergeArray[a any, b comparable](list []a, cover func(a) []b, duplicate bool) []b {
	tagIds := make([]b, 0)
	for _, user := range list {
		bs := cover(user)
		tagIds = append(tagIds, bs...)
	}
	if duplicate {
		return Unique(tagIds)
	}
	return tagIds
}

func SplitArray[T any](arr []T, size int) [][]T {
	var result [][]T
	for i := 0; i < len(arr); i += size {
		end := i + size
		if end > len(arr) {
			end = len(arr)
		}
		result = append(result, arr[i:end])
	}
	return result
}

// SortSpecify 根据指定顺序排序
// 1. 输入参数：
//   - list: 待排序的数组
//   - getSort: 获取集合中的自定字段
//   - sorts: 指定顺序的数组
func SortSpecify[T any, b comparable](list []T, getSort func(T) b, sorts []b) []T {
	ts := make([]T, 0)
	for _, c := range sorts {
		for _, l := range list {
			sort := getSort(l)
			if sort == c {
				ts = append(ts, l)
				break
			}
		}
	}
	return ts
}

func IsEmpty[T any](array []T) bool {
	if array == nil {
		return true
	}
	return len(array) == 0
}

func FromSlice[T any](t ...T) []T {
	return t
}

func GetIndexVal[T any](array []T, index int, nilTo T) (T, bool) {
	if array == nil {
		return nilTo, false
	}
	if index < 0 || index >= len(array) {
		return nilTo, false
	}
	return array[index], true
}

// FirstOrDefault 返回数组的第一个元素，如果数组为空则返回默认值 zeroToDefault 空值是否也返回默认值 默认true
func FirstOrDefault[T any](array []T, defaultValue T, zeroToDefault ...bool) T {
	if len(array) > 0 {
		t := array[0]
		if len(zeroToDefault) == 0 || (len(zeroToDefault) > 0 && zeroToDefault[0]) {
			if xcore.IsZero(t) {
				return defaultValue
			}
		}
		return t
	}
	return defaultValue
}
