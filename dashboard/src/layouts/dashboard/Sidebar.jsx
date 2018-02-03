import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'
import {Menu, Icon, Modal, message} from 'antd'
import {push} from 'react-router-redux'
import {Link} from 'react-router-dom'

import {signOut} from '../../actions'
import {_delete} from '../../ajax'
import {Authorized, USER, ADMIN} from '../../auth'

import plugins from '../../plugins'

const SubMenu = Menu.SubMenu
const MenuItem = Menu.Item
const confirm = Modal.confirm

class Widget extends Component {
  handleMenu = ({key}) => {
    const {push, signOut} = this.props
    const {formatMessage} = this.props.intl

    switch (key) {
      case "users.sign-out":
        confirm({
          title: formatMessage({id: "helpers.are-you-sure"}),
          onOk() {
            _delete('/users/sign-out').then(() => {
              signOut()
              push('/users/sign-in')
              message.success(formatMessage({id: 'helpers.success'}))
            }).catch(message.error)
          }
        });
        break
      default:
        break
    }
  };
  render() {
    return (<Menu theme="dark" mode="inline" onClick={this.handleMenu}>
      {
        plugins.menus.map((it) => Authorized.check(it.roles, (<SubMenu key={it.href} title={(<span >
            <Icon type={it.icon}/>
            <FormattedMessage id={it.label}/>
          </span>)}>
          {
            it.items.map((l) => Authorized.check(l.roles || it.roles, (<MenuItem key={`${it.href}-${l.href}`}>
              <Link to={l.href}>
                <FormattedMessage id={l.label}/>
              </Link>
            </MenuItem>)))
          }
        </SubMenu>)))
      }
      {
        Authorized.check([
          ADMIN, USER
        ], (<MenuItem key='users.sign-out'>
          <Icon type='logout'/>
          <FormattedMessage id='nut.users.sign-out.title'/>
        </MenuItem>))
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
