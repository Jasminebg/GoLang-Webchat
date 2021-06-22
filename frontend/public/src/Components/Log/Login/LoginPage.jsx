import React, {Component} from 'react'
import auth from '../../../authorization/auth';
import "./LoginPage.scss"

class LoginPage extends Component {

  constructor(props){
    super(props);
    this.state={
      name:"",
      colour:"", 
      continue: true
    };
  }
  handleChange = (e) => {
    this.setState({
      [e.target.name]: e.target.value
    });
    // console.log(/^#?(([0-9a-fA-F]{2}){3}|([0-9a-fA-F]){3})$/i.test(this.state.color));
    // console.log(this.state.colour);
  }


  render () {
    return (
      <div className="LoginPage">
        <div className="loginContainer" onKeyPress={this.keyPressed}>
          <div className="form__group field">
            <input type="input" className="form__field" name='name' id="name" value={this.state.name} onChange={e=> this.handleChange(e)}/>
            <label htmlFor="name" className="form__label">Username</label>
          </div>
          <div className="form__group field">
            <input type="input" className="form__field" name='colour' id="colour" value={this.state.color} onChange={e=> this.handleChange(e)}/>
            <label htmlFor="colour" className="form__label">Colour  <span style={{color:'#E92750'}}>(eg default: Red[#E92750]) </span>  </label>
          </div>
          <button className="login-button"  onKeyPress={this.onKeyPress} onClick={this.verifyInput}>Login</button>

        </div>

      </div>

    )
    }

    submitLogin =()=>{
      auth.login(this.state.name, this.state.colour,()=>{
        this.props.history.push("/chat")
      })
    }
    
    verifyInput=()=>{
      if (this.state.name !== ''  ){
        this.submitLogin();
      }
    }

    // checkColor=()=>{
    //   if(/^#?([0-9A-F]{6})$/i.test(this.state.colour)){
    //     return true;
    //   }
    //   this.setState({
    //     colour:"E92750"
    //   }, () => console.log(this.state.colour) );
    //   return true;
      
    // }

    keyPressed =(event)=>{
      if (event.key === "Enter"){
        this.verifyInput();
      }
    }
}

export default LoginPage
