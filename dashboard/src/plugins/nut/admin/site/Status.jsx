import React, {Component} from 'react'
import {Row, Col, Collapse, Table, message} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import PropTypes from 'prop-types'
import SyntaxHighlighter from 'react-syntax-highlighter'
import {docco} from 'react-syntax-highlighter/styles/hljs'

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
      title: <FormattedMessage id="attributes.value"/>,
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
    database: {},
    redis: "",
    routes: [],
    jobber: {},
    network: {}
  }
  componentDidMount() {
    get('/admin/site/status').then((rst) => {
      this.setState(rst)
    }).catch(message.error);
  }
  render() {
    const {
      redis,
      os,
      network,
      database,
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
            <Panel key="database" header={(<FormattedMessage id="nut.admin.site.status.database"/>)}>
              <Hash item={database}/>
            </Panel>
            <Panel key="jobber" header={(<FormattedMessage id="nut.admin.site.status.jobber"/>)}>
              <Hash item={jobber}/>
            </Panel>
            <Panel key="redis" header={(<FormattedMessage id="nut.admin.site.status.redis"/>)}>
              <SyntaxHighlighter style={docco}>{redis}</SyntaxHighlighter>
            </Panel>
            <Panel header={(<FormattedMessage id='nut.admin.site.status.routes'/>)} key="routes">
              <Table rowKey="id" dataSource={routes.map((v, k) => {
                  return {id: k, method: v.Method, path: v.Path}
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
