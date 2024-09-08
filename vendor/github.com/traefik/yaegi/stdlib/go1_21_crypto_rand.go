// Code generated by 'yaegi extract crypto/rand'. DO NOT EDIT.

//go:build go1.21 && !go1.22
// +build go1.21,!go1.22

package stdlib

import (
	"crypto/rand"
	"reflect"
)

func init() {
	Symbols["crypto/rand/rand"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"Int":    reflect.ValueOf(rand.Int),
		"Prime":  reflect.ValueOf(rand.Prime),
		"Read":   reflect.ValueOf(rand.Read),
		"Reader": reflect.ValueOf(&rand.Reader).Elem(),
	}
}
