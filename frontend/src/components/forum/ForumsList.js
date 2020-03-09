import React from 'react';
import Forum from './Forum'
import './ForumsList.css'

const ForumsList = (props) => {
  return (
    <div className={"forums-list"}>
      <div className={'forums-list__header'}>
        <div className={'forums-list__header__title forums__first-col forums__header'}>
          {'Forums'}
        </div>
        <div className={'forums-list__header__threads forums__header forums__col'}>
          {'Threads'}
        </div>
        <div className={'forums-list__header__posts forums__header forums__col'}>
          {'Posts'}
        </div>
      </div>
      {props.forums.map((forum) =>
        <Forum key={forum.slug} forum={forum}/>
      )}
    </div>
  )
};

export default ForumsList;