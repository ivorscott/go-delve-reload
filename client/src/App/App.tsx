import React from "react";
import Header from "./Header";
import Products from "./Products";
import "./App.scss";

class App extends React.Component<any, any> {
  state = {
    products: []
  };

  componentDidMount = async () => {
    console.log(process.env);
    try {
      const response = await fetch(`${process.env.REACT_APP_API_URL}/products`);
      const json = await response.json();
      this.setState(() => ({ products: json }));
    } catch (error) {}
  };

  render() {
    const { products } = this.state;
    return (
      <div className="app">
        <Header
          title="Kingshot"
          subtitle="Second hand games"
          callToActionText="Sign Up"
        />
        <Products products={products} />
      </div>
    );
  }
}

export default App;
