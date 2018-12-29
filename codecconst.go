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

import "reflect"

const (
	mask = byte(127)
	flag = byte(128)
)

const (
	TagRead = int32(-1)

	AsciiGap  = 32
	ChunkSize = 4096

	BcDate       = byte(0x4a) // 64-bit millisecond UTC date
	BcDateMinute = byte(0x4b) // 32-bit minute UTC date

	BcEnd = byte('Z')

	BcListVariable        = byte(0x55)
	BcListFixed           = byte('V')
	BcListVariableUntyped = byte(0x57)
	BcListFixedUntyped    = byte(0x58)

	BcListDirect        = byte(0x70)
	BcListDirectUntyped = byte(0x78)
	ListDirectMax       = byte(0x7)

	BcMap        = byte('M')
	BcMapUntyped = byte('H')

	BcNull = byte('N')

	BcObject    = byte('O')
	BcObjectDef = byte('C')

	BcObjectDirect  = byte(0x60)
	ObjectDirectMax = byte(0x0f)

	BcRef = byte(0x51)

	PPacketChunk = byte(0x4f)
	PPacket      = byte('P')

	PPacketDirect   = byte(0x80)
	PacketDirectMax = byte(0x7f)

	PPacketShort   = byte(0x70)
	PacketShortMax = 0xfff
	ArrayString    = "[string"
	ArrayInt       = "[int"
	ArrayDouble    = "[double"
	ArrayFloat     = "[float"
	ArrayBool      = "[boolean"
	ArrayLong      = "[long"
)

var (
	_interfaceInstance interface{} = 1
	interfaceType                  = reflect.TypeOf(_interfaceInstance)
)
