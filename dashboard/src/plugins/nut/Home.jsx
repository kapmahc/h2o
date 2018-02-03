import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'

import Layout from '../../layouts/application'

class Widget extends Component {
  render() {
    return (<Layout breads={[]} title={{
        id: "nut.home.title"
      }}>
      <div>home</div>
    </Layout>);
  }
}

Widget.propTypes = {
  push: PropTypes.func.isRequired,
  user: PropTypes.object.isRequired,
  info: PropTypes.object.isRequired
}

export default connect(state => ({user: state.currentUser, info: state.siteInfo}), {push})(Widget)
