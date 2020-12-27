import React from 'react';
import ThreadItem from './ThreadItem'
import './ThreadsList.css'
import {observer} from "mobx-react";
import {Link} from "react-router-dom";

const ThreadsList = (props) => {
  return (
    <div className={"threads-list"}>
      <div className={'threads-list__header'}>
        <div className={'threads-list__header__rating threads__header'}>
          {'Rating'}
        </div>
        <div className={'threads-list__header__threads threads__header'}>
          {'Thread'}
        </div>
        <Link to={'/create-thread/'+props.slug} className={'thread-creation button'}>
          {'Create new thread'}
        </Link>
      </div>
      {props.threads.map((thread) =>
        <ThreadItem user={props.user}
                    key={thread.id}
                    thread={thread}
                    onClick={props.onClick}/>
      )}
      {props.threads.length === 0 && <div className={'threads-list_empty'}>
        {'This forum is empty'}
      </div>}
    </div>
  );
};

export default observer(ThreadsList);