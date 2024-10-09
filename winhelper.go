package main

/*
#include "winhelper.h"
*/
import "C"
import (
	"syscall"
	"unsafe"
)

func winUTF16ToStr(str *uint16) string {
	var res []uint16
	
	for charPtr := unsafe.Pointer(str); *(*uint16)(charPtr) != 0; charPtr = unsafe.Add(charPtr, 2) {
		res = append(res, *(*uint16)(charPtr))
	}

	return syscall.UTF16ToString(res)
}

func LPWSTRToStr(str C.LPWSTR) string {
	return winUTF16ToStr((*uint16)(unsafe.Pointer(str)))
}

func LPCWSTRToStr(str C.LPCWSTR) string {
	return winUTF16ToStr((*uint16)(unsafe.Pointer(str)))
}
