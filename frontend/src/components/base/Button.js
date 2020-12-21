import React from "react";
import './Button.css'

const Button = (props) => {
  return (
    <button
      className={props.name + '-button button'}
      disabled={props.disabled}
      onClick={props.action}>
      {props.title? props.title : props.children}
    </button>)
};

export default Button;