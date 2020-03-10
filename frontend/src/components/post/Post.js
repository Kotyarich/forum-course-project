import React from "react";
import {Link} from "react-router-dom";
import Button from "../base/Button";
import './Post.css'

const Post = (props) => {
  const {message, author, created, isEdited, id} = props.post;
  return (
    <div className={'post'}>
      <div className={'post__header'}>
        <Link className={'post__header__author'} to={'/profile/' + author}>
          {author}
        </Link>
        <div className={'post__header__id'}>
          {'#' + id}
        </div>
      </div>
      <div className={'post__main'}>
        <div className={'post__message'}>
          {message}
        </div>
        <hr/>
        <div className={'post__footer'}>
          <div className={'post__created'}>{created}</div>
          <Button title={'Answer'}
                  name={'answer'}
                  action={() => props.onAnswer(id)}/>
          {isEdited && <div className={'post__edited'}>{'edited'}</div>}
        </div>
      </div>
    </div>
  );
};

export default Post;