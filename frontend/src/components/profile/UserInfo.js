import React, {Component} from "react";
import {observer} from "mobx-react";
import Button from "../base/Button";
import ProfileLine from "./ProfileLine";
import './UserInfo.css'

@observer
class UserInfo extends Component {
  constructor(props) {
    super(props);
    this.nickname = this.props.nickname;
    this.state = {isChanging: false};
  }

  componentDidMount() {
    this.props.profileStore.getUserProfile(this.nickname);
  }

  handleChangeInfo(e) {
    e.preventDefault();
    this.setState({isChanging: true});
  }

  handleSubmit(e) {
    e.preventDefault();

    this.setState({isChanging: false});
    this.props.userStore.changeUserProfile({
      nickname: this.props.form.fields.nickname.value,
      fullname: this.props.form.fields.fullName.value,
      email: this.props.form.fields.email.value,
      about: this.props.form.fields.about.value,
    })
      .then(() => {
      });
  }

  render() {
    const {form, onChange} = this.props;
    const isChanging = this.state.isChanging;
    const fields = [
      form.fields.nickname,
      form.fields.fullName,
      form.fields.email,
      form.fields.about,
    ];
    let isCurrent = false;
    const user = this.props.userStore.currentUser;
    if (user) {
      if (user.nickname === form.fields.nickname.value || user.isAdmin) {
        isCurrent = true;
      }
    }

    return (
      <div className={'profile'}>
        <div className={'profile__title'}>Profile</div>
        <hr className={'profile__hr'}/>
        <div className={'profile__info'}>
          {fields.map((field, i) =>
            <ProfileLine key={'line_' + i}
                         title={field.title}
                         name={field.label}
                         onChange={onChange}
                         isChanging={isChanging}
                         value={field.value}
                         error={field.error}
                         type={field.type}/>
          )}
          <hr className={'profile__hr'}/>
          {isCurrent &&
          <div className={'profile__button-group'}>
            {!isChanging ? <Button title={'Change'}
                                   name={'profile-change'}
                                   action={(e) => this.handleChangeInfo(e)}/> :
              <Button title={'Submit'}
                      name={'submit'}
                      action={(e) => this.handleSubmit(e)}/>}
            {isChanging ? <Button title={'Cancel'}
                                  name={'profile-cancel'}
                                  action={(e) => this.handleSubmit(e)}/> : ''
            }
          </div>}
        </div>
      </div>
    );
  }
}

export default UserInfo;