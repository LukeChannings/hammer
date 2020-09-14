interface LiveReloadMessage {
	path: string;
	content: string;
	fileType: string;
}

new EventSource(globalThis.__LIVE_RELOAD_PATH__)
	.addEventListener(
		"message",
		e => {
			if (e.data === "heartbeat") return;
			const msg: LiveReloadMessage = JSON.parse(e.data);
			if (msg.fileType === "css") {
				const style = document.querySelector(`[data-path="${msg.path}"]`);
				if (style) style.innerHTML = msg.content;
			} else {
				location.reload();
			}
		}
	);
