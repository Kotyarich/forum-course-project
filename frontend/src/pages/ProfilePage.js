import React from "react";
import UserInfo from "../components/profile/UserInfo";
import './ProfilePage.css'
import {useParams} from "react-router";

const ProfilePage = (props) => {
  const {nickname} = useParams();
  return (
    <div className={'profile-page'}>
      <UserInfo nickname={nickname}
                userStore={props.userStore}
                profileStore={props.profileStore}
                form={props.profileStore.form}
                onChange={props.profileStore.onFieldChange}/>
    </div>
  )
};

export default ProfilePage;