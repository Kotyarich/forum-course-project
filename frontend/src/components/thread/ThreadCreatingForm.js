import {Component} from "react";
import {observer} from "mobx-react";
import LabeledInput from "../base/LabeledInput"
import React from "react";
import Button from "../base/Button";
import LabeledTextArea from "../base/LabeledTextArea";
import './ThreadCreatingForm.css'

@observer
class ThreadCreatingForm extends Component {
  handleFormCreate(e) {
    e.preventDefault();
    this.props.threadStore.createThread(this.props.slug, 
                      this.props.userStore.currentUser.nickname).then(() => {
      this.props.history.push('/forum/' + this.props.slug);
    });
  }

  handleFormCancel(e) {
    e.preventDefault();
    this.props.history.goBack();
  }

  render() {
    const {form, onChange} = this.props;
    console.log(form);
    return (
      <form className={'thread-creating-form'}
            noValidate={true}
            onSubmit={this.handleFormCreate}>
        <label className={'form-title'}>
          {'Create thread'}
        </label>
        <hr className={'form__hr'}/>
        <LabeledInput title={'Thread name:'}
                      name={'threadname'}
                      type={'text'}
                      value={form.fields.threadname.value}
                      error={form.fields.threadname.error}
                      placeholder={'Enter your thread name'}
                      onChange={onChange}/>
        <LabeledTextArea title={'Initial post:'}
                        name={'initialpost'}
                        value={form.fields.initialpost.value}
                        error={form.fields.initialpost.error}
                        placeholder={'Write your post'}
                        onChange={onChange}/>
        <hr className={'form__hr'}/>
        <div className={'button-group'}>
          <Button title={'Create'}
                  disabled={!form.meta.isValid}
                  name={'create'}
                  action={(e) => this.handleFormCreate(e)}/>
          <Button title={'Cancel'}
                  name={'cancel'}
                  action={(e) => this.handleFormCancel(e)}/>
        </div>
      </form>
    );
  }
}

export default ThreadCreatingForm;