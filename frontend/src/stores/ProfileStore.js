import {observable, runInAction} from 'mobx'
import GenericFormStore from './GenericFormStore'
import UserService from "../services/UserService";

class ProfileStore extends GenericFormStore {
  constructor() {
    super();
    this.userService = new UserService();
  }

  @observable
  status = 'initial';

  @observable
  form = {
    fields: {
      nickname: {
        value: '',
        error: null,
        label: 'nickname',
        title: 'Nickname',
        type: 'input',
        rule: 'required'
      },
      fullName: {
        value: '',
        error: null,
        label: 'fullName',
        title: 'Full Name',
        type: 'input',
        rule: 'required'
      },
      email: {
        value: '',
        error: null,
        label: 'email',
        title: 'Email',
        type: 'input',
        rule: 'required|email'
      },
      about: {
        value: '',
        error: null,
        label: 'about',
        title: 'About',
        type: 'textarea',
        rule: []
      },
    },
    meta: {
      isValid: true,
      error: null,
    },
  };

  getUserProfile = async (nickname) => {
    try {
      const data = await this.userService.get(nickname);
      runInAction(() => {
        this.form.fields.nickname.value = data.nickname;
        this.form.fields.fullName.value = data.fullname;
        this.form.fields.email.value = data.email;
        this.form.fields.about.value = data.about;
      });
    } catch (error) {
      runInAction(() => {
        this.status = 'error';
      })
    }
  };
}

export default ProfileStore;