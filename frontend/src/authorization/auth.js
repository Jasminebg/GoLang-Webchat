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
      username: name,
      password: password,
      color: color,
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
      _name: user.name,
      _token: user.token,
      _password:user.password,
      _color: user.color
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