import React from "react";
import {observer} from "mobx-react";
import ThreadCreatingForm from "../components/thread/ThreadCreatingForm";
import './ThreadCreationPage.css'


@observer
class ThreadCreationPage extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <div className={'thread-creation-page'}>
        <ThreadCreatingForm slug={this.props.slug}
                            history={this.props.history}
                            userStore={this.props.userStore}
                            threadStore={this.props.threadStore}
                            forumStore={this.props.forumStore}
                            form={this.props.threadStore.form}
                            onChange={this.props.threadStore.onFieldChange}/>
      </div>
    );
  }
}



export default ThreadCreationPage;