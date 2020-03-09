import React from "react";
import RegistrationForm from "../components/auth/RegistrationForm";
import './RegistrationPage.css'

const RegistrationPage = (props) => {
  return (
    <div className={'registration-page'}>
      <RegistrationForm history={props.history}
                        userStore={props.userStore}
                        registrationStore={props.registrationStore}
                        form={props.registrationStore.form}
                        onChange={props.registrationStore.onFieldChange}/>
    </div>
  )
};

export default RegistrationPage;