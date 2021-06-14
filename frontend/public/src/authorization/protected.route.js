import React from 'react'
import {Route, Redirect} from 'react-router-dom'
import auth from './auth'

export const ProtectedRoute = ({component: Component, ...rest}) => {


  return(
    <Route 
    {...rest}
    render = {(props) => auth.isAuthenticated()
    ? <Component {...props}/>
    : <Redirect to={{pathname:'/', state:{from:props.location}}} />}
    />
  )
}