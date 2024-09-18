package ranger

import (
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

type Ranger[T int | string] struct {
	Map    *sync.Map
	length int32
}

func New[T int | string](vs ...T) *Ranger[T] {
	r := &Ranger[T]{
		Map:    &sync.Map{},
		length: 0,
	}
	r.Push(vs...)
	return r
}

func ParseNumber(s string) *Ranger[int] {
	numberRanger := New[int]()
	colon := strings.Split(s, ",")
	for _, express := range colon {
		numberRanger.Push(generateNumberRange(express)...)
	}
	return numberRanger
}

func generateNumberRange(s string) (ports []int) {
	dash := strings.Split(s, "-")
	if len(dash) > 2 {
		return ports
	}

	port1, err := strconv.Atoi(dash[0])
	if err != nil {
		return ports
	}

	ports = append(ports, port1)
	if len(dash) == 1 {
		return ports
	}

	port2, err := strconv.Atoi(dash[1])
	if err != nil {
		return ports
	}

	for port := port1 + 1; port <= port2; port++ {
		ports = append(ports, port)
	}
	return ports
}

func (r *Ranger[T]) Length() int32 {
	return r.length
}

func (r *Ranger[T]) Clear() {
	r.length = 0
	r.Map = &sync.Map{}
	return
}

func (r *Ranger[T]) Push(vs ...T) {
	for _, v := range vs {
		if _, ok := r.Map.Load(v); ok == true {
			continue
		}
		r.Map.Store(v, r.length)
		atomic.AddInt32(&r.length, 1)
	}
}

// todo func (r *Ranger[T]) Remove(vs ...T) ，需删除元素，且不影响Value()结果输出，暂无高质量方案

func (r *Ranger[T]) Contains(v any) bool {
	if r == nil {
		return false
	}
	if _, ok := r.Map.Load(v); ok {
		return true
	}
	return false
}

func (r *Ranger[T]) ContainsAny(vs ...any) bool {
	if r == nil {
		return false
	}
	for _, v := range vs {
		if _, ok := r.Map.Load(v); !ok {
			return false
		}
	}
	return true
}

func (r *Ranger[T]) ContainsAll(vs ...any) bool {
	if r == nil {
		return false
	}
	for _, v := range vs {
		if _, ok := r.Map.Load(v); !ok {
			return false
		}
	}
	return true
}

// Value 取值是构造新的切片，会遍历一次全Map
func (r *Ranger[T]) Value() []T {
	if r == nil {
		return []T{}
	}

	var list = make([]T, r.length)

	r.Map.Range(func(key, value any) bool {
		list[int(value.(int32))] = key.(T)
		return true
	})
	return list
}

// Range 过程是无序的
func (r *Ranger[T]) Range(fn func(key T)) {
	r.Map.Range(func(key, value any) bool {
		fn(key.(T))
		return true
	})
}

func (r *Ranger[T]) Sort(getNumber func(key T) int) {
	var list = r.Value()
	sort.Slice(list, func(i, j int) bool {
		return getNumber(list[i]) < getNumber(list[j])
	})

	for i, key := range list {
		r.Map.Store(key, int32(i))
	}
}

func (r *Ranger[T]) Intersection(r1 []T) *Ranger[T] {
	var r2 = New[T]()
	for _, cell := range r1 {
		if r.Contains(cell) {
			r2.Push(cell)
		}
	}
	return r2
}

func (r *Ranger[T]) Subtract(s []T) *Ranger[T] {
	var r1 = New[T](s...)
	var r2 = New[T]()

	r.Map.Range(func(key, value any) bool {
		if r1.ContainsAll(key.(T)) {
			return true
		}
		r2.Push(key.(T))
		return true
	})
	return r2
}
