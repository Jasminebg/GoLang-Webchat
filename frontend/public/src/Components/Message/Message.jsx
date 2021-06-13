import React, { Component } from "react";
import "./Message.scss";

class Message extends Component {
  constructor(props) {
    super(props);
    let temp = JSON.parse(this.props.message);
    this.state = {
      message: temp,
      timeStamp: this.displayTime(temp.timeStamp)
    };
  }

  render() {
    return <div className="Message">
      <span className ="timeStamp">{this.state.timeStamp}</span>
      <span className="userName" style={{color:this.state.message.color}}>
        {this.state.message.user}&nbsp;
      </span>
      <span classname="messageBody">{this.state.message.body}</span>
      </div>;
  }
}

export default Message;