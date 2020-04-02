import React from "react";
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
        <header className="header">
          <div className="header__logo-box">
            <img alt="logo" className="header__logo" src="/logo-white.png" />
          </div>

          <div className="header__text-box">
            <span className="header__text-box-title">Kingshot</span>
            <span className="header__text-box-subtitle">Second hand games</span>

            <a
              href="/#"
              className="header__cta header__cta--white header__cta--animated"
            >
              Sign Up
            </a>
          </div>
        </header>

        <div className="main">
          <div className="main__pitch">
            <h3 className="main__pitch-title">The Best Prices</h3>

            <p className="main__pitch-discount">
              Get 20% Off Any 2nd Purchase!
            </p>

            <div>
              <Products products={this.state.products} />
            </div>
          </div>
        </div>
      </div>
    );
  }
}

export default App;
