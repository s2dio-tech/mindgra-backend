package common

import (
	"reflect"
	"sort"
)

type Stream struct {
	slice reflect.Value
}

func StreamOf(slice any) *Stream {
	rv := reflect.ValueOf(slice)
	return &Stream{
		slice: rv,
	}
}

/*
dst := StreamOf(hoges).
	Filter(func(hoge *Hoge) bool {
		return hoge.Num > 3
	}).Out().([]*Hoge)
*/
// 要素のフィルタリング
func (s *Stream) Filter(fn any) *Stream {
	frv := reflect.ValueOf(fn)
	srv := reflect.MakeSlice(s.slice.Type(), 0, 0)
	for i := 0; i < s.slice.Len(); i++ {
		rv := s.slice.Index(i)
		out := frv.Call([]reflect.Value{rv})
		if out[0].Bool() {
			srv = reflect.Append(srv, rv)
		}
	}
	s.slice = srv
	return s
}

/*
dst := StreamOf(hoges).
	Map(func(hoge *Hoge) string {
		return hoge.ID
	}).Out().([]string)
*/
// 要素の変換
func (s *Stream) Map(fn any) *Stream {
	frv := reflect.ValueOf(fn)
	srt := reflect.SliceOf(frv.Type().Out(0))
	srv := reflect.MakeSlice(srt, 0, 0)
	for i := 0; i < s.slice.Len(); i++ {
		rv := s.slice.Index(i)
		out := frv.Call([]reflect.Value{rv})
		srv = reflect.Append(srv, out[0])
	}
	s.slice = srv
	return s
}

/*
dst := StreamOf(hoges).
	Reduce(func(dst int, num int) int {
		return dst + num
	}).(int)
*/
// 要素の集計
func (s *Stream) Reduce(fn any) any {
	frv := reflect.ValueOf(fn)
	rt := frv.Type().Out(0)
	dst := reflect.New(rt).Elem()
	for i := 0; i < s.slice.Len(); i++ {
		rv := s.slice.Index(i)
		out := frv.Call([]reflect.Value{dst, rv})
		dst = out[0]
	}
	return dst.Interface()
}

/*
dst := StreamOf(hoges).
	Sort(func(prev, next *Hoge) bool {
		return prev.SortNum < next.SortNum
	}).Out().([]*Hoge)
*/
// ソート
func (s *Stream) Sort(fn any) *Stream {
	frv := reflect.ValueOf(fn)
	slice := s.slice.Interface()
	sort.SliceStable(slice, func(i, j int) bool {
		out := frv.Call([]reflect.Value{s.slice.Index(i), s.slice.Index(j)})
		return out[0].Bool()
	})
	s.slice = reflect.ValueOf(slice)
	return s
}

/*
dst := StreamOf(hoges).
    Contains(func(hoge *Hoge) bool {
		return hoge.ID == "abc"
	})
*/
// 要素の存在確認
func (s *Stream) Contains(fn any) bool {
	frv := reflect.ValueOf(fn)
	for i := 0; i < s.slice.Len(); i++ {
		rv := s.slice.Index(i)
		out := frv.Call([]reflect.Value{rv})
		if out[0].Bool() {
			return true
		}
	}
	return false
}

/*
dst := StreamOf(hoges).
    ForEach(func(hoge *Hoge, i int) {
		hoge.ID = "abc"
	})
*/
// 要素のループ
func (s *Stream) ForEach(fn any) *Stream {
	frv := reflect.ValueOf(fn)
	if frv.Type().NumIn() == 1 {
		for i := 0; i < s.slice.Len(); i++ {
			rv := s.slice.Index(i)
			_ = frv.Call([]reflect.Value{rv})
		}
	} else {
		for i := 0; i < s.slice.Len(); i++ {
			rv := s.slice.Index(i)
			_ = frv.Call([]reflect.Value{rv, reflect.ValueOf(i)})
		}
	}
	return s
}

/*
dst := StreamOf(hoges).Count()
*/
// 要素数を取得
func (s *Stream) Count() int {
	return s.slice.Len()
}

// 結果を出力する
func (s *Stream) Out() any {
	return s.slice.Interface()
}
