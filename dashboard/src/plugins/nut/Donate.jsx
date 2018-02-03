import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'
import {Col, Row, message} from 'antd'
import {FormattedMessage} from 'react-intl'

import {get} from '../../ajax'
import Layout from '../../layouts/application'
import Quill from '../../components/Quill'

class Widget extends Component {
  state = {
    body: '',
    paypal: ''
  }
  componentDidMount() {
    get('/donate').then((rst) => {
      this.setState(rst)
    }).catch(message.error);
  }
  render() {
    return (<Layout breads={[]} title={{
        id: "nut.donate.title"
      }}>
      <Row>
        <Col md={{
            span: 18,
            offset: 3
          }}>
          <Quill body={this.state.body}/>
          <FormattedMessage tagName="h3" id="nut.donate.by-paypal"/>
          <div dangerouslySetInnerHTML={{
              __html: this.state.paypal
            }}/>
        </Col>
      </Row>
    </Layout>);
  }
}

Widget.propTypes = {
  push: PropTypes.func.isRequired,
  user: PropTypes.object.isRequired,
  site: PropTypes.object.isRequired
}

export default connect(state => ({user: state.currentUser, site: state.siteInfo}), {push})(Widget)
