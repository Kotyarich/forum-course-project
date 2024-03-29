import {Link} from "react-router-dom";
import {observer} from "mobx-react";
import React, {Component} from "react";
import './Header.css'

const LoggedOutView = props => {
  if (!props.currentUser) {
    return (
      <div className="nav navbar-nav pull-xs-right">
        <div className="nav-item">
          <Link to="/login" className="nav-link button">
            Sign in
          </Link>
        </div>

        <div className="nav-item">
          <Link to="/register" className="nav-link button">
            Sign up
          </Link>
        </div>
      </div>
    );
  }
  return null;
};

const LoggedInView = props => {
  if (props.userStore.currentUser) {
    const onClick = (e) => {
      e.preventDefault();
      props.userStore.signOut();
    };

    return (
      <div className="nav navbar-nav">
        <div className="nav-item">
          <Link
            to={`/profile/${props.userStore.currentUser.nickname}`}
            className="nav-item navbar__username"
          >
            {props.userStore.currentUser.nickname}
          </Link>
          <button className={'nav-item button button_sign-out'}
                  onClick={onClick}>
            Sign Out
          </button>
        </div>

      </div>
    );
  }

  return null;
};

const StatisticView = props => {
  if (props.userStore.currentUser && props.userStore.currentUser.isAdmin) {
    return (
      <div className="nav navbar-nav">
        <div className="nav-item">
          <Link
            to={`/statistic`}
            className="nav-item navbar__username"
          >
            Statistic
          </Link>
        </div>
      </div>
    );
  }

  return null;
};

@observer
class Header extends Component {
  render() {
    return (
      <nav className="navbar navbar-light">
        <Link to="/" className="navbar-brand">
          Codemate finder
        </Link>

        <StatisticView userStore={this.props.userStore}/>
        <LoggedOutView currentUser={this.props.userStore.currentUser}/>
        <LoggedInView userStore={this.props.userStore}/>
      </nav>
    );
  }
}

export default Header;