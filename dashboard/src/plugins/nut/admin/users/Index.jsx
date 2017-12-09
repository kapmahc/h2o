import React, {Component} from 'react'
import {Row, Col, Table, message} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'

import Layout from '../../../../layout'
import {get} from '../../../../ajax'

class Widget extends Component {
  state = {
    items: []
  }
  componentDidMount() {
    get('/api/admin/users').then((rst) => {
      this.setState({items: rst})
    }).catch(message.error);
  }
  render() {
    return (<Layout breads={[{
          href: "/admin/users",
          label: <FormattedMessage id={"nut.admin.users.index.title"}/>
        }
      ]}>
      <Row>
        <Col md={{
            span: 18,
            offset: 2
          }}>
          <Table bordered={true} rowKey="id" dataSource={this.state.items} columns={[
              {
                title: <FormattedMessage id="attributes.user"/>,
                key: 'user',
                render: (text, record) => (<span>{record.name}&lt;{record.email}&gt;[{record.signInCount}]</span>)
              }, {
                title: <FormattedMessage id="nut.attributes.user.lastSignIn"/>,
                key: 'lastSignIn',
                render: (text, record) => (<span>{record.lastSignInAt}[{record.lastSignInIp}]</span>)
              }, {
                title: <FormattedMessage id="nut.attributes.user.currentSignIn"/>,
                key: 'currentSignIn',
                render: (text, record) => (<span>{record.currentSignInAt}[{record.currentSignInIp}]</span>)
              }
            ]}/>
        </Col>
      </Row>
    </Layout>);
  }
}

Widget.propTypes = {
  intl: intlShape.isRequired
}

const WidgetI = injectIntl(Widget)

export default connect(state => ({}), {},)(WidgetI)
