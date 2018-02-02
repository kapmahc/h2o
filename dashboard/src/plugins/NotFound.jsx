import React, {Component} from 'react'
import Exception from 'ant-design-pro/lib/Exception'

import Layout from '../layout'

class Widget extends Component {
  render() {
    return (<Layout breads={[]}><Exception type="404"/></Layout>)
  }
}

export default Widget
