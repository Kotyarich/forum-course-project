import React from "react";
import {observer} from "mobx-react";
import "../components/base/Pagination.css"
import ReactPaginate from 'react-paginate'
import PostsList from "../components/post/PostsList";
import Thread from "../components/thread/Thread";
import PostForm from "../components/post/PostForm";
import "./ThreadPage.css"

@observer
class ThreadPage extends React.Component {
  constructor(props) {
    super(props);
    this.POSTS_LIMIT = 10;
  }

  componentDidMount() {
    this.currentPage = 0;
    const slug = this.props.slug;
    const postStore = this.props.postStore;
    postStore.getThread(slug);
    postStore.getPosts(slug, this.POSTS_LIMIT, "flat");
  };

  onPageChange = (pageNumber) => {
    this.currentPage = pageNumber.selected;
    const slug = this.props.slug;
    const offset = this.currentPage * this.POSTS_LIMIT;
    const postStore = this.props.postStore;
    postStore.getPosts(slug, this.POSTS_LIMIT, "flat", false, offset);
  };

  onAnswer = (id) => {
    this.props.answerStore.form.fields.parent.value = id;
    this.endFormRef.scrollIntoView({behavior: 'smooth'})
  };

  onSend = () => {
    const answerStore = this.props.answerStore;
    const slug = this.props.slug;
    const author = this.props.userStore.currentUser.nickname;
    answerStore.send(slug, author).then(() => {
      answerStore.form.fields.message.value = '';
      answerStore.form.fields.parent.value = '0';
      this.props.postStore.getThread(slug).then(() => {
        const thread = this.props.postStore.thread;
        this.currentPage = thread.posts / this.POSTS_LIMIT | 0;
        const offset = this.currentPage * this.POSTS_LIMIT;
        const postStore = this.props.postStore;
        postStore.getPosts(slug, this.POSTS_LIMIT, "flat", false, offset);
      });
    });
  };

  render() {
    const thread = this.props.postStore.thread;
    const pageCount = thread.posts / this.POSTS_LIMIT;

    const posts = this.props.postStore.posts;

    return (
      <>
        <div className={'thread-page__header'}/>
        <div className={'thread-page'}>
          <Thread thread={thread}/>
          <PostsList store={this.props.postStore}
                     user={this.props.userStore.currentUser}
                     posts={posts}
                     onAnswer={this.onAnswer}/>
          <ReactPaginate pageCount={pageCount}
                         marginPagesDisplayed={1}
                         pageRangeDisplayed={4}
                         forcePage={this.currentPage}
                         onPageChange={this.onPageChange}
                         containerClassName={'pagination'}
                         subContainerClassName={'pages pagination'}
                         breakClassName={'break-me'}
                         disabledClassName={'disabled'}
                         activeClassName={'active'}/>
          <hr/>
          <PostForm form={this.props.answerStore.form}
                    onSend={this.onSend}
                    onChange={this.props.answerStore.onFieldChange}/>
          <div ref={el => {
            this.endFormRef = el
          }}/>
        </div>
      </>
    );
  }
}

export default ThreadPage;