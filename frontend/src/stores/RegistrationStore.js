import {observable, runInAction} from 'mobx'
import GenericFormStore from './GenericFormStore'
import UserService from "../services/UserService";

class RegistrationStore extends GenericFormStore {
  constructor() {
    super();
    this.userService = new UserService();
  }

  @observable
  form = {
    fields: {
      nickname: {
        value: '',
        error: null,
        rule: 'required'
      },
      fullName: {
        value: '',
        error: null,
        rule: 'required'
      },
      email: {
        value: '',
        error: null,
        rule: 'required|email'
      },
      password: {
        value: '',
        error: null,
        rule: 'required'
      },
      passwordSubmission: {
        value: '',
        error: null,
        rule: 'required|same:password'
      },
      about: {
        value: '',
        error: null,
        rule: []
      },
    },
    meta: {
      isValid: true,
      error: null,
    },
  };

  signUp = async () => {
    try {
      await this.userService.create({
        nickname: this.form.fields.nickname.value,
        fullname: this.form.fields.fullName.value,
        email: this.form.fields.email.value,
        about: this.form.fields.about.value,
        password: this.form.fields.password.value,
      });
    } catch (error) {
      runInAction(() => {
        console.log(error);
        this.form.meta.isValid = false;
        this.form.meta.error = error;
      })
    }
  };
}

export default RegistrationStore