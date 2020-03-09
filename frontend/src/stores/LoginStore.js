import {observable, runInAction} from 'mobx'
import GenericFormStore from './GenericFormStore'
import UserService from "../services/UserService";

class LoginStore extends GenericFormStore {
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
      password: {
        value: '',
        error: null,
        rule: 'required'
      },
    },
    meta: {
      isValid: true,
      error: null,
    },
  };

  signIn = async () => {
    try {
      return await this.userService.signIn({
        nickname: this.form.fields.nickname.value,
        password: this.form.fields.password.value,
      });
    } catch (error) {
      runInAction(() => {
        this.form.meta.isValid = false;
        this.form.meta.error = error;
      })
    }
  };
}

export default LoginStore