import React, { Component } from "react";
import Message from "../Message";
import "./ChatHistory.scss";

class ChatHistory extends Component {

  componentDidMount(){
    this.scrollToBottom();
    console.log(this.props.chatHistory)
  }

  componentDidUpdate(){
    this.scrollToBottom();
  }
  scrollToBottom(){
    this.el.scrollIntoView({behvaior:'smooth'});
  }

  render() {
    // console.log(this.props.chatHistory);
    const messages = this.props.chatHistory && this.props.chatHistory.map( (msg,index) => <Message key={index} message={msg} />);

    console.log("chat history");
    console.log(messages);
    return (
      <div className="ChatHistory">
        <div id="chatHistory" className="disable-scrollbars">
          <div id="history"> 
          {messages}
          </div>
          <div ref={el=>{this.el = el;}} />
        </div>
      </div>
    );
  }
}

export default ChatHistory;