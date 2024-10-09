#ifndef NOTIFICATION_CLIENT_H
#define NOTIFICATION_CLIENT_H

#include <windows.h>
#include <mmdeviceapi.h>
#include <endpointvolume.h>
#include <functiondiscoverykeys_devpkey.h>

#include "notification_export.h"

#ifdef __cplusplus

extern "C" {
#endif

	HRESULT RegisterNotificationClient(IMMDeviceEnumerator* deviceEnum, IMMNotificationClient** notifClient);
	HRESULT UnregisterNotificationClient(IMMDeviceEnumerator* deviceEnum, IMMNotificationClient* notifClient);

	HRESULT RegisterControlChangeNotify(IAudioEndpointVolume* volume, IAudioEndpointVolumeCallback** volumeCallback);
	HRESULT UnregisterControlChangeNotify(IAudioEndpointVolume* volume, IAudioEndpointVolumeCallback* volumeCallback);

#ifdef __cplusplus
}
#endif

#endif // NOTIFICATION_CLIENT_H
