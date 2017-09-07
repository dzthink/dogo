//栈结构封装
package dogo

import (
	"sync"
	"errors"
)

const MaxUint = ^uint(0)

type Stack struct {
	element []interface{}
	ptr uint
	mutex *sync.Mutex
}

func NewStack() *Stack {
	return &Stack{
		element : make([]interface{}, 0, 20),
		ptr : 0,
		mutex : new(sync.Mutex),
	}
}


func(s *Stack)Push(e interface{}) error {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	if s.ptr + 1 == MaxUint {
		return errors.New("stack full")
	}
	if s.ptr < uint(len(s.element)) {
		s.element[s.ptr] = e
	} else {
		s.element = append(s.element, e)
	}
	s.ptr = s.ptr + 1
	return nil
}

func(s *Stack)Pop()(interface{}, error) {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	if s.ptr == 0 {
		return nil, errors.New("Statck empty")
	}
	s.ptr = s.ptr - 1
	return s.element[s.ptr], nil
}

func(s *Stack)Peek()(interface{}) {
	if s.ptr == 0 {
		return nil
	}
	return s.element[s.ptr - 1]
}

func(s *Stack)Len() uint {
	return s.ptr
}
