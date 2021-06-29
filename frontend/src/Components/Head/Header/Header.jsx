import React, { Component } from "react";
import auth from "../../../authorization/auth";
import { withRouter } from 'react-router-dom';
import "./Header.scss";

class Header extends Component{

  // handleChange = (e) => {
  //   this.setState({
  //     [e.target.name]: e.target.value
  //   }
  //   , () =>  console.log(this.state.roomName)
  //   );
  // }

  handleLogout(){
    auth.logout(() => {
      return this.props.history.push('/');
    })
  }
  render() {
    return(
    <div className="header">
      <input placeholder="Join room" onKeyDown={this.props.roomName}/>
      <h2>{ this.props.currentRoom  || "Chat App" }</h2>
      <button className="logout-button" onClick={() => { this.handleLogout()}}>
        Logout
      </button>
    </div>
    )
  }
}

export default withRouter(Header);