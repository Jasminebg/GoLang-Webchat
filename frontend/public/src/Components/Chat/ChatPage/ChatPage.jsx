import React,{Component} from 'react';
import { Redirect } from 'react-router-dom';
import ChatSocket from '../../../api/ChatSocket';
import "./ChatPage"
import ChatHistory from "../ChatHistory";
import ChatInput from "../ChatInput";
import UserList from "../Users";
import auth from '../../../authorization/auth';
import Header from '../../Head/Header'


class ChatPage extends Component {
  _chatSocket;
 
  constructor(props){
    super(props);
    this.state={
      isActive:false,
      chatHistory:[],
      userList:[],
      rooms:[]
    }
  }
  
  handleChange = (e) => {
    this.setState({
      [e.target.name]: e.target.value
    }
    );
  }

  handleShow = () => {
    this.setState({
      isActive:true
    });
  };

  handleHide = () => {
    this.setState({
      isActive:false
    });
  };

  componentDidMount(){
    if ( auth.isAuthenticated()){
      // console.log(auth.getUserColour());
      this._chatSocket = new ChatSocket("ws://localhost:8080/ws", auth.getUserName(), auth.getUserColour(), true)

    }
  }
  componentWillUnmount(){
    this._chatSocket.closeSocket();
  }

  
  handleLogout(){
    auth.logout(()=> {
      this.props.history.push("/")
    })
  }

  send(event){
    if(event.keyCode === 13 && event.target.value !== "") {
      // console.log(event.target.value)
      this._chatSocket.sendMessage(this.state.roomName, event.target.value);
      event.target.value = "";
      // console.log(this._chatSocket.room.messages)
      console.log("send event")
      console.log(this._chatSocket.room)
    }
  }
  findRoom(event){
    // console.log(event.target.value);
    if(event.keyCode === 13 && event.target.value !== ""){
      this._chatSocket.joinRoom(event.target.value);
      this._chatSocket.roomInput = event.target.value;
      this.setState({
        roomName : event.target.value
      });
      console.log("join  event")
      console.log(this._chatSocket.rooms)
  
      event.target.value = "";
      // console.log(this._chatSocket.room)
    }
  }
  

  
  render() {
    if(!auth.isAuthenticated){
      return <Redirect to='/' />
    }
    return (
      <div className="ChatPage">
        <Header roomName = {e=> this.findRoom(e)}/>
        <UserList userList={this.state.userList}/>
        <ChatHistory chatHistory={this.state.chatHistory} ></ChatHistory>
        {/* <ChatHistory chatHistory = {this.state.chatHistory}/> */}
        <ChatInput send={e=> this.send(e)}/>
      </div>
    )
  }
}

export default ChatPage
