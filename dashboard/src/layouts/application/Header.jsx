import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {Layout, Menu, Icon} from 'antd';
import {connect} from 'react-redux'
import {Link} from 'react-router-dom'

const {Header} = Layout;
const SubMenu = Menu.SubMenu;
const MenuItemGroup = Menu.ItemGroup;

class Widget extends Component {

  render() {
    const width = document.body.offsetWidth
    const {site} = this.props
    var children = null
    if (width > 1024) {
      children = [(<Menu.Item key="home">
        <Link to="/">{site.subhead}</Link>
      </Menu.Item>)].concat(site.header.map((it, id) => (
        it.items.length > 0
        ? (<SubMenu key={id} title={(<span>{it.label}<Icon type="down"/></span>)}>
          {
            it.items.map((jt, jd) => (<Menu.Item key={`${id}.${jd}`}>
              <a href={it.href} target="_blank">{it.label}</a>
            </Menu.Item>))
          }
        </SubMenu>)
        : (<Menu.Item key={id}>
          <a href={it.href} target="_blank">{it.label}</a>
        </Menu.Item>))))
    } else {
      children = (<SubMenu title={site.subhead}>
        {
          site.header.map((it, id) => (
            it.items.length > 0
            ? (<MenuItemGroup key={id} title={it.label}>
              {
                it.items.map((jt, jd) => (<Menu.Item key={`${id}.${jd}`}>
                  <a href={it.href} target="_blank">{it.label}</a>
                </Menu.Item>))
              }
            </MenuItemGroup>)
            : (<Menu.Item key={id}>
              <a href={it.href} target="_blank">{it.label}</a>
            </Menu.Item>)))
        }
      </SubMenu>)
    }
    return (<Header>
      <div className="logo"/>
      <Menu theme="dark" mode="horizontal" defaultSelectedKeys={[]} style={{
          lineHeight: '64px'
        }}>
        {children}
      </Menu>
    </Header>)
  }
}

Widget.propTypes = {
  user: PropTypes.object.isRequired,
  site: PropTypes.object.isRequired
}

export default connect(state => ({user: state.currentUser, site: state.siteInfo}))(Widget)
