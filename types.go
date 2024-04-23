// Copyright 2024 Arianit Uka
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package openapi

// Basic initialized types that can be easily passed to any Body() in builder.go
var (
	IntType   = int(0)
	Int8Type  = int8(0)
	Int16Type = int16(0)
	Int32Type = int32(0)
	Int64Type = int64(0)

	UintType    = uint(0)
	Uint8Type   = uint8(0)
	Uint16Type  = uint16(0)
	Uint32Type  = uint32(0)
	Uint64Type  = uint64(0)
	UintPtrType = uintptr(0)

	Float32Type = float32(0)
	Float64Type = float64(0)

	Complex64Type  = complex64(0)
	Complex128Type = complex128(0)

	StringType = ""
	BoolType   = false
	ByteType   = byte(0)
	RuneType   = rune(0)
)
