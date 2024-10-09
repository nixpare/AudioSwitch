package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/nixpare/broadcaster"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"github.com/wailsapp/wails/v3/pkg/w32"
	"golang.design/x/hotkey"
)

type WindowState struct {
	X, Y int
	Width, Height int
}

type HotkeyConfig struct {
	Shift bool
	Ctrl bool
	Alt bool
	Meta bool
	Key string
	Code uint16
}

type WindowService struct {
	window  *application.WebviewWindow
	overlay *application.WebviewWindow

	hotkeyBr *broadcaster.Broadcaster[chan <-error]
	hotkeyRegistered bool

	WindowState  WindowState  `json:"window"`
	OverlayState WindowState  `json:"overlay"`
	HotkeyConfig HotkeyConfig `json:"hotkey"`
}

func newWindowService() (*WindowService, error) {
	w := &WindowService{
		hotkeyBr: broadcaster.NewBroadcaster[chan <-error](),
	}

	err := w.loadSaveData()
	if err != nil {
		return nil, err
	}

	err = w.RegisterHotkey(nil)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (w *WindowService) loadSaveData() error {
	saveFile, err := os.OpenFile(windowSaveFilePath, os.O_RDONLY|os.O_CREATE, 0660)
	if err != nil {
		return fmt.Errorf("save file open: %w", err)
	}
	defer saveFile.Close()

	saveData, err := io.ReadAll(saveFile)
	if err != nil {
		return fmt.Errorf("save file read: %w", err)
	}

	if len(saveData) == 0 {
		return nil
	}

	err = json.Unmarshal(saveData, w)
	if err != nil {
		return fmt.Errorf("save data decode: %w", err)
	}

	return nil
}

func (w *WindowService) Close() error {
	err := w.updateSaveData()
	if err != nil {
		return err
	}

	err = w.UnregisterHotkey()
	if err != nil {
		return err
	}

	return nil
}

func (w *WindowService) updateSaveData() error {
	saveData, err := json.MarshalIndent(w, "", "\t")
	if err != nil {
		return fmt.Errorf("save data encode: %w", err)
	}

	saveFile, err := os.OpenFile(windowSaveFilePath, os.O_WRONLY | os.O_TRUNC, 0)
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

func (w *WindowService) CreateWindow() {
	if w.window != nil {
		w.window.Show()
		w.window.UnMinimise()
		return
	}

	w.window = createWindowOptions(&w.WindowState, application.WebviewWindowOptions{
		Title:            "AudioSwitch Dashboard",
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	w.window.OnWindowEvent(events.Windows.WindowClose, func(event *application.WindowEvent) {
		w.window = nil
	})

	app.OnEvent("window-resize", func(event *application.CustomEvent) {
		w.WindowState.Width, w.WindowState.Height = w.window.Size()
	})
}

func (w *WindowService) CreateOverlay() {
	if w.overlay != nil {
		w.overlay.Show()
		return
	}

	w.overlay = createWindowOptions(&w.OverlayState, application.WebviewWindowOptions{
		Title:            "AudioSwitch Overlay",
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/overlay.html",

		Frameless:     true,
		DisableResize: true,

		Windows: application.WindowsWindow{
			ExStyle: w32.WS_EX_TOOLWINDOW|w32.WS_EX_TOPMOST,
		},

		DefaultContextMenuDisabled: true,
	})

	w.overlay.OnWindowEvent(events.Windows.WindowClose, func(event *application.WindowEvent) {
		w.overlay = nil
	})

	app.OnEvent("overlay-resize", func(event *application.CustomEvent) {
		w.OverlayState.Width, w.OverlayState.Height = w.overlay.Size()
	})
}

func (w *WindowService) Exit() {
	w.window.Close()
	w.overlay.Close()
}

func createWindowOptions(state *WindowState, options application.WebviewWindowOptions) *application.WebviewWindow {
	if state.X == 0 && state.Y == 0 {
		options.Centered = true
		options.X, options.Y = 0, 0
	} else {
		options.Centered = false
		options.X, options.Y = state.X, state.Y
	}

	if state.Width != 0 && state.Height != 0 {
		options.Width, options.Height = state.Width, state.Height
	} else {
		options.Width, options.Height = 0, 0
	}

	window := app.NewWebviewWindowWithOptions(options)

	window.OnWindowEvent(events.Common.WindowDidMove, func(event *application.WindowEvent) {
		updateWindowState(window, state)
	})
	window.OnWindowEvent(events.Common.WindowDidResize, func(event *application.WindowEvent) {
		updateWindowState(window, state)
		println("resize")
	})

	return window
}

func updateWindowState(window *application.WebviewWindow, state *WindowState) {
	state.X, state.Y = window.Position()
	state.Width, state.Height = window.Width(), window.Height()
}

func (w *WindowService) GetHotkeyConfig() HotkeyConfig {
	return w.HotkeyConfig
}

func (w *WindowService) RegisterHotkey(data any) error {
	if w.hotkeyRegistered {
		return nil
	}

	if data != nil {
		hotkeyConfig := data.(map[string]any)
		w.HotkeyConfig = HotkeyConfig{
			Shift: hotkeyConfig["Shift"].(bool),
			Ctrl: hotkeyConfig["Ctrl"].(bool),
			Alt: hotkeyConfig["Alt"].(bool),
			Meta: hotkeyConfig["Meta"].(bool),
			Key: hotkeyConfig["Key"].(string),
			Code: uint16(hotkeyConfig["Code"].(float64)),
		}
	}

	if w.HotkeyConfig.Key == "" {
		return nil
	}

	log.Printf("New hotkey registered: %+v\n", w.HotkeyConfig)

	var modifiers []hotkey.Modifier
	if w.HotkeyConfig.Shift {
		modifiers = append(modifiers, hotkey.ModShift)
	}
	if w.HotkeyConfig.Ctrl {
		modifiers = append(modifiers, hotkey.ModCtrl)
	}
	if w.HotkeyConfig.Alt {
		modifiers = append(modifiers, hotkey.ModAlt)
	}
	if w.HotkeyConfig.Meta {
		modifiers = append(modifiers, hotkey.ModWin)
	}

	hk := hotkey.New(modifiers, hotkey.Key(w.HotkeyConfig.Code))
	
	err := hk.Register()
	if err != nil {
		return err
	}

	listener := w.hotkeyBr.Register(0)
	var resultCh chan <-error

	go func() {
		defer listener.Unregister()

	loop:
		for {
			select {
			case <-hk.Keydown():
				err := audioService.ToggleSelected()
				if err != nil {
					log.Printf("toggle selected error: %v\n", err)
				}
			case resultCh = <-listener.Ch():
				break loop
			}
		}

		err := hk.Unregister()
		if err != nil {
			resultCh <- fmt.Errorf("failed to unregister hotkey: %w", err)
		}

		resultCh <- nil
	}()

	w.hotkeyRegistered = true
	return nil
}

func (w *WindowService) UnregisterHotkey() error {
	if !w.hotkeyRegistered {
		return nil
	}

	resultCh := make(chan error)
	defer close(resultCh)

	w.hotkeyBr.Send(resultCh)
	err := <-resultCh
	if err != nil {
		return err
	}

	w.HotkeyConfig = HotkeyConfig{}
	w.hotkeyRegistered = false
	return nil
}
