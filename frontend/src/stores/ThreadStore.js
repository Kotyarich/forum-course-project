import {observable, runInAction} from 'mobx'
import ThreadService from "../services/ThreadService";
import ForumService from "../services/ForumService";

class ThreadStore {
  constructor() {
    this.threadService = new ThreadService();
    this.forumService = new ForumService();
  }

  @observable
  status = 'initial';

  @observable
  threads = [];

  @observable
  forum = [];

  getForum = async (slug) => {
    const data = await this.forumService.getBySlug(slug);
    if (data.message) {
      throw new Error(data.message);
    }

    runInAction(() => {
      this.forum = data;
    });
  };

  getThreads = async (slug, limit = 10, offset = 0) => {
    const data = await this.threadService.getAll(slug, '', limit, offset);
    if (data.message) {
      this.status = 'error';
      throw new Error(data.message);
    }

    runInAction(() => {
      this.threads = data;
      this.status = "ok"
    });
  };

  voteForThread = async (threadSlug, nickname, vote) => {
    try {
      const thread = await this.threadService.vote(threadSlug, nickname, vote);
      const threadNumber = this.threads.findIndex((t) => t.id === thread.id);
      runInAction(() => {
        this.threads[threadNumber] = thread;
        this.status = "ok"
      });
    } catch (error) {
      runInAction(() => {
        this.status = 'error';
      })
    }
  };
}

export default ThreadStore;