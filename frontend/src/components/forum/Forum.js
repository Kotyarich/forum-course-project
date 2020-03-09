import React from "react";
import './Forum.css'
import {Link} from "react-router-dom";

const Forum = (props) => {
  const {title, user, threads, posts, slug} = props.forum;
  return (
    <div className={"forum"}>
      <div className={'forum__header'}>
        <Link to={'/forum/' + slug} className={'forum__title forums__first-col'}>
          {title}
        </Link>
        <div className={'forum__user'}>
          Author: {user}
        </div>
      </div>
      <div className={'forum__numeric-info forums__col'}>
        {threads}
      </div>
      <div className={'forum__numeric-info forums__col'}>
        {posts}
      </div>
    </div>
  )
};

export default Forum;