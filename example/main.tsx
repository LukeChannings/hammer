import React, { useState } from "https://cdn.skypack.dev/react";
import { render } from "https://cdn.skypack.dev/react-dom";
import { test } from "./test.css";

const App = () => {
	const [n, setN] = useState<number>(0);

	return (
		<p className={test}>
			Count {n} <button onClick={() => setN(n + 1)}>Inc!</button>
		</p>
	);
}

render(<App />, document.getElementById('root'))
