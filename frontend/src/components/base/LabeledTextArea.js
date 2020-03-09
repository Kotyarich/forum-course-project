import React from "react";
import './LabeledTextArea.css'
import TextArea from "./TextArea";

const LabeledTextArea = (props) => {
  return (
    <div className="form-group">
      <label htmlFor={props.name} className="form-label">{props.title}</label>
      <TextArea value={props.value}
                error={props.error}
                placeholder={props.placeholder}
                onChange={props.onChange}
                name={props.name}/>
    </div>
  )
};

export default LabeledTextArea;