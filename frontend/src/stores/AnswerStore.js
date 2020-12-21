import {observable, runInAction} from 'mobx'
import GenericFormStore from './GenericFormStore'
import PostService from "../services/PostService";

class AnswerStore extends GenericFormStore {
  constructor() {
    super();
    this.postService = new PostService();
  }

  @observable
  form = {
    fields: {
      parent: {
        value: '0',
        error: null,
        rule: 'required'
      },
      message: {
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

  send = async (slug, author) => {
    try {
      const message = this.form.fields.message.value;
      console.log(message);
      const parent = this.form.fields.parent.value;
      return await this.postService.create(slug, author, message, parent);
    } catch (error) {
      runInAction(() => {
        this.form.meta.isValid = false;
        this.form.meta.error = error;
      });
      return error;
    }
  };
}

export default AnswerStore;