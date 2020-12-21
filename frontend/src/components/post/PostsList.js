import React from 'react';
import Post from './Post'
import './PostsList.css'
import {observer} from "mobx-react";

const PostsList = (props) => {
  return (
    <div className={"posts-list"}>
      {props.posts.map((post) =>
        <Post key={post.id} post={post} onAnswer={() => {
          props.onAnswer(post.id)
        }}/>
      )}
    </div>
  );
};

export default observer(PostsList);