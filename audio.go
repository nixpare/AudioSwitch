package main

/*
#include "winaudio_wrapper.h"
#include "notification.h"
#include "winhelper.h"
*/
import "C"
import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
)

type SaveState struct {
	Prefs    map[string]*Device
	Selected string
	Muted    bool
}

type State struct {
	SaveState

	Devices map[string]*Device
}

type AudioService struct {
	State

	deviceEnum  *C.IMMDeviceEnumerator
	notifClient *C.IMMNotificationClient

	running bool
	m       sync.Mutex
}

var (
	ErrAudioServiceNotRunning = errors.New("audio service not running")
	ErrDeviceNotFound         = errors.New("device not found")
)

func initCOMLibraryMultithreaded() error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var wg sync.WaitGroup
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()

		hr := C.CoInitializeEx(nil, C.COINIT_MULTITHREADED)
		if int32(hr) < 0 {
			err = fmt.Errorf("COM library init failed with code 0x%x", uint32(hr))
		}
	}()

	wg.Wait()
	return err
}

func newAudioService() *AudioService {
	return &AudioService{
		State: State{
			Devices: make(map[string]*Device),
		},
	}
}

func (s *AudioService) Start() error {
	s.m.Lock()
	defer s.m.Unlock()

	if s.running {
		return nil
	}

	err := s.loadSaveData()
	if err != nil {
		return err
	}

	err = initCOMLibraryMultithreaded()
	if err != nil {
		return err
	}

	// Create IMMDeviceEnumerator instance
	hr := C.CreateInstance(&s.deviceEnum)
	if hr < 0 {
		return fmt.Errorf("device enumerator create: 0x%x", uint32(hr))
	}

	err = s.updateDeviceList()
	if err != nil {
		return fmt.Errorf("device list update: %w", err)
	}

	// Start listening for device events
	err = s.listenForDeviceEvents()
	if err != nil {
		return err
	}

	s.running = true
	return nil
}

func (s *AudioService) Stop() error {
	s.m.Lock()
	defer s.m.Unlock()

	var errs []error

	if !s.running {
		return nil
	}

	if err := s.updateSaveData(); err != nil {
		return err
	}

	for _, device := range s.Devices {
		device.release()
	}
	clear(s.Devices)

	// Unregister IMMNotificationClient callbacks
	if hr := C.UnregisterNotificationClient(s.deviceEnum, s.notifClient); hr < 0 {
		errs = append(errs, fmt.Errorf("device enumerator create: 0x%x", uint32(hr)))
	}
	// Release IMMDeviceEnumerator instance
	C.IMMDeviceEnumerator_Release(s.deviceEnum)
	C.CoUninitialize()

	s.running = false
	return errors.Join(errs...)
}

func (s *AudioService) GetState() (State, error) {
	s.m.Lock()
	defer s.m.Unlock()

	if !s.running {
		return State{}, ErrAudioServiceNotRunning
	}

	return s.State, nil
}

