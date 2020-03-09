import {Component} from "react";
import {observer} from "mobx-react";
import LabeledInput from "../base/LabeledInput"
import React from "react";
import Button from "../base/Button";
import LabeledTextArea from "../base/LabeledTextArea";
import './RegistrationForm.css'

@observer
class RegistrationForm extends Component {
  handleFormSubmit(e) {
    e.preventDefault();
    this.props.registrationStore.signUp().then(() => {
      const nickname = this.props.form.fields.nickname.value;
      this.props.userStore.currentUser = {
        nickname: nickname,
      };
      this.props.history.push('/profile/' + nickname);
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
      <form className={'registration-form'}
            noValidate={true}
            onSubmit={this.handleFormSubmit}>
        <label className={'form-title'}>
          {'Registration'}
        </label>
        <hr className={'form__hr'}/>
        <LabeledInput title={'Nickname:'}
                      name={'nickname'}
                      type={'text'}
                      value={form.fields.nickname.value}
                      error={form.fields.nickname.error}
                      placeholder={'Enter your nickname'}
                      onChange={onChange}/>
        <LabeledInput title={'Full Name:'}
                      name={'fullName'}
                      type={'text'}
                      value={form.fields.fullName.value}
                      error={form.fields.fullName.error}
                      placeholder={'Enter your full name'}
                      onChange={onChange}/>
        <LabeledInput title={'Email:'}
                      name={'email'}
                      type={'email'}
                      value={form.fields.email.value}
                      error={form.fields.email.error}
                      placeholder={'Enter your email'}
                      onChange={onChange}/>
        <LabeledInput title={'Password:'}
                      name={'password'}
                      type={'password'}
                      value={form.fields.password.value}
                      error={form.fields.password.error}
                      placeholder={'Enter your password'}
                      onChange={onChange}/>
        <LabeledInput title={'Password submission:'}
                      name={'passwordSubmission'}
                      type={'password'}
                      value={form.fields.passwordSubmission.value}
                      error={form.fields.passwordSubmission.error}
                      placeholder={'Enter your password once more'}
                      onChange={onChange}/>
        <LabeledTextArea title={'About:'}
                         name={'about'}
                         value={form.fields.about.value}
                         error={form.fields.about.error}
                         placeholder={'Write something about you'}
                         onChange={onChange}/>
        <hr className={'form__hr'}/>
        <div className={'button-group'}>
          <Button title={'Register'}
                  disabled={!form.meta.isValid}
                  name={'submit'}
                  action={(e) => this.handleFormSubmit(e)}/>
          <Button title={'Cancel'}
                  name={'cancel'}
                  action={(e) => this.handleFormCancel(e)}/>
        </div>
      </form>
    );
  }
}

export default RegistrationForm;