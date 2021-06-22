// App.js
import React, { Component } from "react";
import {Route, Switch} from 'react-router-dom';
// import Header from "./Components/Head/Header";
// import ServerList from "./Components/ServerList";
import ChatPage from './Components/Chat/ChatPage';
import LoginPage from './Components/Log/Login';
import {ProtectedRoute} from './authorization/protected.route'
import "./App.css";

class App extends Component {
  render() {
    return (
      <div className="App">
        {/* <Header/> */}
        {/* <ServerList/> */}
        <Route> 
          <Switch>
            <Route exact path="/" component={LoginPage}/>
            <ProtectedRoute exact path="/chat" component={ChatPage} />
            <Route path="*" component={() => "404 NOT FOUND"}/>
          </Switch>
         </Route>
      </div>
    );
  }
}

export default App;