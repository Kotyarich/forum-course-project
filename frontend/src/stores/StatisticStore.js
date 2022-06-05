import {observable, runInAction} from 'mobx'
import GenericFormStore from './GenericFormStore'
import StatisticService from "../services/StatisticService";

class ThreadStore extends GenericFormStore  {
  constructor() {
    super();
    this.statisticService = new StatisticService();
  }

  @observable
  statistics = {
    users: 0,
    posts: 0,
    votes: 0,
    threads: 0,
    forums: 0,
  };

  get = async () => {
    const data = await this.statisticService.get();
    if (data.message) {
      throw new Error(data.message);
    }

    runInAction(() => {
      this.statistics = data;
    });
  };
}

export default ThreadStore;