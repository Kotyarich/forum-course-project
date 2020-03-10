import React from "react";
import {observer} from "mobx-react";
import "../components/base/Pagination.css"
import ReactPaginate from 'react-paginate'
import PostsList from "../components/post/PostsList";
import Thread from "../components/thread/Thread";
import Button from "../components/base/Button";

@observer
class ThreadPage extends React.Component {
  constructor(props) {
    super(props);
    this.POSTS_LIMIT = 10;
  }

  componentDidMount() {
    const slug = this.props.slug;
    const postStore = this.props.postStore;
    postStore.getThread(slug);
    postStore.getPosts(slug, this.POSTS_LIMIT, "flat");
  };

  onPageChange = (pageNumber) => {
    const slug = this.props.slug;
    const offset = pageNumber.selected * this.POSTS_LIMIT;
    const postStore = this.props.postStore;
    postStore.getPosts(slug, this.POSTS_LIMIT, "flat", false, offset);
  };

  render() {
    const thread = this.props.postStore.thread;
    const pageCount = thread.posts / this.POSTS_LIMIT;

    const posts = this.props.postStore.posts;
    console.log(thread);

    return (
      <div className={'threads-page'}>
        <Thread thread={thread}/>
        <PostsList posts={posts}/>
        <ReactPaginate pageCount={pageCount}
                       marginPagesDisplayed={1}
                       pageRangeDisplayed={4}
                       onPageChange={this.onPageChange}
                       containerClassName={'pagination'}
                       subContainerClassName={'pages pagination'}
                       breakClassName={'break-me'}
                       disabledClassName={'disabled'}
                       activeClassName={'active'}/>
        <div className={'thread__footer'}>
          <Button title={'Answer'}
                  name={'answer'}
                  action={() => {
                  }}/>
        </div>
      </div>
    );
  }
}

export default ThreadPage;