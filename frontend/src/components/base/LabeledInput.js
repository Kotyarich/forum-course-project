import React from "react";
import Input from './Input'
import './LabeledInput.css'

const LabeledInput = (props) => {
  const className = "form-input" + (props.error ? " form-input-error" : "");
  return (
    <div className="form-group">
      <label htmlFor={props.name} className="form-label">{props.title}</label>
      <Input
        className={className}
        id={props.name}
        type={props.type || 'text'}
        name={props.name}
        value={props.value}
        error={props.error}
        onChange={props.onChange}
        placeholder={props.placeholder}
      />
    </div>
  )
};

export default LabeledInput;