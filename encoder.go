/*
 *
 *  * Copyright 2012-2016 Viant.
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License"); you may not
 *  * use this file except in compliance with the License. You may obtain a copy of
 *  * the License at
 *  *
 *  * http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  * Unless required by applicable law or agreed to in writing, software
 *  * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 *  * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 *  * License for the specific language governing permissions and limitations under
 *  * the License.
 *
 */

package hessian

import (
	"bytes"
	"io"
	"reflect"
	"unsafe"
)

//Encoder type
type Encoder struct {
	writer     io.Writer
	clsDefList []ClassDef
	nameMap    map[string]string
	refMap     map[unsafe.Pointer]_refElem
}

//NewEncoder new
func NewEncoder(w io.Writer, np map[string]string) *Encoder {
	if np == nil {
		np = make(map[string]string, 11)
	}
	encoder := &Encoder{
		nameMap: np,
	}
	if w != nil {
		encoder.Reset(w)
	}
	return encoder
}

//Reset reset
func (e *Encoder) Reset(w io.Writer) {
	e.writer = w
	e.clsDefList = make([]ClassDef, 0, 11)
	e.refMap = make(map[unsafe.Pointer]_refElem, 11)
}

//RegisterNameType register name type
func (e *Encoder) RegisterNameType(key string, objectName string) {
	e.nameMap[key] = objectName
}

//RegisterNameMap register name map
func (e *Encoder) RegisterNameMap(mp map[string]string) {
	e.nameMap = mp
}

//WriteObject write object
func (e *Encoder) WriteObject(data interface{}) error {
	_, err := e.WriteData(data)
	return err
}

//WriteTo write object to target writer
func (e *Encoder) WriteTo(w io.Writer, data interface{}) error {
	e.Reset(w)
	return e.WriteObject(data)
}

// Encode encode object to bytes
func (e *Encoder) Encode(object interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	err := e.WriteTo(buffer, object)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

//WriteData write object
func (e *Encoder) WriteData(data interface{}) (int, error) {
	if data == nil {
		e.writeBT(_nilTag)
		return 1, nil
	}
	source := data
	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		v = UnpackPtr(v)

		if !v.IsValid() {
			e.writeBT(_nilTag)
			return 1, nil
		}

		data = v.Interface()
	}

	switch v.Kind() {
	case reflect.Bool:
		value := data.(bool)
		return e.writeBoolean(value)
	case reflect.String:
		value := data.(string)
		return e.writeString(value)
	case reflect.Int8: // as int
		value := int32(data.(int8))
		return e.writeInt(value)
	case reflect.Int16: // as int
		value := int32(data.(int16))
		return e.writeInt(value)
	case reflect.Int32: // as int
		value := data.(int32)
		return e.writeInt(value)
	case reflect.Int: // as int
		value := int32(data.(int))
		return e.writeInt(value)
	case reflect.Uint8: // as int
		value := int32(data.(uint8))
		return e.writeInt(value)
	case reflect.Uint16: // as int
		value := int32(data.(uint16))
		return e.writeInt(value)
	case reflect.Int64: // as long
		value := data.(int64)
		return e.writeLong(value)
	case reflect.Uint: // as long
		value := int64(data.(uint))
		return e.writeLong(value)
	case reflect.Uint32: // as long
		value := int64(data.(uint32))
		return e.writeLong(value)
	case reflect.Uint64: // as long
		value := int64(data.(uint64))
		return e.writeLong(value)
	case reflect.Float32:
		value := data.(float32)
		return e.writeDouble(float64(value))
	case reflect.Float64:
		value := data.(float64)
		return e.writeDouble(value)
	case reflect.Slice, reflect.Array:
		return e.writeList(source)
	case reflect.Map:
		return e.writeMap(source)
	case reflect.Struct:
		return e.writeObject(source)
	}
	return 0, newCodecError("WriteData", "unsupported object:%v, kind:%v, type:%v", data, v.Kind(), v.Kind())
}

func (e *Encoder) writeString(value string) (int, error) {
	return e.writer.Write(encodeString(value))
}

func (e *Encoder) writeInt(value int32) (int, error) {
	return e.writer.Write(encodeInt(value))
}

func (e *Encoder) writeLong(value int64) (int, error) {
	return e.writer.Write(encodeLong(value))
}

func (e *Encoder) writeDouble(value float64) (int, error) {
	bytes, err := encodeDouble(value)
	if err != nil {
		return 0, err
	}
	return e.writer.Write(bytes)
}

func (e *Encoder) writeBoolean(value bool) (int, error) {
	return e.writer.Write(encodeBoolean(value))
}

func (e *Encoder) writeBinary(value []byte) (int, error) {
	return e.writer.Write(encodeBinary(value))
}

func (e *Encoder) writeBT(bs ...byte) (int, error) {
	return e.writer.Write(bs)
}

func (e *Encoder) writeBytes(bytes []byte) (int, error) {
	return e.writer.Write(bytes)
}
