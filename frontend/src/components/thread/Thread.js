import React from 'react'
import {Link} from "react-router-dom";
import {observer} from "mobx-react";
import "./Thread.css"

@observer
class Thread extends React.Component{
  render() {
    const {author, message, created, title, isEdited} = this.props.thread;
    return(
      <div className={'thread'}>
        <div className={'thread__author'}>
          <Link className={'thread__author__link'} to={'/profile/' + author}>
            {author}
          </Link>
        </div>
        <div className={'thread__main'}>
          <div className={'thread__title'}>{title}</div>
          <div className={'thread__message'}>
            {message}
          </div>
          <div className={'thread__created'}>{created}</div>
        </div>
        {isEdited && <div className={'thread__edited'}>{'edited'}</div>}
      </div>
    );
  }
}

export default Thread;