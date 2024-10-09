#include "winaudio_wrapper.h"

HRESULT CreateInstance(IMMDeviceEnumerator** deviceEnum) {
	return CoCreateInstance(CLSID_MMDeviceEnumerator, NULL, CLSCTX_ALL, IID_IMMDeviceEnumerator, (void**)deviceEnum);
}

//
// IMMDeviceEnumerator
//

HRESULT IMMDeviceEnumerator_EnumAudioEndpoints(IMMDeviceEnumerator* deviceEnum, IMMDeviceCollection** collection) {
	return deviceEnum->EnumAudioEndpoints(eCapture, DEVICE_STATE_ACTIVE, collection);
}

void IMMDeviceEnumerator_Release(IMMDeviceEnumerator* deviceEnum) {
	deviceEnum->Release();
}

//
// IMMDeviceCollection
//

HRESULT IMMDeviceCollection_GetCount(IMMDeviceCollection* collection, UINT* count) {
	return collection->GetCount(count);
}

HRESULT IMMDeviceCollection_Item(IMMDeviceCollection* collection, UINT index, IMMDevice** device) {
	HRESULT hr = collection->Item(index, device);
	if (SUCCEEDED(hr)) {
		return hr;
	}
		
	return hr;
}

void IMMDeviceCollection_Release(IMMDeviceCollection* collection) {
	collection->Release();
}

//
// IMMDevice
//

// After use, free id with CoTaskMemFree
HRESULT IMMDevice_GetId(IMMDevice* device, LPWSTR* id) {
	return device->GetId(id);
}

HRESULT IMMDevice_OpenPropertyStore(IMMDevice* device, IPropertyStore** store) {
	return device->OpenPropertyStore(STGM_READ, store);
}

HRESULT IMMDevice_Activate(IMMDevice* device, const IID* iid, DWORD dwClsCtx, PROPVARIANT* pActivationParams, void** ppInterface) {
	return device->Activate(*iid, dwClsCtx, pActivationParams, ppInterface);
}

void IMMDevice_Release(IMMDevice* device) {
	device->Release();
}

//
// IPropertyStore
//

HRESULT IPropertyStore_GetValue(IPropertyStore* store, const PROPERTYKEY* key, PROPVARIANT* prop) {
	PropVariantInit(prop);
	return store->GetValue(*key, prop);
}

void IPropertyStore_Release(IPropertyStore* store) {
	store->Release();
}

//
// PROPVARIANT
//

LPWSTR PROPVARIANT_GetStringValue(PROPVARIANT* prop) {
	return prop->pwszVal;
}

//
// IAudioEndpointVolume
//

HRESULT IAudioEndpointVolume_GetMute(IAudioEndpointVolume* volume, BOOL* muted) {
	return volume->GetMute(muted);
}

HRESULT IAudioEndpointVolume_SetMute(IAudioEndpointVolume* volume, BOOL muted, LPCGUID context) {
	return volume->SetMute(muted, context);
}

void IAudioEndpointVolume_Release(IAudioEndpointVolume* volume) {
	volume->Release();
}
