class ChatSocket{
  cs;

  constructor( userName, userColour){
    this.cs = new WebSocket(`ws://localhost:8080/ws?user=${userName}&userColour=${userColour}`)

  }
  connect(cb) {
    console.log("Attempting Connection...");
  
    this.cs.onopen = () => {
      console.log("Successfully Connected");
    };
  
    this.cs.onmessage = (msg) => {
      cb(msg)
    };
  
    this.cs.onclose = (event) => {
      console.log("Socket Closed Connection: ", event);
      cb(event)
    };
  
    this.cs.onerror = (error) => {
      console.log("Socket Error: ", error);
    };
  };


  sendMessage(room, msg){
    if (msg !== ""){
      this.cs.send(JSON.stringify({
        action:"send-message",
        message:msg,
        roomid:room.ID,
        room:room.name
      }));      
    }
  }

  createSocket(){
    console.log("WS open!");
  }

  closeSocket(){
    if (this.cs != null){
      this.cs.close();
      this.cs=null;
    }
  }

  joinRoom(roomName){
    this.cs.send(JSON.stringify({action:'join-room', message:roomName}))
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
}
export default ChatSocket;