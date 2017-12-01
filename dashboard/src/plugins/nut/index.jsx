import React from 'react'
import {Route} from 'react-router'
import Exception from 'ant-design-pro/lib/Exception'

import Home from './Home'
const NotFound = () => (<Exception type="404"/>)
const routes = [
  (< Route key = "nut.home" exact path = "/" component = {
    Home
  } />),

  (<Route key="nut.no-match" component={NotFound}/>)
]

export default routes
