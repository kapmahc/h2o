import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'
import {Col, Row} from 'antd'

import Layout from '../../layouts/application'
import Quill from '../../components/Quill'

class Widget extends Component {
  render() {
    const {site} = this.props
    return (<Layout breads={[]} title={{
        id: "nut.home.title"
      }}>
      <Row>
        <Col md={{
            span: 18,
            offset: 3
          }}>
          <Quill body={site.home.body}/>
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
