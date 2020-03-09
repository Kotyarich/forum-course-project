import React from "react";
import './ThreadItem.css'
import {Link} from "react-router-dom";
import Rating from "../base/Rating";

const ThreadItem = (props) => {
  const {title, author, created, votes, id} = props.thread;
  return (
    <div className={'thread-item'}>
      <Rating id={id}
              votes={votes}
              onClick={props.onClick}/>
      <div className={'thread-item__header'}>
        <Link to={'/thread/' + id} className={'thread-item__title'}>{title}</Link>
        <div className={'thread-item__started'}>
          <span>{'Started by '}</span>
          <Link className={'thread-item__author'} to={'/profile/' + author}>
            {author}
          </Link>
          <span>{', ' + created}</span>
        </div>
      </div>
    </div>
  )
};

export default ThreadItem;