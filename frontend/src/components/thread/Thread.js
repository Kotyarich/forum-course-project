import React from 'react'
import {Link} from "react-router-dom";
import {observer} from "mobx-react";
import "./Thread.css"

@observer
class Thread extends React.Component {
  render() {
    const {author, message, created, title, isEdited} = this.props.thread;
    return (
      <div className={'thread'}>
        <div className={'thread__header'}>
          <div className={'thread__title'}>{title}</div>
          <Link className={'thread__author'} to={'/profile/' + author}>
            {'by ' + author}
          </Link>
        </div>
        <div className={'thread__message'}>
          <hr/>
          {message}
          <hr/>
        </div>
        <div className={'thread__created'}>{created}</div>
        {isEdited && <div className={'thread__edited'}>{'edited'}</div>}
        <hr/>
      </div>
    );
  }
}

export default Thread;