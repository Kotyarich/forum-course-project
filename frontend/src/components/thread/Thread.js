import React from 'react'
import {Link} from "react-router-dom";
import Button from "../base/Button";
import {observer} from "mobx-react";

@observer
class Thread extends React.Component{
  render() {
    const {author, message, created, isEdited, id} = this.props.thread;
    return(
      <div className={'thread'}>
        <div className={'thread__author'}>
          <Link className={'thread__author__link'} to={'/profile' + author}>
            {author}
          </Link>
        </div>
        <div className={'thread__main'}>
          <div className={'thread__created'}>{created}</div>
          <div className={'thread__message'}>
            {message}
          </div>
          <div className={'thread__footer'}>
            <Button title={'Answer'}
                    name={'answer'}
                    action={() => props.onAnswer(id)}/>
            {isEdited && <div className={'thread__edited'}>{'edited'}</div>}
          </div>
        </div>
      </div>
    );
  }
}

export default Thread;