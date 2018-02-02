import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {Layout, BackTop} from 'antd'
import {injectIntl, intlShape} from 'react-intl'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'

import Footer from '../layouts/Footer'
import LeftNavPanel from './LeftNavPanel'
import TopNavBar from './TopNavBar'
import {signOut} from '../actions'

const {Header, Content, Sider} = Layout

class Widget extends Component {

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
  signOut: PropTypes.func.isRequired,
  user: PropTypes.object.isRequired,
  info: PropTypes.object.isRequired,
  breads: PropTypes.array.isRequired,
  intl: intlShape.isRequired
}

const WidgetI = injectIntl(Widget)

export default connect(state => ({user: state.currentUser, info: state.siteInfo}), {push, signOut})(WidgetI)
