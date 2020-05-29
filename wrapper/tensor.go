package wrapper

// #include <stdlib.h>
import "C"

import (
	// "fmt"
	"reflect"

	gotch "github.com/sugarme/gotch"
	lib "github.com/sugarme/gotch/libtch"
)

type Tensor struct {
	ctensor *lib.C_tensor
}

// NewTensor creates a new tensor
func NewTensor() Tensor {
	ctensor := lib.AtNewTensor()
	return Tensor{ctensor}
}

// FOfSlice creates tensor from a slice data
func (ts Tensor) FOfSlice(data interface{}, dtype gotch.DType) (retVal *Tensor, err error) {

	dataLen := reflect.ValueOf(data).Len()
	shape := []int64{int64(dataLen)}
	elementNum := ElementCount(shape)
	// eltSizeInBytes := dtype.EltSizeInBytes() // Element Size in Byte for Int dtype
	eltSizeInBytes := gotch.DTypeSize(dtype)

	nbytes := int(eltSizeInBytes) * int(elementNum)

	dataPtr, buff := CMalloc(nbytes)

	if err = EncodeTensor(buff, reflect.ValueOf(data), shape); err != nil {
		return nil, err
	}

	ctensor := lib.AtTensorOfData(dataPtr, shape, uint(len(shape)), uint(eltSizeInBytes), int(gotch.DType2CInt(dtype)))

	retVal = &Tensor{ctensor}

	// Read back created tensor values by C libtorch
	// readDataPtr := lib.AtDataPtr(retVal.ctensor)
	// readDataSlice := (*[1 << 30]byte)(readDataPtr)[:nbytes:nbytes]
	// // typ := typeOf(dtype, shape)
	// typ := reflect.TypeOf(int32(0)) // C. type `int` ~ Go type `int32`
	// val := reflect.New(typ)
	// if err := DecodeTensor(bytes.NewReader(readDataSlice), shape, typ, val); err != nil {
	// panic(fmt.Sprintf("unable to decode Tensor of type %v and shape %v - %v", dtype, shape, err))
	// }
	//
	// tensorData := reflect.Indirect(val).Interface()
	//
	// fmt.Println("%v", tensorData)

	return retVal, nil
}

func (ts Tensor) Print() {
	lib.AtPrint(ts.ctensor)
}