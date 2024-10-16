import { AudioService, WindowService } from "../../bindings/github.com/nixpare/AudioSwitch";
import * as wails from "@wailsio/runtime";

import { render } from "solid-js/web";
import { createEffect, createSignal, onMount, Show } from "solid-js";

const [muted, setMuted] = createSignal(false)

AudioService.GetState()
	.then(state => setMuted(state.Muted))
	.catch(err => console.error(err))

wails.Events.On("audio-device-update", (ev) => {
	setMuted(ev.data[0].Muted)
});

document.addEventListener('contextmenu', async () => {
	await WindowService.CreateWindow();
})

const background = document.querySelector('body > .background')
let backgroundEffectTimeout;

createEffect(() => {
	muted()
	background.style = '--body-background-alpha: 1'
	
	if (backgroundEffectTimeout) {
		window.clearTimeout(backgroundEffectTimeout)
	}
	
	backgroundEffectTimeout = window.setTimeout(() => {
		background.removeAttribute('style')
		backgroundEffectTimeout = undefined;
	}, 1200)
})

function MuteButton() {
	async function toggleSelected() {
		await AudioService.ToggleSelected();
	}

	onMount(async () => {
		await resizeWindow()
	})
	
	return (
		<button class="btn" onclick={toggleSelected}>
			<Show when={!muted()}>
				<svg xmlns="http://www.w3.org/2000/svg"
					viewBox="0 0 384 512">{/*<!--!Font Awesome Free 6.6.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license/free Copyright 2024 Fonticons, Inc.-->*/}
					<path
						d="M192 0C139 0 96 43 96 96l0 160c0 53 43 96 96 96s96-43 96-96l0-160c0-53-43-96-96-96zM64 216c0-13.3-10.7-24-24-24s-24 10.7-24 24l0 40c0 89.1 66.2 162.7 152 174.4l0 33.6-48 0c-13.3 0-24 10.7-24 24s10.7 24 24 24l72 0 72 0c13.3 0 24-10.7 24-24s-10.7-24-24-24l-48 0 0-33.6c85.8-11.7 152-85.3 152-174.4l0-40c0-13.3-10.7-24-24-24s-24 10.7-24 24l0 40c0 70.7-57.3 128-128 128s-128-57.3-128-128l0-40z" />
				</svg>
			</Show>
			<Show when={muted()}>
				<svg class="muted" xmlns="http://www.w3.org/2000/svg"
					viewBox="0 0 640 512">{/*<!--!Font Awesome Free 6.6.0 by @fontawesome - https://fontawesome.com License - https://fontawesome.com/license/free Copyright 2024 Fonticons, Inc.-->*/}
					<path
						d="M38.8 5.1C28.4-3.1 13.3-1.2 5.1 9.2S-1.2 34.7 9.2 42.9l592 464c10.4 8.2 25.5 6.3 33.7-4.1s6.3-25.5-4.1-33.7L472.1 344.7c15.2-26 23.9-56.3 23.9-88.7l0-40c0-13.3-10.7-24-24-24s-24 10.7-24 24l0 40c0 21.2-5.1 41.1-14.2 58.7L416 300.8 416 96c0-53-43-96-96-96s-96 43-96 96l0 54.3L38.8 5.1zM344 430.4c20.4-2.8 39.7-9.1 57.3-18.2l-43.1-33.9C346.1 382 333.3 384 320 384c-70.7 0-128-57.3-128-128l0-8.7L144.7 210c-.5 1.9-.7 3.9-.7 6l0 40c0 89.1 66.2 162.7 152 174.4l0 33.6-48 0c-13.3 0-24 10.7-24 24s10.7 24 24 24l72 0 72 0c13.3 0 24-10.7 24-24s-10.7-24-24-24l-48 0 0-33.6z" />
				</svg>
			</Show>
		</button>
	)
}

render(() => <MuteButton />, document.querySelector('[mute-button]'))

async function resizeWindow() {
	const width = document.body.offsetWidth;
	const height = document.body.offsetHeight;
	const size = await wails.Window.Size();

	if (size.width == width && size.height == height) {
		return;
	}

	await wails.Window.SetSize(width, height);
	wails.Events.Emit({ name: "overlay-resize" });
}
