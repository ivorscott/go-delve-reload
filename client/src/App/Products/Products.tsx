import React from "react";
import { IProduct } from "./types";
import "./Products.scss";

const Products: React.FC<{ products: IProduct[] }> = ({ products }) => (
  <div className="products">
    {products.map(product => {
      return (
        <div key={product.id} className="products__card">
          <header className="products__card-header">
            <h2 className="products__card-name">{product.name}</h2>
          </header>
          <footer className="products__card-footer">
            <p className="products__card-price">${product.price}</p>
          </footer>
        </div>
      );
    })}
  </div>
);

export default Products;
