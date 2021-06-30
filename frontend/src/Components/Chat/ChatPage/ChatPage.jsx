import React,{Component} from 'react';
import { Redirect } from 'react-router-dom';
import ChatSocket from '../../../api/ChatSocket';
import "./ChatPage.scss"
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
      rooms:{},
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
      this._chatSocket = new ChatSocket(auth.getUserName(), auth.getUserColour())
      this._chatSocket.connect((event) => {
        this.handleSocketEvent(event)
      });
    }
  }

  handleSocketEvent(event) {
    switch (event.type) {
      case "close":
        this.handleLogout()
        break;
      case "message":
        this.handleNewMessage(event)
        break;
      default:
    }
  }
  handleNewMessage(event){
    let data = event.data;
    data = data.split(/\r?\n/);

    for (let i = 0; i < data.length; i++){
      let msg = JSON.parse(data[i]);
      switch (msg.action){
        case "send-message":
          this.handleChatMessage(msg);
          break;
        case "user-join":
          this.handleUserJoined(msg);
          break;
        case "user-left":
          this.handleUserLeft(msg);
          break;
        case "room-joined":
          this.handleRoomJoined(msg)
          break;
        case "user-join-room":
          this.handleUserJoinedRoom(msg)
          break;
        case "list-clients":
          this.handleUserJoinedRoom(msg)
          break;
        default:
          break;

      }
      this.updateState();
    }
  };
  handleUserJoinedRoom(msg){
    let user = {
      name: msg.user,
      id: msg.id,
      color: msg.color
    };
    this.state.rooms[msg.roomid].users.push(user)

    this.sortUsersinRoom(msg.roomid);
    // if( this.state.rooms[msg.roomid].users.length > 0){
    //   this.state.rooms[msg.roomid].users.sort( (a,b) =>{
    //     if(a.name < b.name){return -1;}
    //     if(a.name > b.name){return 1;}
    //     {return 0;}
    //   } );
    // }
  };
  handleChatMessage(msg){
   if (typeof this.state.rooms[msg.roomid] !== "undefined"){
    let message = {
      msg:msg.message,
      user:msg.user,
      color:msg.color,
      timeStamp: msg.timestamp
    }
    this.state.rooms[msg.roomid].messages.push(message);
   } 
  };
  handleUserJoined(msg){
    let user = {
      name: msg.user,
      id: msg.id,
      color: msg.color
    };
    this.state.userList.push(user);
    
  };
  handleRoomJoined(msg){
    let user = {
      name: msg.user,
      id: msg.id,
      color: msg.color
    };
    if (typeof this.state.rooms[msg.roomid] === "undefined"){
      let room = {
        name:msg.room,
        ID: msg.roomid,
        messages:[],
        users:[user]
      };
      this.setState({
        ...this.state,
        rooms: {
            ...this.state.rooms,
            [msg.roomid]: room
        }
     });    
    } else {
      this.state.rooms[msg.roomid].users.push(user);

    }
    this.sortUsersinRoom(msg.roomid);
  };
  sortUsersinRoom(roomid){

    if( this.state.rooms[roomid].users.length > 0){
      this.state.rooms[roomid].users.sort( (a,b) =>{
        if(a.name < b.name){return -1;}
        if(a.name > b.name){return 1;}
        {return 0;}
      } );
    }
  }

  handleUserLeft(msg){
    for (let i =0; i< this.state.rooms[msg.roomid].users.length;i++){
      if (this.state.rooms[msg.roomid].users[i].id === msg.id){
        this.state.rooms[msg.roomid].users.splice(i,1);
      }
    }

  };

  componentWillUnmount(){
    this._chatSocket.closeSocket();
  };

  
  handleLogout(){
    auth.logout(()=> {
      this.props.history.push("/")
    })
  };

  send(event, room){
    if(event.keyCode === 13 && event.target.value !== "") {
      this._chatSocket.sendMessage(room, event.target.value);
      event.target.value = "";
    }
  };
  findRoom(event){
    if(event.keyCode === 13 && event.target.value !== ""){
      this._chatSocket.joinRoom(event.target.value);
      this._chatSocket.roomInput = event.target.value;
      event.target.value = "";
    }
  };
  updateState (){
    this.setState({
      rooms: this.state.rooms
    });
  };
  changeRoom= (roomID)=>{
    this.setState({
      room: this.state.rooms[roomID]
    });
  };
  privateMessage=(room)=>{
    this._chatSocket.joinPrivateRoom(room);
  };
    
  render() {
    if(!auth.isAuthenticated){
      return <Redirect to='/' />
    }
    
    return (
        <div className="ChatPage" >
          <Header roomName = {e=> this.findRoom(e)} currentRoom={this.state.room.name}/>
            {/* add leave room button in Serverlist? */}
            <div className="chatBody"> 
              <div className = "serverList"> 
                <ServerList rooms = {Object.values(this.state.rooms)} changeRoom={ this.changeRoom}/>
              </div>
              <div className = "roomPage"> 
                <ChatHistory chatHistory={this.state.room.messages}></ChatHistory>
                <ChatInput send={e=> this.send(e, this.state.room) } roomName={this.state.room.name}/>
              </div>
              <div className = "usersList"> 
                <UserList userList={this.state.room.users} privateMessage={this.privateMessage}/>
              </div>
            </div>
        </div>
    )
  }
}

export default ChatPage
