class ChatSocket{
  _socketEndpoint;
  _socket;
  users = [];
  user= {
    username: this.userName,
    userID: ''
  };
  // roomInput;
  rooms = [];
  room = {
    name:'',
    ID: '',
    messages:[]
  };
  cs;

  constructor(socketEndpoint, userName, userColour, connect=false){
    // this._socketEndpoint = `${socketEndpoint}?user=${userName}&userColour=${userColour}&userId=${userId}`;
    // this._socket = connect ? new WebSocket(this._socketEndpoint):null;
    this.cs = new WebSocket(`ws://localhost:8080/ws?user=${userName}&userColour=${userColour}`)
    this.cs.addEventListener('open', (event) => {this.createSocket(event)});
    this.cs.addEventListener('message', (event) => {this.handleNewMessage(event)});

  }

  handleNewMessage(event){
    // console.log(event);
    let data = event.data;
    console.log("event");
    console.log(event);
    data = data.split(/\r?\n/);
    // console.log("data");
    // console.log(JSON.parse(data));
    // console.log(".")

    for (let i = 0; i < data.length; i++){
      // console.log(data[i]);
      let msg = JSON.parse(data[i]);
      this.user.userID = msg.id;
      console.log("handlenewmessage");
      console.log(msg);
      // console.log(msg.message);
      // console.log("switch");
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
        default:
          break;

      }
      // console.log(this.room.messages);
    }
  }
  handleChatMessage(msg){
   this.room = this.findRoom(msg.roomID);
   if (typeof this.room !== "undefined"){
     this.room.messages.push(msg);
   } 
  };
  handleUserJoined(msg){
    this.users.push(msg.user);
  };
  handleRoomJoined(msg){
    this.room.name = msg.room;
    this.room.ID = msg.roomid;
    this.rooms[this.room.name] = this.room;
  };

  handleUserLeft(msg){
    for (let i =0; i< this.users.length;i++){
      if (this.users[i].id === msg.user){
        this.users.splice(i,1);
      }
    }
  };

  sendMessage(roomname, msg){
    // this.room = this.findRoom(roomname);
    // console.log(this.room); 
    // console.log("send") 
    if (msg !== ""){
      this.cs.send(JSON.stringify({
        action:"send-message",
        message:msg,
        roomid:this.room.ID,
        room:this.room.name
      }));
      
      // this.room.message.push(msg);
      // this.rooms[roomname].messages.push(msg);
    }
    // console.log("send message");
    // console.log(this.room);
    // console.log(this.rooms);
  }

  createSocket(){
    // if (this.cs == null){
    //   this.cs = new WebSocket(this.cs)
    // }
    console.log("WS open!");
  }

  closeSocket(){
    if (this.cs != null){
      this.cs.close();
      this.cs=null;
    }
  }

  findRoom(roomID){
    for (let i = 0; i < this.rooms.length; i++){
      if(this.rooms[i].id === roomID){
        return this.rooms[i];
      }
    }
  };
  joinRoom(roomName){
    this.cs.send(JSON.stringify({action:'join-room', message:roomName}))
    // this.room.name = roomName;
    // this.roomInput="";
  };
  leaveRoom(room){
    this.cs.send(JSON.stringify({action:'leave-room', message:room}))
    for (let i=0;i<this.rooms.length;i++){
      if(this.rooms[i].ID === room){
        this.rooms.splice(i,1);
        break;
      }
    }
  };
  joinPrivateRoom(room){
    this.cs.send(JSON.stringify({action:'join-room-private', message:room}))
  };
  // getRoomMessages(){
  //   return this.room.messages;
  // }
}
export default ChatSocket;