import {observable, runInAction} from 'mobx'
import GenericFormStore from './GenericFormStore'
import ThreadService from "../services/ThreadService";
import ForumService from "../services/ForumService";

class ThreadStore extends GenericFormStore  {
  constructor() {
    super();
    this.threadService = new ThreadService();
    this.forumService = new ForumService();
  }

  @observable
  form = {
    fields: {
      threadname: {
        value: '',
        error: null,
        rule: 'required'
      },
      initialpost: {
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

  createThread = async (forum_slug, author) => {
    try {
      await this.forumService.createThread({
         author: author,
         created: "2017-01-01T00:00:00.000Z",
         message: this.form.fields.initialpost.value,
         title: this.form.fields.threadname.value
      },
      forum_slug
      );
    } catch (error) {
      runInAction(() => {
        this.status = 'error';
      })
    }
  };

}

export default ThreadStore;