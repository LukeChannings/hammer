import React, { useState } from "https://cdn.skypack.dev/react";
import { render } from "https://cdn.skypack.dev/react-dom";
import "./test.css";

const App = () => {
	const [n, setN] = useState<number>(0);

	return (
		<p className="Test">Count {n} <button onClick={() => setN(n + 1)}>Inc!</button></p>
	)
}

render(<App />, document.getElementById('root'))
