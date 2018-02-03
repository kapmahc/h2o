import React, {Component} from 'react'
import {message} from 'antd'

import {get} from '../../../ajax'
import List from './List'

class Widget extends Component {
  state = {
    items: []
  }
  componentDidMount() {
    get('/forum/latest/articles').then((rst) => {
      this.setState({items: rst})
    }).catch(message.error)
  }
  render() {
    return (<List title={{
        id: "forum.articles.latest.title"
      }} items={this.state.items}/>)
  }
}

export default Widget
