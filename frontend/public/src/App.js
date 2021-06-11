// App.js
import React, { Component } from "react";
import Header from "./Components/Header";
import ChatHistory from "./Components/ChatHistory";
import ServerList from "./Components/ServerList";
import "./App.css";
import { connect, sendMsg } from "./api";
import ChatInput from "./Components/ChatInput";

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      chatHistory: [],
      typingStatus:""
    }
    // this.send = this.send.bind(this);
  }
  componentDidMount() {
    connect((msg) => {
      console.log("New Message")
      this.setState(prevState => ({
        chatHistory: [...this.state.chatHistory, msg]
      }))
      console.log(this.state);
    });
  }
  send(event) {
    this.setState ((prevState) => ({
      typingStatus:"User is typing..."
    }))
    // this.state.typingStatus="User is typing...";
    if(event.keyCode === 13) {
      sendMsg(event.target.value);
      event.target.value = "";
      // typingStatus="";

    }
  }

  render() {
    return (
      <div className="App">
        <Header/>
        <div style={{display:'grid', gridTemplate:' 1fr /  1fr 15fr ',  alignContent:'center', justifyContent:'center'}}>
        <ServerList/>
        <ChatHistory chatHistory={this.state.chatHistory} />
        <p>{this.typingStatus}</p>
        <ChatInput send={this.send} />  
        </div>
      </div>
    );
  }
}

export default App;