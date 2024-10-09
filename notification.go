//go:generate go tool cgo -exportheader notification_export.h .\notification.go
//go:generate pwsh -nop -c rm -r .\_obj
package main

/*
#include "notification.h"
*/
import "C"
import (
	"log"
)

//export OnDeviceStateChangedCallback
func OnDeviceStateChangedCallback(pwstrDeviceId C.LPCWSTR, dwNewState C.DWORD) C.HRESULT {
	err := audioService.updateFrontend(true)
	if err != nil {
		log.Printf("OnDeviceStateChangedCallback error: %v\n", err)
		return C.E_FAIL
	}
	return C.S_OK
}

//export OnDeviceAddedCallback
func OnDeviceAddedCallback(pwstrDeviceId C.LPCWSTR) C.HRESULT {
	err := audioService.updateFrontend(true)
	if err != nil {
		log.Printf("OnDeviceAddedCallback error: %v\n", err)
		return C.E_FAIL
	}
	return C.S_OK
}

//export OnDeviceRemovedCallback
func OnDeviceRemovedCallback(pwstrDeviceId C.LPCWSTR) C.HRESULT {
	err := audioService.updateFrontend(true)
	if err != nil {
		log.Printf("OnDeviceRemovedCallback error: %v\n", err)
		return C.E_FAIL
	}
	return C.S_OK
}

//export OnDefaultDeviceChangedCallback
func OnDefaultDeviceChangedCallback(flow C.EDataFlow, role C.ERole, pwstrDefaultDeviceId C.LPCWSTR) C.HRESULT {
	return C.S_OK
}

//export OnPropertyValueChangedCallback
func OnPropertyValueChangedCallback(pwstrDeviceId C.LPCWSTR, key C.PROPERTYKEY) C.HRESULT {
	return C.S_OK
}

//export OnEndpointVolumeChangeNotify
func OnEndpointVolumeChangeNotify(pNotify C.PAUDIO_VOLUME_NOTIFICATION_DATA) C.HRESULT {
	err := audioService.setPrefMute(pNotify.bMuted != 0)
	if err != nil {
		log.Printf("OnEndpointVolumeChangeNotify error: %v\n", err)
		return C.E_FAIL
	}

	err = audioService.updateFrontend(false)
	if err != nil {
		log.Printf("OnEndpointVolumeChangeNotify error: %v\n", err)
		return C.E_FAIL
	}

	return C.S_OK
}
