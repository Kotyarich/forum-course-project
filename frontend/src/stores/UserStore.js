import {observable, runInAction} from 'mobx'
import UserService from "../services/UserService";

class UserStore {
  constructor() {
    this.userService = new UserService();
    this.userService.checkAuth().then(user => {
      if (!user.error) {
        this.currentUser = user;
      }
    })
  }

  @observable
  status = 'initial';

  @observable
  currentUser = null;

  signOut = () => {
    this.currentUser = null;
    this.userService.singOut();
  };

  getUserProfile = async () => {
    try {
      const data = await this.userService.get(this.nickname);
      runInAction(() => {
        this.currentUser = {
          nickname: data.nickname,
          fullName: data.fullname,
          email: data.email,
          about: data.about,
        };
        this.status = 'ok'
      });
    } catch (error) {
      runInAction(() => {
        this.status = 'error';
      })
    }
  };
}

export default UserStore;