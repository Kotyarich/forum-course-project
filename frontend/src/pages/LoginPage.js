import React from "react";
import LoginForm from "../components/auth/LoginForm";
import './LoginPage.css'

const LoginPage = (props) => {
  return (
    <div className={'login-page'}>
      <LoginForm history={props.history}
                 userStore={props.userStore}
                 loginStore={props.loginStore}
                 form={props.loginStore.form}
                 onChange={props.loginStore.onFieldChange}/>
    </div>
  )
};

export default LoginPage;