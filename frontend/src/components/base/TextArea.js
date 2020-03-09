import React from "react";
import './TextArea.css'

const TextArea = (props) => {
  const className = "form-textarea" + (props.error ? " form-textarea-error" : "");
  return (
    <div className={"form-textarea__container"}>
        <textarea
          className={className}
          id={props.name}
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

export default TextArea;