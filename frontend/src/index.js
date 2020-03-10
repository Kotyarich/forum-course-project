import React from 'react';
import ReactDOM from 'react-dom';
import {withRouter} from "react-router-dom";
import {
  BrowserRouter as Router,
  Switch,
  Route
} from "react-router-dom";
import RegistrationStore from "./stores/RegistrationStore";
import {observer} from "mobx-react";
import LoginStore from "./stores/LoginStore";
import UserStore from "./stores/UserStore";
import MainPage from "./pages/MainPage";
import RegistrationPage from "./pages/RegistrationPage";
import LoginPage from "./pages/LoginPage";
import ProfilePage from "./pages/ProfilePage";
import './index.css';
import Header from "./components/header/Header";
import ProfileStore from "./stores/ProfileStore";
import ForumStore from "./stores/ForumStore";
import ForumPage from "./pages/ForumPage";
import ThreadStore from "./stores/ThreadStore";
import ThreadPage from "./pages/ThreadPage";
import PostStore from "./stores/PostStore";

let registrationStore = new RegistrationStore();
let loginStore = new LoginStore();
let userStore = new UserStore();
let profileStore = new ProfileStore();
let forumStore = new ForumStore();
let threadStore = new ThreadStore();
let postStore = new PostStore();

@observer
class App extends React.Component {
  render() {
    return (
      <div>
        <Header userStore={userStore}/>
        <Switch>
          <Route exact path="/">
            <MainPage forumStore={forumStore}/>
          </Route>
          <Route path="/register">
            <RegistrationPage history={this.props.history}
                              userStore={userStore}
                              registrationStore={registrationStore}/>
          </Route>
          <Route path="/login">
            <LoginPage history={this.props.history}
                       userStore={userStore}
                       loginStore={loginStore}/>
          </Route>
          <Route path="/profile/:nickname">
            <ProfilePage userStore={userStore}
                         profileStore={profileStore}/>
          </Route>
          <Route path="/forum/:slug" render={(props) =>
            <ForumPage slug={props.match.params.slug}
                       userStore={userStore}
                       threadStore={threadStore}/>
          }/>
          <Route path="/thread/:slug" render={(props) =>
            <ThreadPage slug={props.match.params.slug}
                        userStore={userStore}
                        postStore={postStore}/>
          }/>
        </Switch>
      </div>
    );
  }
}

// ========================================
const AppWithRouter = withRouter(App);

ReactDOM.render(
  <Router>
    <AppWithRouter/>
  </Router>,
  document.getElementById('root')
);
