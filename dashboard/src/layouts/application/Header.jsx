import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {Layout, Menu} from 'antd';
import {connect} from 'react-redux'

const {Header} = Layout;

class Widget extends Component {
  render() {
    console.log(document.body.offsetWidth);
    // TODO 根据宽度计算headernav
    const {site} = this.props
    return (<Header>
      <div className="logo"/>
      <Menu theme="dark" mode="horizontal" defaultSelectedKeys={[]} style={{
          lineHeight: '64px'
        }}>
        <Menu.Item key="home">{site.subhead}</Menu.Item>
        <Menu.Item key="2">nav 2</Menu.Item>
        <Menu.Item key="3">nav 3</Menu.Item>
      </Menu>
    </Header>)
  }
}

Widget.propTypes = {
  user: PropTypes.object.isRequired,
  site: PropTypes.object.isRequired
}

export default connect(state => ({user: state.currentUser, site: state.siteInfo}))(Widget)
