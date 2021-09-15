import axios from 'axios';

class Auth {
  constructor(){
    this.sessionStorageUser = 'ChatUser';
    this.loginError = "";
  }

  async login(name, password, color, cb) {
    if(color.charAt(0) === '#'){
      color = color.replace('#', '');
    }
    let user = {
      username: this.name,
      password: this.password,
      color: this.color,
      token: ""
    };

    try{
      const result = await axios.post("http://" + window.location.host + '/api/login', user)
      if (result.data.status !== "undefined" && result.data.status == "error"){
        this.loginError = "Login failed";
        cb();
      }else {
        this.user.token = result.data;
      }
    }catch(e){
      this.loginError = "Login failed";
      console.log(e);
      cb();
    }



    sessionStorage.setItem(this.sessionStorageUser, JSON.stringify({
      _name: this.user.name,
      _token: this.user.token,
      _password:this.user.password,
      _color: this.user.color
    }));
    cb();
  }

  logout(cb) {
    sessionStorage.removeItem(this.sessionStorageUser)
    cb();
  }

  isAuthenticated() {
    var test = sessionStorage.getItem(this.sessionStorageUser);
    return test;
  }

  getUserToken(){
    return JSON.parse(sessionStorage.getItem(this.sessionStorageUser))._token;
  }

  getUserName() {
    return JSON.parse(sessionStorage.getItem(this.sessionStorageUser))._name;
  }

  getUserColor(){
    return JSON.parse(sessionStorage.getItem(this.sessionStorageUser))._color;
  }  
}

export default new Auth()