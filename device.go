package main

/*
#include "winaudio_wrapper.h"
#include "notification.h"
#include "winhelper.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type DeviceState struct {
	ID    string
	Name  string
}

type Device struct {
	device   *C.IMMDevice
	volume   *C.IAudioEndpointVolume
	callback *C.IAudioEndpointVolumeCallback

	DeviceState
}

func newDevice(immDevice *C.IMMDevice) (*Device, error) {
	device := &Device{device: immDevice}

	err := device.getID()
	if err != nil {
		device.release()
		return nil, fmt.Errorf("device id: %w", err)
	}

	err = device.getName()
	if err != nil {
		device.release()
		return nil, fmt.Errorf("device name: %w", err)
	}

	return device, nil
}

func (d *Device) getID() error {
	var id C.LPWSTR
	hr := C.IMMDevice_GetId(d.device, &id)
	if hr < 0 {
		return fmt.Errorf("device id: 0x%x", uint32(hr))
	}
	defer C.freeUTF16String(id)

	d.ID = LPWSTRToStr(id)
	return nil
}

func (d *Device) getName() error {
	var store *C.IPropertyStore
	if hr := C.IMMDevice_OpenPropertyStore(d.device, &store); hr < 0 {
		return fmt.Errorf("device %s property store: 0x%x", d.ID, uint32(hr))
	}
	defer C.IPropertyStore_Release(store)

	var prop C.PROPVARIANT
	if hr := C.IPropertyStore_GetValue(store, &C.PKEY_Device_FriendlyName, &prop); hr < 0 {
		return fmt.Errorf("device %s name value: 0x%x", d.ID, uint32(hr))
	}
	defer C.PropVariantClear(&prop) // Will handle also the name free

	name := C.PROPVARIANT_GetStringValue(&prop)
	d.Name = LPWSTRToStr(name)
	return nil
}

func (d *Device) initVolume() error {
	if d.volume != nil {
		return nil
	}

	hr := C.IMMDevice_Activate(
		d.device,
		&C.IID_IAudioEndpointVolume,
		C.CLSCTX_ALL,
		nil,
		(*unsafe.Pointer)(unsafe.Pointer(&d.volume)),
	)
	if hr < 0 {
		return fmt.Errorf("device %s endpoint volume: 0x%x", d.ID, uint32(hr))
	}

	return nil
}

func (d *Device) releaseVolume() error {
	if d.volume == nil {
		return nil
	}
	defer func() { d.volume = nil }()

	// This must be done to avoid the release function
	// to block the thread, don't ask me why
	volume := d.volume
	go C.IAudioEndpointVolume_Release(volume)

	return nil
}

func (d *Device) getMuted() (bool, error) {
	if d.volume == nil {
		return false, nil
	}

	var mutedInt C.int
	hr := C.IAudioEndpointVolume_GetMute(d.volume, &mutedInt)
	if hr < 0 {
		return false, fmt.Errorf("device %s get mute: 0x%x", d.ID, uint32(hr))
	}
	return mutedInt != 0, nil
}

func (d *Device) setMuted(muted bool) error {
	if d.volume == nil {
		return nil
	}

	var cMuted C.BOOL
	if muted {
		cMuted = 1
	}

	hr := C.IAudioEndpointVolume_SetMute(d.volume, cMuted, nil)
	if hr < 0 {
		return fmt.Errorf("device %s get mute: 0x%x", d.ID, uint32(hr))
	}

	return nil
}

func (d *Device) registerControlChangeNotify() error {
	if d.callback != nil {
		return nil
	}

	if d.volume == nil {
		return nil
	}

	hr := C.RegisterControlChangeNotify(d.volume, &d.callback)
	if hr < 0 {
		return fmt.Errorf("device %s register notify: 0x%x", d.ID, uint32(hr))
	}

	return nil
}

func (d *Device) unregisterControlChangeNotify() error {
	if d.callback == nil {
		return nil
	}

	if d.volume == nil {
		panic("unexpected device state: callback present with volume released")
	}

	defer func() { d.callback = nil }()

	hr := C.UnregisterControlChangeNotify(d.volume, d.callback)
	if hr < 0 {
		return fmt.Errorf("device %s unregister notify: 0x%x", d.ID, uint32(hr))
	}

	return nil
}

func (d *Device) activate() error {
	err := d.initVolume()
	if err != nil {
		return err
	}

	audioService.Muted, err = d.getMuted()
	if err != nil {
		return err
	}

	err = d.registerControlChangeNotify()
	if err != nil {
		return err
	}
	
	return nil
}

func (d *Device) deactivate() error {
	err := d.unregisterControlChangeNotify()
	if err != nil {
		return err
	}

	err = d.releaseVolume()
	if err != nil {
		return err
	}

	return nil
}

func (d *Device) copyStateFrom(other *Device) {
	d.DeviceState = other.DeviceState
}

func (d *Device) release() {
	d.deactivate()
	C.IMMDevice_Release(d.device)
}
