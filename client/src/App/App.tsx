import React from "react";
import Header from "./Header";
import Products from "./Products";
import { IProduct } from "./Products/types";
import "./App.scss";

type State = typeof initialState;

const initialState: { products: IProduct[] } = { products: [] };

class App extends React.Component<{}, State> {
  state = initialState;

  componentDidMount = async () => {
    try {
      const response = await fetch(`${process.env.REACT_APP_API_URL}/products`);
      const json = await response.json();
      this.setState(() => ({ products: json }));
    } catch (error) {}
  };

  render() {
    return (
      <div className="app">
        <Header
          title="Kingshot"
          subtitle="Second hand games"
          callToActionText="Sign Up"
        />
        <Products products={this.state.products} />
      </div>
    );
  }
}

export default App;
