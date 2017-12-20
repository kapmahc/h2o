import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {Layout, message, BackTop} from 'antd'
import {injectIntl, intlShape} from 'react-intl'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'

import Footer from './Footer'
import LeftNavPanel from './LeftNavPanel'
import TopNavBar from './TopNavBar'
import {signIn, signOut, refresh} from '../actions'
import {get, token} from '../ajax'

const {Header, Content, Sider} = Layout

class Widget extends Component {
  componentDidMount() {
    const {signIn, refresh, info} = this.props
    var tkn = token()
    if (tkn) {
      signIn(tkn)
    }
    if (info.languages.length === 0) {
      get('/layout').then((rst) => refresh(rst)).catch(message.error)
    }
  }
  render() {
    const {children, breads, info} = this.props

    return (<Layout>
      <Sider breakpoint="lg" collapsedWidth="0" onCollapse={(collapsed, type) => {
          console.log(collapsed, type);
        }}>
        <div className="logo">{info.subhead}</div>
        <LeftNavPanel/>
      </Sider>
      <Layout>
        <Header style={{
            background: '#fff',
            padding: 0
          }}/>
        <Content style={{
            margin: '0 16px'
          }}>
          <TopNavBar items={breads}/>
          <div style={{
              padding: 24,
              background: '#fff',
              minHeight: 360
            }}>
            {children}
          </div>
        </Content>
        <Footer/>
        <BackTop/>
      </Layout>
    </Layout>);
  }
}

Widget.propTypes = {
  children: PropTypes.node.isRequired,
  push: PropTypes.func.isRequired,
  refresh: PropTypes.func.isRequired,
  signIn: PropTypes.func.isRequired,
  signOut: PropTypes.func.isRequired,
  user: PropTypes.object.isRequired,
  info: PropTypes.object.isRequired,
  breads: PropTypes.array.isRequired,
  intl: intlShape.isRequired
}

const WidgetI = injectIntl(Widget)

export default connect(state => ({user: state.currentUser, info: state.siteInfo}), {push, signIn, refresh, signOut})(WidgetI)
