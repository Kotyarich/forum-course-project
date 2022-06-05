import React, {Component} from "react";
import {Link} from "react-router-dom";
import Button from "../base/Button";
import './Post.css'
import TextArea from "../base/TextArea";

class Post extends Component {
  constructor(props) {
    super(props);
    this.state = {isChanging: false, message: this.props.post.message};
  }

  handleChangeInfo(e) {
    e.preventDefault();
    this.setState({isChanging: true});
  }

  onChange = (name, value) => {
    this.setState({message: value})
  };

  save = (e) => {
    e.preventDefault();
    this.setState({isChanging: false});
    this.props.store.change(this.props.post.id, this.state.message);
  };

  render() {
    const currentUser = this.props.user;
    console.log(currentUser);
    const {author, parent, created, isEdited, id} = this.props.post;
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
            {this.state.isChanging ? <TextArea value={this.state.message}
                                               placeholder={''}
                                               onChange={this.onChange}
                                               name={''}/> : this.state.message}
            <hr/>
          </div>
          <div className={'post__footer'}>
            <div className={'post__created'}>{created}</div>
            {isEdited && <div className={'post__edited'}>{'edited'}</div>}
            {(currentUser.nickname === author || currentUser.isAdmin) && !this.state.isChanging &&
            <Button title={'Edit'}
                    name={'post__answer'}
                    action={(e) => this.handleChangeInfo(e)}/>}
            {this.state.isChanging ?
              <Button title={'Save'}
                      name={'post__answer'}
                      action={this.save}/> :
              <Button title={'Answer'}
                      name={'post__answer'}
                      action={() => this.props.onAnswer(id)}/>}
          </div>
          <hr/>
        </div>
      </div>
    );
  };
}

export default Post;