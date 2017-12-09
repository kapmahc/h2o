import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'
import {Link} from 'react-router-dom'
import {Card, Row, Col, Icon} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'

import dashboard from '../../dashboard'
import Layout from '../../layout'

class Panel extends Component {
  render() {
    const {item} = this.props
    return (<Card title={(<span >
        <Icon type={item.icon}/>&nbsp;
        <FormattedMessage id={item.label}/>
      </span>)}>
      {
        item.items.map((l) => (<p key={l.key}>
          <Link to={l.key}><FormattedMessage id={l.label}/></Link>
        </p>))
      }
    </Card>)

  }
}

Panel.propTypes = {
  item: PropTypes.object.isRequired
}

class Widget extends Component {
  render() {
    const {user} = this.props
    var zone = user.uid
      ? dashboard(user).filter((it) => it.items).map((it) => (<Col key={it.key} md={{
          span: 4
        }}><Panel item={it}/><br/></Col>))
      : (<Col md={{
          offset: 2,
          span: 8
        }}>
        <Card>
          {
            dashboard(user).filter((it) => !it.items).map((l) => (<p key={l.key}>
              <Icon type={l.icon}/>&nbsp;
              <Link to={l.key}>
                <FormattedMessage id={l.label}/>
              </Link>
            </p>))
          }
        </Card>
      </Col>)

    return (<Layout breads={[]}>
      <Row gutter={16}>
        {zone}
      </Row>
    </Layout>);
  }
}

Widget.propTypes = {
  push: PropTypes.func.isRequired,
  user: PropTypes.object.isRequired,
  info: PropTypes.object.isRequired,
  intl: intlShape.isRequired
}

const WidgetI = injectIntl(Widget)

export default connect(state => ({user: state.currentUser, info: state.siteInfo}), {push})(WidgetI)
