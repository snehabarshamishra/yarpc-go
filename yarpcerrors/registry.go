// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package yarpcerrors

import (
	"fmt"
	"reflect"
)

var (
	_errorReflectType  = reflect.TypeOf((*error)(nil)).Elem()
	_reflectTypeToCode = make(map[reflect.Type]Code, 0)
)

// RegisterErrorCode registers a type of error to a specific Code.
//
// This information will be used by transports to set a corresponding
// error code over the wire if an error of such type is returned from
// a handler.
//
// An empty instance of the error should be passed as err.
// Only structs or struct pointers can be registered.
// The Code cannot be CodeOK.
// If these conditions are not met, this will panic.
//
// This is not thread-safe an should only be called at initialization.
// Values registered later will overwrite values registered earlier.
//
// Example:
//
//   func init() {
//     yarpcerrors.RegisterErrorCode(&NotFoundError{}, yarpcerrors.CodeNotFound)
//   }
//
//   type NotFoundError struct {
//     Key string
//   }
//
//   func (e *NotFoundError) Error() string {
//     return fmt.Sprintf("key not found: %s", e.Key)
//   }
func RegisterErrorCode(emptyError interface{}, code Code) {
	if code == CodeOK {
		panic("yarpcerrors: registered code cannot be CodeOK")
	}
	if emptyError == nil {
		panic("yarpcerrors: given error is nil")
	}
	reflectType := reflect.TypeOf(emptyError)
	if reflectType == nil {
		panic("yarpcerrors: reflect type for given error is nil")
	}
	if !reflectType.AssignableTo(_errorReflectType) {
		panic(fmt.Sprintf("yarpcerrors: %T is not assignable to error", emptyError))
	}
	if !((reflectType.Kind() == reflect.Ptr && reflectType.Elem().Kind() == reflect.Struct) || (reflectType.Kind() == reflect.Ptr)) {
		panic(fmt.Sprintf("yarpcerrors: %T is not a struct or struct pointer", emptyError))
	}
	_reflectTypeToCode[reflectType] = code
}
