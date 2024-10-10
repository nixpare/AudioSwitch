import { AudioService, WindowService } from "../../bindings/github.com/nixpare/AudioSwitch";
import * as types from "../../bindings/github.com/nixpare/AudioSwitch";
import * as wails from "@wailsio/runtime";

import { render } from "solid-js/web";
import { createSignal, createEffect, For, Show } from "solid-js";
import { createStore, reconcile } from "solid-js/store"

const [appState, setAppState] = createStore(new types.State())

AudioService.GetState()
    .then(state => setAppState(reconcile(state)))
    .catch(err => console.error(err))

wails.Events.On("audio-device-update", (ev) => {
    const newState = ev.data[0];
    setAppState(reconcile(newState, { merge: true }))
});

window.addEventListener('resize', (ev) => {
    wails.Events.Emit({ name: "window-resize" })
})

window.addEventListener('blur', () => {
    document.body.style = '--body-background-alpha: 1'
})
window.addEventListener('focus', () => {
    document.body.removeAttribute('style')
})

const [devices, setDevices] = createStore({})

createEffect(() => {
    Object.entries(appState.Devices).forEach(([key, value]) => {
        if (appState.Prefs[key]) return

        setDevices(key, reconcile({
            ...value,
            pref: false
        }, { merge: true }))
    })

    Object.entries(appState.Prefs).forEach(([key, value]) => {
        setDevices(key, reconcile({
            ...value,
            pref: true
        }, { merge: true }))
    })
})

function DeviceList() {
    return (
        <>
            <For each={Object.values(devices)} >{
                (device) => <Device device={device} />
            }</For>
        </>
    )
}

function SelectedDevice() {
    const [selected, setSelected] = createSignal('')
    const [size, setSize] = createSignal(0)

    createEffect(() => {
        setSelected(devices[appState.Selected]?.Name || 'No Device Selected')

        // for some reason the items does not report the correct final size
        // immediately, so this fixes it
        setTimeout(() => {
            let maxSize = 0
            document.querySelectorAll('[device-list] [device').forEach(el => {
                const btn = el.querySelector('button')
                maxSize = Math.max(maxSize, el.offsetWidth - btn.offsetWidth)
            })
            setSize(maxSize)
        }, 0)
    })

    async function toggleSelected() {
        await AudioService.ToggleSelected();
    }

    return (
        <>
            <h3>Select device:</h3>
            <div class="device" style={`min-width: ${size()}px`}>{selected()}</div>
            <button class="btn" onclick={toggleSelected}>
                <Show when={!appState.Muted}>
                    <svg xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 384 512">{/*<!--!Font Awesome Free 6.6.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license/free Copyright 2024 Fonticons, Inc.-->*/}
                        <path
                            d="M192 0C139 0 96 43 96 96l0 160c0 53 43 96 96 96s96-43 96-96l0-160c0-53-43-96-96-96zM64 216c0-13.3-10.7-24-24-24s-24 10.7-24 24l0 40c0 89.1 66.2 162.7 152 174.4l0 33.6-48 0c-13.3 0-24 10.7-24 24s10.7 24 24 24l72 0 72 0c13.3 0 24-10.7 24-24s-10.7-24-24-24l-48 0 0-33.6c85.8-11.7 152-85.3 152-174.4l0-40c0-13.3-10.7-24-24-24s-24 10.7-24 24l0 40c0 70.7-57.3 128-128 128s-128-57.3-128-128l0-40z" />
                    </svg>
                </Show>
                <Show when={appState.Muted}>
                    <svg class="muted" xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 640 512">{/*<!--!Font Awesome Free 6.6.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license/free Copyright 2024 Fonticons, Inc.-->*/}
                        <path
                            d="M38.8 5.1C28.4-3.1 13.3-1.2 5.1 9.2S-1.2 34.7 9.2 42.9l592 464c10.4 8.2 25.5 6.3 33.7-4.1s6.3-25.5-4.1-33.7L472.1 344.7c15.2-26 23.9-56.3 23.9-88.7l0-40c0-13.3-10.7-24-24-24s-24 10.7-24 24l0 40c0 21.2-5.1 41.1-14.2 58.7L416 300.8 416 96c0-53-43-96-96-96s-96 43-96 96l0 54.3L38.8 5.1zM344 430.4c20.4-2.8 39.7-9.1 57.3-18.2l-43.1-33.9C346.1 382 333.3 384 320 384c-70.7 0-128-57.3-128-128l0-8.7L144.7 210c-.5 1.9-.7 3.9-.7 6l0 40c0 89.1 66.2 162.7 152 174.4l0 33.6-48 0c-13.3 0-24 10.7-24 24s10.7 24 24 24l72 0 72 0c13.3 0 24-10.7 24-24s-10.7-24-24-24l-48 0 0-33.6z" />
                    </svg>
                </Show>
            </button>
        </>
    )
}

