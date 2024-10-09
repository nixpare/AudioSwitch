#include "notification.h"

class NotificationClient : public IMMNotificationClient {
private:
	LONG _cRef;  // Conteggio dei riferimenti per la gestione del ciclo di vita dell'oggetto

public:
	// Costruttore
	NotificationClient() : _cRef(1) {}

	// Implementazione di IUnknown
	ULONG STDMETHODCALLTYPE AddRef() {
		return InterlockedIncrement(&_cRef);
	}

	ULONG STDMETHODCALLTYPE Release() {
		ULONG ulRef = InterlockedDecrement(&_cRef);
		if (0 == ulRef) {
			delete this;
		}
		return ulRef;
	}

	HRESULT STDMETHODCALLTYPE QueryInterface(REFIID riid, VOID** ppvInterface) {
		if (IID_IUnknown == riid || __uuidof(IMMNotificationClient) == riid) {
			AddRef();
			*ppvInterface = static_cast<IMMNotificationClient*>(this);
			return S_OK;
		}
		*ppvInterface = NULL;
		return E_NOINTERFACE;
	}

	// Implementazione dei metodi IMMNotificationClient
	HRESULT STDMETHODCALLTYPE OnDeviceStateChanged(LPCWSTR pwstrDeviceId, DWORD dwNewState) {
		return OnDeviceStateChangedCallback(pwstrDeviceId, dwNewState);
	}

	HRESULT STDMETHODCALLTYPE OnDeviceAdded(LPCWSTR pwstrDeviceId) {
		return OnDeviceAddedCallback(pwstrDeviceId);
	}

	HRESULT STDMETHODCALLTYPE OnDeviceRemoved(LPCWSTR pwstrDeviceId) {
		return OnDeviceRemovedCallback(pwstrDeviceId);
	}

	HRESULT STDMETHODCALLTYPE OnDefaultDeviceChanged(EDataFlow flow, ERole role, LPCWSTR pwstrDefaultDeviceId) {
		return OnDefaultDeviceChangedCallback(flow, role, pwstrDefaultDeviceId);
	}

	HRESULT STDMETHODCALLTYPE OnPropertyValueChanged(LPCWSTR pwstrDeviceId, const PROPERTYKEY key) {
		return OnPropertyValueChangedCallback(pwstrDeviceId, key);
	}
};

class EndpointVolumeCallback : public IAudioEndpointVolumeCallback {
private:
	LONG _cRef;  // Conteggio dei riferimenti per la gestione del ciclo di vita dell'oggetto

public:
	// Costruttore
	EndpointVolumeCallback() : _cRef(1) {}

	// Implementazione di IUnknown
	ULONG STDMETHODCALLTYPE AddRef() {
		return InterlockedIncrement(&_cRef);
	}

	ULONG STDMETHODCALLTYPE Release() {
		ULONG ulRef = InterlockedDecrement(&_cRef);
		if (0 == ulRef) {
			delete this;
		}
		return ulRef;
	}

	HRESULT STDMETHODCALLTYPE QueryInterface(REFIID riid, VOID** ppvInterface) {
		if (IID_IUnknown == riid || __uuidof(IAudioEndpointVolumeCallback) == riid) {
			AddRef();
			*ppvInterface = static_cast<IAudioEndpointVolumeCallback*>(this);
			return S_OK;
		}
		*ppvInterface = NULL;
		return E_NOINTERFACE;
	}

	// Implementazione dei metodi IAudioEndpointVolumeCallback
	HRESULT STDMETHODCALLTYPE OnNotify(PAUDIO_VOLUME_NOTIFICATION_DATA pNotify) {
		return OnEndpointVolumeChangeNotify(pNotify);
	}
};

extern "C" {

	HRESULT RegisterNotificationClient(IMMDeviceEnumerator* deviceEnum, IMMNotificationClient** notifClient) {
		*notifClient = new NotificationClient();
		return deviceEnum->RegisterEndpointNotificationCallback(*notifClient);
	}

	HRESULT UnregisterNotificationClient(IMMDeviceEnumerator* deviceEnum, IMMNotificationClient* notifClient) {
		HRESULT hr = deviceEnum->UnregisterEndpointNotificationCallback(notifClient);
		notifClient->Release();
		return hr;
	}

	HRESULT RegisterControlChangeNotify(IAudioEndpointVolume* volume, IAudioEndpointVolumeCallback** volumeCallback) {
		*volumeCallback = new EndpointVolumeCallback();
		return volume->RegisterControlChangeNotify(*volumeCallback);
	}

	HRESULT UnregisterControlChangeNotify(IAudioEndpointVolume* volume, IAudioEndpointVolumeCallback* volumeCallback) {
		HRESULT hr = volume->UnregisterControlChangeNotify(volumeCallback);
		volumeCallback->Release();
		return hr;
	}

}
