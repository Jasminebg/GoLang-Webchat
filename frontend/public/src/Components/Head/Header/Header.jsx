import React, { Component } from "react";
import auth from "../../../authorization/auth";
import { withRouter } from 'react-router-dom';
import "./Header.scss";

class Header extends Component{

  constructor(props){
    super(props);
    this.state={
      roomName:""
    };
  }
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
      <input  onKeyDown={this.props.roomName}/>
      {/* <input type="input" className="form__field" name='roomName' id="roomName"  value={this.state.roomName} placeholder="room name..." onChange={ this.props.send}/> */}
      <h2>Chat App</h2>
      <button className="logout-button" onClick={() => { this.handleLogout()}}>
        Logout
      </button>
    </div>
    )
  }
}

export default withRouter(Header);