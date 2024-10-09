#ifndef WINAUDIO_WRAPPER_H
#define WINAUDIO_WRAPPER_H

#define INITGUID

#include <windows.h>
#include <mmdeviceapi.h>
#include <functiondiscoverykeys_devpkey.h>
#include <endpointvolume.h>
#include <combaseapi.h>
#include <stdio.h>

#ifdef __cplusplus
extern "C" {
#endif // __cplusplus

	HRESULT CreateInstance(IMMDeviceEnumerator** deviceEnum);

	HRESULT IMMDeviceEnumerator_EnumAudioEndpoints(IMMDeviceEnumerator* deviceEnum, IMMDeviceCollection** collection);
	void IMMDeviceEnumerator_Release(IMMDeviceEnumerator* deviceEnum);

	HRESULT IMMDeviceCollection_GetCount(IMMDeviceCollection* collection, UINT* count);
	HRESULT IMMDeviceCollection_Item(IMMDeviceCollection* collection, UINT index, IMMDevice** device);
	void IMMDeviceCollection_Release(IMMDeviceCollection* collection);

	HRESULT IMMDevice_GetId(IMMDevice* device, LPWSTR* id);
	HRESULT IMMDevice_OpenPropertyStore(IMMDevice* device, IPropertyStore** store);
	HRESULT IMMDevice_Activate(IMMDevice* device, const IID* iid, DWORD dwClsCtx, PROPVARIANT* pActivationParams, void** ppInterface);
	void IMMDevice_Release(IMMDevice* device);

	HRESULT IPropertyStore_GetValue(IPropertyStore* store, const PROPERTYKEY* key, PROPVARIANT* prop);
	void IPropertyStore_Release(IPropertyStore* store);

	LPWSTR PROPVARIANT_GetStringValue(PROPVARIANT* prop);

	HRESULT IAudioEndpointVolume_GetMute(IAudioEndpointVolume* volume, BOOL* muted);
	HRESULT IAudioEndpointVolume_SetMute(IAudioEndpointVolume* volume, BOOL muted, LPCGUID context);
	void IAudioEndpointVolume_Release(IAudioEndpointVolume* volume);

#ifdef __cplusplus
}
#endif // __cplusplus

#endif // WINAUDIO_WRAPPER_H