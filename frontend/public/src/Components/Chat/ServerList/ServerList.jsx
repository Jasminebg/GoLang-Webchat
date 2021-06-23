import React, { Component } from "react";
import "./ServerList.scss";

class ServerList extends Component {
  render() {
    console.log("Servers");
    console.log(this.props.rooms);

    // gets messages from app.js through props

    return (
      // <div className="ServerIcons" onClick={}> A </div> 
      <div className="Servers">
        <div className="server-list">
          { this.props.rooms && this.props.rooms.map((room, index) =>
          <div key={index} className="ServerIcons" value = {room.ID} onClick={this.props.changeRoom} > 
            {room.name.split(' ').map(function(item){return item[0]}).join('')}
          </div>
           )} 
        </div>
      </div>
    );
  }
}

export default ServerList;