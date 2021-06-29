class Auth {
  constructor(){
    this.sessionStorageUser = 'ChatUser';
  }

  login(name, colour, cb) {
    if(colour.charAt(0) === '#'){
      colour = colour.replace('#', '');
    }
    sessionStorage.setItem(this.sessionStorageUser, JSON.stringify({
      _name: name,
      _colour: colour
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

  getUserName() {
    return JSON.parse(sessionStorage.getItem(this.sessionStorageUser))._name;
  }

  getUserColour(){
    return JSON.parse(sessionStorage.getItem(this.sessionStorageUser))._colour;
  }  
}

export default new Auth()