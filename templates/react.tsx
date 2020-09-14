import React, { useState } from "https://cdn.skypack.dev/react@latest";
import ReactDOM from "https://cdn.skypack.dev/react-dom@latest";

const App = () => {
	const [count, setCount] = useState(0);

	return (
		<div>
			<p>By my count it's {count}!</p>
			<button onClick={() => setCount(count + 1)}>No!</button>
		</div>
	);
};

ReactDOM.render(<App />, document.getElementById("mount"));
