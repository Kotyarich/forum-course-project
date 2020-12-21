import React from "react";
import {observer} from "mobx-react";
import ForumsList from "../components/forum/ForumsList";

@observer
class MainPage extends React.Component {
  componentDidMount() {
    this.props.forumStore.getForums();
  }

  render() {
    const forums = this.props.forumStore.forums;
    return(
      <div className={'main-page'}>
        <ForumsList forums={forums}/>
      </div>
    )
  }
}

export default MainPage;