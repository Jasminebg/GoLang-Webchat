import React, { Component } from "react";
import "./ChatInput.scss";


class ChatInput extends Component {
  render() {
    let chatInput;
   if((this.props.roomName).length > 0){
     // visible = true;
     chatInput= <input onKeyDown={this.props.send} placeholder={"Message " +  this.props.roomName }/>
   } else {
    chatInput= <div/>
 
   }
    return (
      <div className="ChatInput" >
            {chatInput}
      </div>
    );
  }
}

export default ChatInput;