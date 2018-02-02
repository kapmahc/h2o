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

import Layout from '../../../layouts/dashboard'
import {USER, ADMIN} from '../../../auth'
import {get, _delete, backend} from '../../../ajax'

class Widget extends Component {
  state = {
    items: []
  }
  componentDidMount() {
    get('/survey/forms').then((rst) => {
      this.setState({items: rst})
    }).catch(message.error);
  }
  handleRemove = (id) => {
    const {formatMessage} = this.props.intl
    _delete(`/survey/forms/${id}`).then((rst) => {
      message.success(formatMessage({id: 'helpers.success'}))
      var items = this.state.items.filter((it) => it.id !== id)
      this.setState({items})
    }).catch(message.error)
  }
  render() {
    const {push} = this.props
    const title = {
      id: "survey.forms.index.title"
    }
    return (<Layout breads={[{
          href: "/survey/forms",
          label: title
        }
      ]} title={title} roles={[USER, ADMIN]}>
      <Row>
        <Col>
          <Button onClick={(e) => push('/survey/forms/new')} type='primary' shape="circle" icon="plus"/>
          <Table bordered={true} rowKey="id" dataSource={this.state.items} columns={[
              {
                title: <FormattedMessage id="attributes.title"/>,
                key: 'title',
                render: (text, record) => (<a href={backend(`/survey/htdocs/forms/apply/${record.id}`)} target="_blank">{record.title}</a>)
              }, {
                title: 'Action',
                key: 'action',
                render: (text, record) => (<span>
                  <CopyToClipboard text={backend(`/survey/htdocs/forms/apply/${record.id}`)}><Button shape="circle" icon="copy"/></CopyToClipboard>
                  <Button onClick={(e) => push(`/survey/forms/edit/${record.id}`)} shape="circle" icon="edit"/>
                  <Button onClick={(e) => push(`/survey/fields/${record.id}`)} shape="circle" icon="profile"/>
                  <Button onClick={(e) => push(`/survey/records/${record.id}`)} shape="circle" icon="export"/>
                  <Popconfirm title={<FormattedMessage id = "helpers.are-you-sure" />} onConfirm={(e) => this.handleRemove(record.id)}>
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
