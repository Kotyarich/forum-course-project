import React from "react";
import {observer} from "mobx-react";
import ThreadsList from "../components/thread/ThreadsList";
import "../components/base/Pagination.css"
import ReactPaginate from 'react-paginate'

@observer
class ForumPage extends React.Component {
  constructor(props) {
    super(props);
    this.THREADS_LIMIT = 10;
  }

  componentDidMount() {
    const slug = this.props.slug;
    const threadStore = this.props.threadStore;
    threadStore.getForum(slug);
    threadStore.getThreads(slug, this.THREADS_LIMIT, 0);
  };

  onRatingChange = (id, value) => {
    const nickname = this.props.userStore.currentUser.nickname;
    this.props.threadStore.voteForThread(id, nickname, value);
  };

  onPageChange = (pageNumber) => {
    console.log(pageNumber);
    const slug = this.props.slug;
    const offset = pageNumber.selected * this.THREADS_LIMIT;
    this.props.threadStore.getThreads(slug, this.THREADS_LIMIT, offset);
  };

  render() {
    const forum = this.props.threadStore.forum;
    const pageCount = forum.threads / this.THREADS_LIMIT;

    const threads = this.props.threadStore.threads;

    return (
      <div className={'threads-page'}>
        <ThreadsList threads={threads} onClick={this.onRatingChange}/>
        <ReactPaginate pageCount={pageCount}
                       marginPagesDisplayed={1}
                       pageRangeDisplayed={4}
                       onPageChange={this.onPageChange}
                       containerClassName={'pagination'}
                       subContainerClassName={'pages pagination'}
                       breakClassName={'break-me'}
                       disabledClassName={'disabled'}
                       activeClassName={'active'}/>
      </div>
    );
  }
}

export default ForumPage;