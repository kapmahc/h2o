import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'
import {Menu, Icon, Modal, message} from 'antd'
import {push} from 'react-router-redux'

import {signOut} from '../actions'
import {_delete} from '../ajax'

import dashboard from '../dashboard'

const SubMenu = Menu.SubMenu
const confirm = Modal.confirm

class Widget extends Component {
  handleMenu = ({key}) => {
    const {push, signOut} = this.props
    const {formatMessage} = this.props.intl

    switch (key) {
      case "/users/sign-out":
        confirm({
          title: formatMessage({id: "messages.are-you-sure"}),
          onOk() {
            _delete('/api/users/sign-out').then(() => {
              signOut()
              push('/users/sign-in')
              message.success(formatMessage({id: 'messages.success'}))
            }).catch(message.error)
          }
        });
        break
      default:
        push(key)
    }
  };
  render() {
    const {user} = this.props
    return (<Menu theme="dark" mode="inline" onClick={this.handleMenu}>
      {
        dashboard(user).map(
          (item) => item.items
          ? (<SubMenu key={item.key} title={(<span >
              <Icon type={item.icon}/>
              <FormattedMessage id={item.label}/>
            </span>)}>
            {item.items.map((l) => (<Menu.Item key={l.key}><FormattedMessage id={l.label}/></Menu.Item>))}
          </SubMenu>)
          : (<Menu.Item key={item.key}>
            <Icon type={item.icon}/>
            <FormattedMessage id={item.label}/>
          </Menu.Item>))
      }
    </Menu>)
  }
}

Widget.propTypes = {
  intl: intlShape.isRequired,
  push: PropTypes.func.isRequired,
  signOut: PropTypes.func.isRequired,
  user: PropTypes.object.isRequired
}

const WidgetI = injectIntl(Widget)
export default connect(state => ({user: state.currentUser}), {push, signOut})(WidgetI)