function Device(props) {
    async function setDevice() {
        await AudioService.SetDevice(props.device.ID);
    }

    return (
        <li id={props.device.ID} class="btn" device onclick={setDevice}>
            {props.device.Name}
            <PrefButton id={props.device.ID} pref={props.device.pref} />
        </li>
    )
}

function PrefButton(props) {
    async function togglePref(ev) {
        ev.stopPropagation();
        await AudioService.TogglePref(props.id);
    }

    return (
        <button class="btn" onclick={togglePref} pref-button>
            <Show when={!props.pref}>
                <svg xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 576 512">{/*<!--!Font Awesome Free 6.6.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license/free Copyright 2024 Fonticons, Inc.-->*/}
                    <path d="M287.9 0c9.2 0 17.6 5.2 21.6 13.5l68.6 141.3 153.2 22.6c9 1.3 16.5 7.6 19.3 16.3s.5 18.1-5.9 24.5L433.6 328.4l26.2 155.6c1.5 9-2.2 18.1-9.7 23.5s-17.3 6-25.3 1.7l-137-73.2L151 509.1c-8.1 4.3-17.9 3.7-25.3-1.7s-11.2-14.5-9.7-23.5l26.2-155.6L31.1 218.2c-6.5-6.4-8.7-15.9-5.9-24.5s10.3-14.9 19.3-16.3l153.2-22.6L266.3 13.5C270.4 5.2 278.7 0 287.9 0zm0 79L235.4 187.2c-3.5 7.1-10.2 12.1-18.1 13.3L99 217.9 184.9 303c5.5 5.5 8.1 13.3 6.8 21L171.4 443.7l105.2-56.2c7.1-3.8 15.6-3.8 22.6 0l105.2 56.2L384.2 324.1c-1.3-7.7 1.2-15.5 6.8-21l85.9-85.1L358.6 200.5c-7.8-1.2-14.6-6.1-18.1-13.3L287.9 79z" />
                </svg>
            </Show>
            <Show when={props.pref}>
                <svg xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 576 512">{/*<!--!Font Awesome Free 6.6.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license/free Copyright 2024 Fonticons, Inc.-->*/}
                    <path d="M316.9 18C311.6 7 300.4 0 288.1 0s-23.4 7-28.8 18L195 150.3 51.4 171.5c-12 1.8-22 10.2-25.7 21.7s-.7 24.2 7.9 32.7L137.8 329 113.2 474.7c-2 12 3 24.2 12.9 31.3s23 8 33.8 2.3l128.3-68.5 128.3 68.5c10.8 5.7 23.9 4.9 33.8-2.3s14.9-19.3 12.9-31.3L438.5 329 542.7 225.9c8.6-8.5 11.7-21.2 7.9-32.7s-13.7-19.9-25.7-21.7L381.2 150.3 316.9 18z" />
                </svg>
            </Show>
        </button>
    )
}

