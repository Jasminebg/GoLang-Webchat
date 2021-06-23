import React,{Component} from 'react';
import { Redirect } from 'react-router-dom';
import ChatSocket from '../../../api/ChatSocket';
import "./ChatPage"
import ChatHistory from "../ChatHistory";
import ChatInput from "../ChatInput";
import UserList from "../Users";
import auth from '../../../authorization/auth';
import Header from '../../Head/Header'
import ServerList from '../ServerList'


class ChatPage extends Component {
  _chatSocket;
 
  constructor(props){
    super(props);
    this.state={
      isActive:false,
      chatHistory:[],
      userList:[],
      rooms:[],
      room : {
        name:'',
        ID: '',
        messages:[],
        users:[]
      }
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
      this._chatSocket = new ChatSocket(auth.getUserName(), auth.getUserColour())

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

  send(event, room){
    if(event.keyCode === 13 && event.target.value !== "") {
      // console.log(event.target.value)
      this._chatSocket.sendMessage(room, event.target.value);
      event.target.value = "";
      // console.log(this._chatSocket.room.messages)
      this.setState ({
        rooms : Object.values(this._chatSocket.rooms)
      });
      // this.state.rooms = Object.values(this._chatSocket.rooms);
      // console.log("send event")
      // console.log(this._chatSocket.room)
    }
  }
  findRoom(event){
    // console.log(event.target.value);
    if(event.keyCode === 13 && event.target.value !== ""){
      this._chatSocket.joinRoom(event.target.value);
      this._chatSocket.roomInput = event.target.value;
      console.log("join  event")
      this.setState ({
        rooms : Object.values(this._chatSocket.rooms)
      });
      // this.state.rooms = Object.values(this._chatSocket.rooms);
      // this.room = this.rooms[event.target.value];
      event.target.value = "";
    }
  }

  changeRoom(event){
    this.setState({
      room: this._chatSocket.findRoom(event.target.value)
    });

  }
  

  
  render() {
    if(!auth.isAuthenticated){
      return <Redirect to='/' />
    }
    
          {/* { this._chatSocket.rooms && this._chatSocket.rooms.map((room, index) => */}
    return (
      // <div style = {{height:'100vw', backgroundColor:'#36393F'}}> 
        <div className="ChatPage" >
          <Header roomName = {e=> this.findRoom(e)}/>
          <ServerList rooms = {this.state.rooms} changeRoom={e=> this.changeRoom(e)}/>
          {/* { this.state.rooms && this.state.rooms.map((room, index) =>
          <ServerList key = {index} roomName={room.name} roomID = {room.ID} 
            changeRoom= {this.changeRoom(room.ID) } />

              )} */}
            <div className = "roomPage"> 
              <UserList userList={this.state.room.users}/>
              <ChatHistory chatHistory={this.state.room.messages}></ChatHistory>
              <ChatInput send={e=> this.send(e, this.state.room)}/>
            </div>
        </div>
      // </div>
      //   <div className="ChatPage" >
      //   <Header roomName = {e=> this.findRoom(e)}/>
      //   { this.state.rooms && this.state.rooms.map((room, index) =>

      //     <div key={index} className = "roomPage"> 
      //       <UserList userList={room.users}/>
      //       <ChatHistory chatHistory={room.messages}></ChatHistory>
      //       <ChatInput send={e=> this.send(e, room)}/>
      //     </div>
      //   )}
      // </div>
    )
  }
}

export default ChatPage
