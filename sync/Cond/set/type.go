package set

import "sort"

type Empty struct{}

type Sets interface {
	Insert(items ...string) Set
	Delete(items ...string) Set
	Has(item string) bool
	HasAll(items ...string) bool
	HasAny(items ...string) bool
	Difference(s2 Set) Set
	Len() int
	Intersection(s2 Set) Set
	IsSuperset(s2 Set) bool
	Equal(s2 Set) bool
	List() []string
	UnsortedList() []string
}
type Set map[string]Empty

func NewSet(items ...string) Set {
	ss := Set{}
	ss.Insert(items...)
	return ss
}

// 返回长度
func (s Set) Len() int {
	return len(s)
}

// 插入元素
func (s Set) Insert(items ...string) Set {
	for _, item := range items {
		// 这里实际上利用map实现了去重
		s[item] = Empty{}
	}
	return s
}

// 删除元素
func (s Set) Delete(items ...string) Set {
	for _, item := range items {
		delete(s, item)
	}
	return s
}

// 判断某个元素是否存在
func (s Set) Has(item string) bool {
	_, contained := s[item]
	return contained
}

// 判断所有元素是否存在
func (s Set) HasAll(items ...string) bool {
	for _, item := range items {
		if !s.Has(item) {
			return false
		}
	}
	return true
}

// 判断是否存在任意一个
func (s Set) HasAny(items ...string) bool {
	for _, item := range items {
		if s.Has(item) {
			return true
		}
	}
	return false
}

// Difference returns a set of objects that are not in s2
// For example:
// s1 = {a1, a2, a3}
// s2 = {a1, a2, a4, a5}
// s1.Difference(s2) = {a3}
// s2.Difference(s1) = {a4, a5}
func (s Set) Difference(s2 Set) Set {
	result := NewSet()
	for key := range s {
		if !s2.Has(key) {
			result.Insert(key)
		}
	}
	return result
}

// Intersection returns a new set which includes the item in BOTH s1 and s2
// For example:
// s1 = {a1, a2}
// s2 = {a2, a3}
// s1.Intersection(s2) = {a2}
func (s Set) Intersection(s2 Set) Set {
	var less, more Set
	result := NewSet()

	if s.Len() < s2.Len() {
		less = s
		more = s2
	} else {
		less = s2
		more = s
	}

	// 从小的遍历，减少次数
	for key := range less {
		// 当s2也有这个key时
		if more.Has(key) {
			result.Insert(key)
		}
	}
	return result
}

// IsSuperset returns true if and only if s1 is a superset of s2.
func (s Set) IsSuperset(s2 Set) bool {
	for item := range s2 {
		if !s.Has(item) {
			return false
		}
	}
	return true
}

// 当s1和s2长度一致，且s2中每个元素s1都有
func (s Set) Equal(s2 Set) bool {
	return len(s) == len(s2) && s.IsSuperset(s2)
}

// 返回集合所有的元素(排序好的)
func (s Set) List() []string {
	res := make(sortableSliceOfSet, 0, len(s))
	for key := range s {
		res = append(res, key)
	}
	// 排序
	sort.Sort(res)
	return []string(res)
}

// 返回集合所有的元素(未排序)
func (s Set) UnsortedList() []string {
	res := make([]string, 0, len(s))
	for key := range s {
		res = append(res, key)
	}
	return res
}

// Returns a single element from the set
func (s Set) PopAny() (string, bool) {
	for key := range s {
		s.Delete(key)
		return key, true
	}
	var zeroValue string
	return zeroValue, false
}

type sortableSliceOfSet []string

func (s sortableSliceOfSet) Len() int           { return len(s) }
func (s sortableSliceOfSet) Less(i, j int) bool { return lessString(s[i], s[j]) }
func (s sortableSliceOfSet) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func lessString(lhs, rhs string) bool {
	return lhs < rhs
}
