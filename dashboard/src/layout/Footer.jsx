import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {Layout} from 'antd'
import {connect} from 'react-redux'
import {FormattedMessage} from 'react-intl'
import {set as setLocale} from '../intl'

const {Footer} = Layout

class Widget extends Component {
  switchLanguage = (l) => {
    setLocale(l)
    window.location.reload()
  }
  render() {
    const {site} = this.props
    return (<Footer style={{
        textAlign: 'center'
      }}>
      &copy;{site.copyright}
      {
        site.languages.map((l, i) => (<a style={{
            paddingLeft: '8px'
          }} key={i} onClick={(e) => this.switchLanguage(l)}><FormattedMessage id={`languages.${l}`}/></a>))
      }
    </Footer>);
  }
}
Widget.propTypes = {
  site: PropTypes.object.isRequired
}

export default connect(state => ({site: state.siteInfo}))(Widget)
