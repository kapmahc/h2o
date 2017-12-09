import React, {Component} from 'react'
import {Row, Col, Collapse, Table, message} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import PropTypes from 'prop-types'

import Layout from '../../../../layout'
import {get} from '../../../../ajax'

const Panel = Collapse.Panel

const Hash = ({item}) => (<Table rowKey="key" dataSource={Object.entries(item).map((v) => {
    return {key: v[0], val: v[1]}
  })} columns={[
    {
      title: <FormattedMessage id="attributes.key"/>,
      key: 'key',
      dataIndex: 'key'
    }, {
      title: <FormattedMessage id="attributes.val"/>,
      key: 'val',
      dataIndex: 'val'
    }
  ]}/>)

Hash.propTypes = {
  item: PropTypes.object.isRequired
}

class Widget extends Component {
  state = {
    os: {},
    postgresql: {},
    redis: [],
    routes: [],
    jobber: {},
    network: {}
  }
  componentDidMount() {
    get('/api/admin/site/status').then((rst) => {
      this.setState(rst)
    }).catch(message.error);
  }
  render() {
    const {
      redis,
      os,
      network,
      postgresql,
      jobber,
      routes
    } = this.state

    return (<Layout breads={[{
          href: "/admin/site/status",
          label: <FormattedMessage id={"nut.admin.site.status.title"}/>
        }
      ]}>
      <Row>
        <Col md={{
            span: 16,
            offset: 2
          }}>
          <Collapse>
            <Panel key="os" header={(<FormattedMessage id="nut.admin.site.status.os"/>)}>
              <Hash item={os}/>
            </Panel>
            <Panel key="network" header={(<FormattedMessage id="nut.admin.site.status.network"/>)}>
              <Hash item={network}/>
            </Panel>
            <Panel key="postgresql" header={(<FormattedMessage id="nut.admin.site.status.postgresql"/>)}>
              <Hash item={postgresql}/>
            </Panel>
            <Panel key="jobber" header={(<FormattedMessage id="nut.admin.site.status.jobber"/>)}>
              <Hash item={jobber}/>
            </Panel>
            <Panel key="redis" header={(<FormattedMessage id="nut.admin.site.status.redis"/>)}>
              {redis.map((v, k) => (<div key={k}>{v}</div>))}
            </Panel>
            <Panel header={(<FormattedMessage id='nut.admin.site.status.routes'/>)} key="routes">
              <Table rowKey="id" dataSource={routes.map((v, k) => {
                  return Object.assign({}, v, {id: k})
                })} columns={[
                  {
                    title: 'METHOD',
                    key: 'method',
                    dataIndex: 'method'
                  }, {
                    title: 'PATH',
                    key: 'path',
                    dataIndex: 'path'
                  }
                ]}/>
            </Panel>
          </Collapse>
        </Col>
      </Row>
    </Layout>);
  }
}
Widget.propTypes = {
  intl: intlShape.isRequired
}

const WidgetI = injectIntl(Widget)

export default WidgetI
