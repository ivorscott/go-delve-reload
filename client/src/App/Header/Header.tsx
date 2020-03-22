import React from "react";

export const Header: React.FC<{
  title: string;
  subtitle: string;
  callToActionText: string;
}> = ({ title, subtitle, callToActionText }) => (
  <header className="header">
    <div className="header__logo-box">
      <img
        alt="logo"
        className="header__logo"
        src="https://github.com/jonasschmedtmann/advanced-css-course/blob/master/Natours/starter/img/logo-white.png?raw=true"
      />
    </div>
    <div className="header__text-box">
      <h1 className="heading-primary">
        <span className="heading-primary--main">{title}</span>
        <span className="heading-primary--sub">{subtitle}</span>
      </h1>

      <a href="/#" className="btn btn--white btn--animated">
        {callToActionText}
      </a>
    </div>
  </header>
);
