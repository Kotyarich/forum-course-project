import React from 'react';
import Post from './Post'
import './PostsList.css'
import {observer} from "mobx-react";

const PostsList = (props) => {
  console.log(props.posts);

  return (
    <div className={"posts-list"}>
      {props.posts.map((post) =>
        <Post key={post.id} post={post}/>
      )}
    </div>
  );
};

export default observer(PostsList);