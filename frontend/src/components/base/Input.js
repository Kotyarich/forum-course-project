import React from "react";
import './Input.css'

const Input = (props) => {
  const className = "form-input" + (props.error ? " form-input-error" : "");
  return (
      <div className={"form-input__container"}>
        <input
          className={className}
          id={props.name}
          type={props.type || 'text'}
          name={props.name}
          value={props.value}
          onChange={(e) => props.onChange(e.target.name, e.target.value)}
          placeholder={props.placeholder}
        />
        {props.error ? <div className="form-input__error">
          {props.error}
        </div> : null}
      </div>
  )
};

export default Input;