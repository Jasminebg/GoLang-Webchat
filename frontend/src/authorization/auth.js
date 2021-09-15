class Auth {
  constructor(){
    this.sessionStorageUser = 'ChatUser';
    this.loginError = "";
  }

  async login(name, password, color, cb) {
    if(colour.charAt(0) === '#'){
      colour = colour.replace('#', '');
    }
    user = {
      username: this.name,
      password: this.password,
      color: this.color
    };

    try{
      const result = await axios.post("http://" + window.location.host + '/api/login', user)
      if (result.data.status !== "undefined" && result.data.status == "error"){
        this.loginError = "Login failed";
        cb();
      }else {
        token = result.data;
      }
    }catch(e){
      this.loginError = "Login failed";
      console.log(e);
      cb();
    }



    sessionStorage.setItem(this.sessionStorageUser, JSON.stringify({
      _name: name,
      _token: token,
      _password:password,
      _color: color
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