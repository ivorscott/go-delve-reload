import React from "react";
import { IProduct } from "./types";
import "./Products.scss";

const Products: React.FC<{ products: IProduct[] }> = ({ products }) => (
  <main className="products-feature">
    <div className="products">
      {products.map(product => {
        return (
          <div key={product.id} className="products__card">
            <b className="products__card-name">{product.name}</b>
            <p className="products__card-price">{product.price}</p>
          </div>
        );
      })}
    </div>
  </main>
);

export default Products;
