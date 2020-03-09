import {observable, runInAction} from 'mobx'
import ForumService from "../services/ForumService";

class ForumStore {
  constructor() {
    this.forumService = new ForumService();
  }

  @observable
  status = 'initial';

  @observable
  forums = [];

  getForums = async () => {
    try {
      const data = await this.forumService.getAll();
      runInAction(() => {
        this.forums = data;
        this.status = "ok"
      });
    } catch (error) {
      runInAction(() => {
        this.status = 'error';
      })
    }
  };
}

export default ForumStore;