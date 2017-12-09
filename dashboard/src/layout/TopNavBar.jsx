import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {Breadcrumb, Icon} from 'antd'
import {Link} from 'react-router-dom'

class Widget extends Component {
  render() {
    return (
      <Breadcrumb style={{
        margin: '16px 0'
      }}>
        <Breadcrumb.Item >
          <Link to="/">
            <Icon type="home"/>
          </Link>
        </Breadcrumb.Item>
        {this.props.items.map((it) => (
          <Breadcrumb.Item key={it.href}>
            <Link to={it.href}>
              {it.label}
            </Link>
          </Breadcrumb.Item>
        ))}
      </Breadcrumb>
    );
  }
}

Widget.propTypes = {
  items: PropTypes.array.isRequired
}

export default Widget
