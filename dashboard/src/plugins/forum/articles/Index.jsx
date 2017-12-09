import React, {Component} from 'react'
import {
  Row,
  Col,
  Table,
  Popconfirm,
  Button,
  message
} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'
import {CopyToClipboard} from 'react-copy-to-clipboard'

import Layout from '../../../layout'
import {get, _delete} from '../../../ajax'
import {UEditor} from '../../../components/form'

class Widget extends Component {
  state = {
    items: []
  }
  componentDidMount() {
    get('/api/forum/articles').then((rst) => {
      this.setState({items: rst})
    }).catch(message.error);
  }
  handleRemove = (id) => {
    const {formatMessage} = this.props.intl
    _delete(`/api/forum/articles/${id}`).then((rst) => {
      message.success(formatMessage({id: 'messages.success'}))
      var items = this.state.items.filter((it) => it.id !== id)
      this.setState({items})
    }).catch(message.error)
  }
  render() {
    const {push} = this.props
    return (<Layout breads={[{
          href: "/forum/articles",
          label: <FormattedMessage id={"forum.articles.index.title"}/>
        }
      ]}>
      <Row>
        <Col>
          <Button onClick={(e) => push('/forum/articles/new')} type='primary' shape="circle" icon="plus"/>
          <Table bordered={true} rowKey="id" dataSource={this.state.items} columns={[
              {
                title: <FormattedMessage id="attributes.title"/>,
                key: 'title',
                render: (text, record) => (<a href={`/forum/articles/show/${record.id}`} target="_blank">{record.title}</a>)
              }, {
                title: 'Action',
                key: 'action',
                render: (text, record) => (<span>
                  <CopyToClipboard text={`/forum/articles/show/${record.id}`}><Button shape="circle" icon="copy"/></CopyToClipboard>
                  <Button onClick={(e) => push(`/forum/articles/edit/${record.id}`)} shape="circle" icon="edit"/>
                  <UEditor target={record.id} action="/forum/articles/body/edit"/>
                  <Popconfirm title={<FormattedMessage id = "messages.are-you-sure" />} onConfirm={(e) => this.handleRemove(record.id)}>
                    <Button type="danger" shape="circle" icon="delete"/>
                  </Popconfirm>
                </span>)
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

export default connect(state => ({}), {
  push
},)(WidgetI)
