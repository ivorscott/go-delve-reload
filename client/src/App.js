import React from "react";
import logo from "./logo.svg";
import "./App.css";

class App extends React.Component {
  state = {
    products: []
  };

  componentDidMount = async () => {
    try {
      const response = await fetch(`http://localhost:4000/products`);
      const json = await response.json();
      this.setState(() => ({ products: json }));
    } catch (error) {}
  };

  render() {
    const { products } = this.state;
    return (
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <p>
            Edit <code>src/App.tsx</code> and save to reload.
          </p>
          <a
            className="App-link"
            href="https://reactjs.org"
            target="_blank"
            rel="noopener noreferrer"
          >
            Learn React
          </a>
        </header>
        <footer style={{ textAlign: "left" }}>
          <pre>
            <code>
              {products.length > 0 && JSON.stringify(products, null, 4)}
            </code>
          </pre>
        </footer>
      </div>
    );
  }
}

export default App;
