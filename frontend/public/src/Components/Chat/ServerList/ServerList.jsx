import React, { Component } from "react";
import "./ServerList.scss";

class ServerList extends Component {
  render() {
    console.log(this.props.chatHistory);

    // gets messages from app.js through props

    return (
      <div className="Servers">
        {/* <div className="ServerIcons" onClick={}> A </div> */}
        <div className="ServerIcons" > A </div>
        <div className="ServerIcons" > B </div>
        <div className="ServerIcons" > C </div>
        <div className="ServerIcons" > D </div>
        <div className="ServerIcons" > E </div>
      </div>
    );
  }
}

export default ServerList;