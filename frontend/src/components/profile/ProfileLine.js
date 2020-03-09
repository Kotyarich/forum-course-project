import React from "react";
import Input from '../base/Input'
import './ProfileLine.css'
import TextArea from "../base/TextArea";

const ProfileLine = (props) => {
  const className = "form-input" + (props.error ? " form-input-error" : "");
  return (
    <div className="profile-line">
      <div className={"profile-line__name"}>{props.title}</div>
      {props.isChanging ?
        props.type === 'input' ? <Input
          className={className}
          id={props.name}
          name={props.name}
          value={props.value}
          error={props.error}
          onChange={props.onChange}
          placeholder={''}
        /> : <TextArea value={props.value}
                       error={props.error}
                       placeholder={''}
                       onChange={props.onChange}
                       name={props.name}/> :
        <div className={"profile-line__value"}>{props.value}</div>}
    </div>
  )
};

export default ProfileLine;