func (s *AudioService) updateDeviceList() error {
	devices, err := s.getDeviceCollection()
	if err != nil {
		return err
	}

	check := make(map[string]bool)
	for id := range s.Devices {
		check[id] = false
	}

	// Iterate over each device and retrieve its properties
	for _, device := range devices {
		check[device.ID] = true

		if oldDev, ok := s.Devices[device.ID]; ok {
			oldDev.copyStateFrom(device)
			device.release()
		} else {
			s.Devices[device.ID] = device
		}
	}

	for id, found := range check {
		if !found {
			s.Devices[id].release()
			delete(s.Devices, id)
		}
	}

	if device, ok := s.Devices[s.Selected]; ok {
		err = device.activate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *AudioService) SetDevice(id string) error {
	s.m.Lock()
	defer s.m.Unlock()

	if !s.running {
		return ErrAudioServiceNotRunning
	}

	for deviceID, device := range s.Devices {
		if deviceID != id {
			continue
		}

		if id == s.Selected {
			return s.updateFrontend(false)
		}

		// TODO: handle error
		if oldDevice, ok := s.Devices[s.Selected]; ok {
			oldDevice.deactivate()
		}

		s.Selected = id
		err := device.activate()
		if err != nil {
			return err
		}

		return s.updateFrontend(false)
	}

	for deviceID := range s.Prefs {
		if deviceID != id {
			continue
		}

		if id == s.Selected {
			return s.updateFrontend(false)
		}

		// TODO: handle error
		if oldDevice, ok := s.Devices[s.Selected]; ok {
			oldDevice.deactivate()
		}

		s.Selected = id
		return s.updateFrontend(false)
	}

	return ErrDeviceNotFound
}

func (s *AudioService) TogglePref(id string) error {
	s.m.Lock()
	defer s.m.Unlock()

	if !s.running {
		return ErrAudioServiceNotRunning
	}

	_, ok := s.Prefs[id]
	if ok {
		delete(s.Prefs, id)
		return s.updateFrontend(false)
	}

	device, ok := s.Devices[id]
	if !ok {
		return ErrDeviceNotFound
	}

	s.Prefs[id] = device
	return s.updateFrontend(false)
}

func (s *AudioService) ToggleSelected() error {
	s.m.Lock()
	defer s.m.Unlock()

	if !s.running {
		return ErrAudioServiceNotRunning
	}

	return s.setPrefMute(!s.Muted)
}

func (s *AudioService) setPrefMute(muted bool) error {
	device, ok := s.Devices[s.Selected]
	if ok {
		err := device.setMuted(muted)
		if err != nil {
			return err
		}

		s.Muted = muted
		return nil
	}

	_, ok = s.Prefs[s.Selected]
	if ok {
		return nil
	}
	
	return ErrDeviceNotFound
}

func (s *AudioService) loadSaveData() error {
	saveFile, err := os.OpenFile(audioSaveFilePath, os.O_RDONLY|os.O_CREATE, 0660)
	if err != nil {
		return fmt.Errorf("save file open: %w", err)
	}
	defer saveFile.Close()

	saveData, err := io.ReadAll(saveFile)
	if err != nil {
		return fmt.Errorf("save file read: %w", err)
	}

	if len(saveData) == 0 {
		s.State.SaveState.Prefs = make(map[string]*Device)
		return nil
	}

	err = json.Unmarshal(saveData, &s.State.SaveState)
	if err != nil {
		return fmt.Errorf("save data decode: %w", err)
	}

	return nil
}

func (s *AudioService) updateSaveData() error {
	saveData, err := json.MarshalIndent(s.State.SaveState, "", "\t")
	if err != nil {
		return fmt.Errorf("save data encode: %w", err)
	}

	saveFile, err := os.OpenFile(audioSaveFilePath, os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		return fmt.Errorf("save file open: %w", err)
	}
	defer saveFile.Close()

	_, err = saveFile.Write(saveData)
	if err != nil {
		return fmt.Errorf("save file write: %w", err)
	}

	return nil
}

func (s *AudioService) listenForDeviceEvents() error {
	if hr := C.RegisterNotificationClient(s.deviceEnum, &s.notifClient); hr < 0 {
		return fmt.Errorf("audio notification registration: 0x%x", uint32(hr))
	}
	return nil
}

func (s *AudioService) updateFrontend(regenerateList bool) error {
	if regenerateList {
		err := s.updateDeviceList()
		if err != nil {
			log.Printf("audio service: device list update error: %v\n", err)
			return err
		}
	}

	app.EmitEvent("audio-device-update", s.State)
	return nil
}

func (s *AudioService) getDeviceCollection() ([]*Device, error) {
	// Enumerate audio endpoints (eRender for playback devices, eCapture for recording devices)
	var deviceCollection *C.IMMDeviceCollection
	if hr := C.IMMDeviceEnumerator_EnumAudioEndpoints(s.deviceEnum, &deviceCollection); hr < 0 {
		return nil, fmt.Errorf("audio device collection: 0x%x", uint32(hr))
	}
	defer C.IMMDeviceCollection_Release(deviceCollection)

	// Get the number of audio devices
	var count C.uint
	if hr := C.IMMDeviceCollection_GetCount(deviceCollection, &count); hr < 0 {
		return nil, fmt.Errorf("audio device collection count: 0x%x", uint32(hr))
	}

	devices := make([]*Device, 0, count)

	// Iterate over each device and retrieve its properties
	for i := range count {
		var immDevice *C.IMMDevice
		if hr := C.IMMDeviceCollection_Item(deviceCollection, i, &immDevice); hr < 0 {
			log.Printf("audio service: device collection item error: 0x%x\n", hr)
			continue
		}

		device, err := newDevice(immDevice)
		if err != nil {
			log.Printf("audio service: device error: %v\n", err)
			continue
		}

		devices = append(devices, device)
	}

	return devices, nil
}
