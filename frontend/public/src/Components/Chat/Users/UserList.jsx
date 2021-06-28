import React,{Component} from 'react'
import "./UserList.scss"

class  UserList extends Component {
  
 render() { 
    return (
      <div className="UserList">
        <div className="user-list">
          <p className="usersLabel"> Users:</p>
           {this.props.userList && this.props.userList.map((user, index) => 
          <p key={index} className="user" style={{ color: "#"+user.color }} onClick={()=> this.props.privateMessage(user.id)}>
            {user.name}
            </p>
        )}
        </div>
      </div>
    );
  }
}

export default UserList
