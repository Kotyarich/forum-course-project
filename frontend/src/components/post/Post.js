import React from "react";
import {Link} from "react-router-dom";
import Button from "../base/Button";
import './Post.css'

const Post = (props) => {
  const {message, author, parent, created, isEdited, id} = props.post;
  return (
    <div className={'post'}>
      <div className={'post__header'}>
        <Link className={'post__header__author'} to={'/profile/' + author}>
          {author}
        </Link>
        {parent !== 0 && <div className={'post__reply-to'}>
          {'reply to post #' + parent}
        </div>}
        <div className={'post__header__id'}>
          {'#' + id}
        </div>
      </div>
      <div className={'post__main'}>
        <div className={'post__message'}>
          <hr/>
          {message}
          <hr/>
        </div>
        <div className={'post__footer'}>
          <div className={'post__created'}>{created}</div>
          {isEdited && <div className={'post__edited'}>{'edited'}</div>}
          <Button title={'Change'}
                  name={'post__change'}
                  action={() => {}}/>
          <Button title={'Answer'}
                  name={'post__answer'}
                  action={() => props.onAnswer(id)}/>
        </div>
        <hr/>
      </div>
    </div>
  );
};

export default Post;