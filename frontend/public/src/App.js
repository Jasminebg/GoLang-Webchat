// App.js
import React, { Component } from "react";
import Header from "./Components/Header";
import ChatHistory from "./Components/ChatHistory";
import "./App.css";
import { connect, sendMsg } from "./api";

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      chatHistory: [],
      typingStatus:""
    }
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
    this.typingStatus="User is typing";
    if(event.keyCode === 13) {
      sendMsg(event.target.value);
      event.target.value = "";
      this.typingStatus="";

    }
  }

  render() {
    return (
      <div className="App">
        <Header/>
        <ChatHistory chatHistory={this.state.chatHistory} />
        <p>{this.typingStatus}</p>
        <ChatInput send={this.send} />  
      </div>
    );
  }
}

export default App;