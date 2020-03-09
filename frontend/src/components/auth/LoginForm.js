import React, {Component} from "react";
import {observer} from "mobx-react";
import LabeledInput from "../base/LabeledInput"
import Button from "../base/Button"
import './LoginForm.css'

@observer
class LoginForm extends Component {
  handleFormSubmit = (e) => {
    e.preventDefault();

    this.props.loginStore.signIn().then((user) => {
      this.props.userStore.currentUser = user;
      this.props.history.push('/profile/' + user.nickname);
    });
  };

  handleFormCancel(e) {
    e.preventDefault();
    this.props.history.goBack();
  }

  render() {
    const {form, onChange} = this.props;
    return (
      <form className={'login-form'}
            noValidate={true}
            onSubmit={this.handleFormSubmit}>
        <label className={'form-title'}>
          {'Login'}
        </label>
        <hr className={'form__hr'}/>
        <LabeledInput title={'Nickname:'}
                      name={'nickname'}
                      type={'text'}
                      value={form.fields.nickname.value}
                      error={form.fields.nickname.error}
                      placeholder={'Enter your nickname'}
                      onChange={onChange}/>
        <LabeledInput title={'Password:'}
                      name={'password'}
                      type={'password'}
                      value={form.fields.password.value}
                      error={form.fields.password.error}
                      placeholder={'Enter your password'}
                      onChange={onChange}/>
        <hr className={'form__hr'}/>
        <div className={'login-button-group'}>
          <Button title={'Login'}
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

export default LoginForm;