function HotkeyManager() {
    const [hotkeyConfig, setHotkeyConfig] = createStore(new types.HotkeyConfig())
    const [listening, setListening] = createSignal(true)
    
    WindowService.GetHotkeyConfig()
        .then(config => {
            setHotkeyConfig(reconcile(config, { merge: true }))
        })
        .catch(err => console.error(err))

    function keyDown(ev) {
        if (!listening()) return

        ev.preventDefault()

        setHotkeyConfig(reconcile({
            Shift: ev.shiftKey,
            Ctrl: ev.ctrlKey,
            Alt: ev.altKey,
            Meta: ev.metaKey,
        }, { merge: true }))

        if (ev.location == 0 && ev.key.length == 1) {
            setListening(false)

            setHotkeyConfig(reconcile({
                ...hotkeyConfig,
                Key: ev.key.toLocaleUpperCase(),
                Code: ev.keyCode
            }, { merge: true }))
        }
    }

    async function deleteHotkey() {
        setHotkeyConfig(reconcile(new types.HotkeyConfig(), { merge: true }))
        await WindowService.UnregisterHotkey()
        setListening(true)
    }

    async function saveHotkey() {
        await WindowService.RegisterHotkey(hotkeyConfig);
    }
    
    return (
        <>
            <div class="btn hotkey" tabindex="0" onkeydown={keyDown}>
                <div class={`key ${hotkeyConfig.Shift ? '' : 'hidden'}`} shift-key>SHIFT</div>
                <div class={`key ${hotkeyConfig.Ctrl ? '' : 'hidden'}`} ctrl-key>CTRL</div>
                <div class={`key ${hotkeyConfig.Alt ? '' : 'hidden'}`} alt-key>ALT</div>
                <div class={`key ${hotkeyConfig.Meta ? '' : 'hidden'}`} meta-key>META</div>
                <div class={`key ${hotkeyConfig.Key ? '' : 'hidden'}`} normal-key>{hotkeyConfig.Key}</div>
            </div>
            <button class="btn" onclick={deleteHotkey}>
                <svg xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 448 512">{/*<!--!Font Awesome Free 6.6.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license/free Copyright 2024 Fonticons, Inc.-->*/}
                    <path
                        d="M135.2 17.7L128 32 32 32C14.3 32 0 46.3 0 64S14.3 96 32 96l384 0c17.7 0 32-14.3 32-32s-14.3-32-32-32l-96 0-7.2-14.3C307.4 6.8 296.3 0 284.2 0L163.8 0c-12.1 0-23.2 6.8-28.6 17.7zM416 128L32 128 53.2 467c1.6 25.3 22.6 45 47.9 45l245.8 0c25.3 0 46.3-19.7 47.9-45L416 128z" />
                </svg>
            </button>
            <button class="btn" onclick={saveHotkey}>
                <svg xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 448 512">{/*<!--!Font Awesome Free 6.6.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license/free Copyright 2024 Fonticons, Inc.-->*/}
                    <path
                        d="M48 96l0 320c0 8.8 7.2 16 16 16l320 0c8.8 0 16-7.2 16-16l0-245.5c0-4.2-1.7-8.3-4.7-11.3l33.9-33.9c12 12 18.7 28.3 18.7 45.3L448 416c0 35.3-28.7 64-64 64L64 480c-35.3 0-64-28.7-64-64L0 96C0 60.7 28.7 32 64 32l245.5 0c17 0 33.3 6.7 45.3 18.7l74.5 74.5-33.9 33.9L320.8 84.7c-.3-.3-.5-.5-.8-.8L320 184c0 13.3-10.7 24-24 24l-192 0c-13.3 0-24-10.7-24-24L80 80 64 80c-8.8 0-16 7.2-16 16zm80-16l0 80 144 0 0-80L128 80zm32 240a64 64 0 1 1 128 0 64 64 0 1 1 -128 0z" />
                </svg>
            </button>
        </>
    )
}

function ExitButton() {
    function exit() {
        WindowService.Exit();
    }

    return (
        <button class="btn exit-btn" onclick={exit}>Exit</button>
    )
}

render(() => <SelectedDevice />, document.querySelector('[selected-device]'))
render(() => <DeviceList />, document.querySelector('[device-list]'))
render(() => <HotkeyManager />, document.querySelector('[hotkey-manager]'))
render(() => <ExitButton />, document.querySelector('[exit-button]'))
