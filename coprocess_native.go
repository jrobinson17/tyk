// +build coprocess
// +build !grpc

package main

/*
#cgo python CFLAGS: -DENABLE_PYTHON
#include <stdio.h>
#include <stdlib.h>

#include "coprocess/sds/sds.h"

#include "coprocess/api.h"

#ifdef ENABLE_PYTHON
#include "coprocess/python/dispatcher.h"
#include "coprocess/python/binding.h"
#endif

*/
import "C"

import (
	"github.com/golang/protobuf/proto"

	"github.com/TykTechnologies/tyk/coprocess"

	"encoding/json"
	"unsafe"
)

// Dispatch prepares a CoProcessMessage, sends it to the GlobalDispatcher and gets a reply.
func (c *CoProcessor) Dispatch(object *coprocess.Object) (*coprocess.Object, error) {

	var objectMsg []byte
	if MessageType == coprocess.ProtobufMessage {
		objectMsg, _ = proto.Marshal(object)
	} else if MessageType == coprocess.JsonMessage {
		objectMsg, _ = json.Marshal(object)
	}

	objectMsgStr := string(objectMsg)

	CObjectStr := C.CString(objectMsgStr)

	objectPtr := (*C.struct_CoProcessMessage)(C.malloc(C.size_t(unsafe.Sizeof(C.struct_CoProcessMessage{}))))
	objectPtr.p_data = unsafe.Pointer(CObjectStr)
	objectPtr.length = C.int(len(objectMsg))

	newObjectPtr := (*C.struct_CoProcessMessage)(GlobalDispatcher.Dispatch(unsafe.Pointer(objectPtr)))

	newObjectBytes := C.GoBytes(newObjectPtr.p_data, newObjectPtr.length)

	newObject := &coprocess.Object{}

	if MessageType == coprocess.ProtobufMessage {
		proto.Unmarshal(newObjectBytes, newObject)
	} else if MessageType == coprocess.JsonMessage {
		json.Unmarshal(newObjectBytes, newObject)
	}

	C.free(unsafe.Pointer(CObjectStr))
	C.free(unsafe.Pointer(objectPtr))
	C.free(unsafe.Pointer(newObjectPtr))

	return newObject, nil
}
