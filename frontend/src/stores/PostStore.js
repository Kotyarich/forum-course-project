import {observable, runInAction} from 'mobx'
import ThreadService from "../services/ThreadService";
import PostService from "../services/PostService";

class PostStore {
  constructor() {
    this.threadService = new ThreadService();
    this.postService = new PostService();
  }

  @observable
  status = 'initial';

  @observable
  posts = [];

  @observable
  thread = [];

  getThread = async (slug) => {
    const data = await this.threadService.getBySlug(slug);
    if (!data.author && data.message) {
      throw new Error(data.message);
    }

    runInAction(() => {
      this.thread = data;
    });
  };

  getPosts = async (slug, limit = 10, sort = "flat", desc = false,
                    offset = 0, since = 0) => {
    console.log(offset);
    const data = await this.postService
      .getByThreadSlug(slug, limit, sort, false, offset, since);
    if (data.message) {
      throw new Error(data.message);
    }

    runInAction(() => {
      this.posts = data;
    });
  };
}

export default PostStore;