import React, { Component } from "react";
import "./ServerList.scss";

class ServerList extends Component {
  render() {
    // gets messages from app.js through props
    // document.addEventListener("contextmenu", (event) => {
    //   event.preventDefault();
    //   const xPos = event.pageX + "px";
    //   const yPos = event.pageY + "px";
    //   //
    // });
    return (
      // <div className="ServerIcons" onClick={}> A </div> 
      <div className="Servers">
        <div className="server-list">
          { this.props.rooms && this.props.rooms.map((room, index) =>
          <button key={index} className="ServerIcons" value = {room.ID} 
           onClick={()=> this.props.changeRoom(room.ID) } > 
            {room.name.split(' ').map(function(item){return item[0]}).join('').toUpperCase()}
          </button>
           )} 
        </div>
      </div>
    );
  }
}

export default ServerList;