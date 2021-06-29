import React,{Component} from 'react'
import "./UserList.scss"
class  UserList extends Component {
  
 render() { 
   let users;
  if((this.props.userList).length > 0){
    // visible = true;
    users= <p className="usersLabel" > {"Users:" }</p> 
  } else {
    users= <p className="usersLabel" > {"" }</p> 

  }
  
    return (
      <div className="UserList">
        <div className="user-list">
          {users}
          {/* <p className="usersLabel" > {"Users:" && visible }</p> */}
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
