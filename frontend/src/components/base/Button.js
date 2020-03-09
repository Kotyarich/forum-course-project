import React from "react";
import './Button.css'

const Button = (props) => {
  return (
    <button
      className={props.name + '-button button'}
      disabled={props.disabled}
      onClick={props.action}>
      {props.title}
    </button>)
};

export default Button;