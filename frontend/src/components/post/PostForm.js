import React, {Component} from "react";
import {observer} from "mobx-react";
import Button from "../base/Button"
import TextArea from "../base/TextArea";
import LabeledInput from "../base/LabeledInput";
import './PostForm.css'

@observer
class PostForm extends Component {
  handleFormSubmit = (e) => {
    e.preventDefault();
    this.props.onSend();
  };

  render() {
    const {form, onChange} = this.props;
    return (
      <form className={'post-form'}
            id={'post-form'}
            noValidate={true}
            onSubmit={(e) => {this.handleFormSubmit(e)}}>
        <LabeledInput name={'parent'}
                      type={'text'}
                      title={'Answer to post #'}
                      placeholder={''}
                      value={form.fields.parent.value}
                      error={form.fields.parent.error}
                      onChange={onChange}/>
        <TextArea value={form.fields.message.value}
                  error={form.fields.message.error}
                  placeholder={''}
                  onChange={onChange}
                  name={'message'}/>
        <Button title={'Answer'}
                name={'answer'}
                action={(e) => {this.handleFormSubmit(e)}}/>
      </form>
    );
  }
}

export default PostForm;
