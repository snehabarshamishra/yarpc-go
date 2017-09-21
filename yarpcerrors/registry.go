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
	"errors"
	"fmt"
	"reflect"
)

var (
	_errorReflectType  = reflect.TypeOf((*error)(nil)).Elem()
	_reflectTypeToCode = make(map[reflect.Type]Code, 0)
	_reflectTypeToName = make(map[reflect.Type]string, 0)
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
	reflectType, err := getAndCheckErrorReflectType(emptyError)
	if err != nil {
		panic(err.Error())
	}
	_reflectTypeToCode[reflectType] = code
}

// RegisterErrorName registers a type of error to a specific name.
//
// This information will be used by transports to set a corresponding
// error name over the wire if an error of such type is returned from
// a handler. If no code is registered with RegisterErrorCode,
// CodeUnknown will be passed over the wire.
//
// An empty instance of the error should be passed as err.
// Only structs or struct pointers can be registered.
// If these conditions are not met, this will panic.
//
// This is not thread-safe an should only be called at initialization.
// Values registered later will overwrite values registered earlier.
//
// Example:
//
//   func init() {
//     yarpcerrors.RegisterErrorName(&PaymentDeniedError{}, "payment-denied")
//   }
//
//   type PaymentDeniedError struct {
//     Provider string
//   }
//
//   func (e *PaymentDeniedError) Error() string {
//     return fmt.Sprintf("payment denied for provider: %s", e.Provider)
//   }
func RegisterErrorName(emptyError interface{}, name string) {
	if code == CodeOK {
		panic("yarpcerrors: registered code cannot be CodeOK")
	}
	reflectType, err := getAndCheckErrorReflectType(emptyError)
	if err != nil {
		panic(err.Error())
	}
	_reflectTypeToName[reflectType] = name
}

// GetCodeForRegisteredError gets the Code for the given error,
// or CodeOK if there is no Code registered.
func GetCodeForRegisteredError(err interface{}) Code {
	if err == nil {
		return CodeOK
	}
	reflectType = reflect.Type(err)
	if reflectType == nil {
		return CodeOK
	}
	code, ok := _reflectTypeToCode[reflectType]
	if !ok {
		return CodeOK
	}
	return code
}

// GetNameForRegisteredError gets the name for the given error,
// or "" if there is no name registered.
func GetNameForRegisteredError(err interface{}) string {
	if err == nil {
		return ""
	}
	reflectType = reflect.Type(err)
	if reflectType == nil {
		return ""
	}
	name, ok := _reflectTypeToName[reflectType]
	if !ok {
		return ""
	}
	return name
}

func getAndCheckErrorReflectType(err interface{}) (reflect.Type, error) {
	if err == nil {
		return nil, errors.New("yarpcerrors: given error is nil")
	}
	reflectType := reflect.TypeOf(err)
	if reflectType == nil {
		return nil, errors.New("yarpcerrors: reflect type for given error is nil")
	}
	if !reflectType.AssignableTo(_errorReflectType) {
		return nil, fmt.Errorf("yarpcerrors: %T is not assignable to error", err)
	}
	if !((reflectType.Kind() == reflect.Ptr && reflectType.Elem().Kind() == reflect.Struct) || (reflectType.Kind() == reflect.Ptr)) {
		return nil, fmt.Errorf("yarpcerrors: %T is not a struct or struct pointer")
	}
	return reflectType, nil
}
