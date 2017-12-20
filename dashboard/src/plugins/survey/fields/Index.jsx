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
import {get, _delete, backend} from '../../../ajax'
import PlainText from '../../../components/PlainText'

class Widget extends Component {
  state = {
    items: []
  }
  componentDidMount() {
    const {formId} = this.props.match.params
    get(`/survey/fields?formId=${formId}`).then((rst) => {
      this.setState({items: rst})
    }).catch(message.error);
  }
  handleRemove = (id) => {
    const {formatMessage} = this.props.intl
    _delete(`/survey/fields/${id}`).then((rst) => {
      message.success(formatMessage({id: 'helpers.success'}))
      var items = this.state.items.filter((it) => it.id !== id)
      this.setState({items})
    }).catch(message.error)
  }
  render() {
    const {push} = this.props
    const {formId} = this.props.match.params
    return (<Layout breads={[
        {
          href: "/survey/forms",
          label: <FormattedMessage id={"survey.forms.index.title"}/>
        }, {
          href: `/survey/fields/${formId}`,
          label: <FormattedMessage id={"survey.fields.index.title"}/>
        }
      ]}>
      <Row>
        <Col>
          <Button onClick={(e) => push(`/survey/fields/new/${formId}`)} type='primary' shape="circle" icon="plus"/>
          <Table bordered={true} rowKey="id" dataSource={this.state.items} columns={[
              {
                title: <FormattedMessage id="attributes.label"/>,
                key: 'label',
                dataIndex: 'label'
              }, {
                title: <FormattedMessage id="attributes.type"/>,
                key: 'type',
                render: (text, record) => (<span>{record.type}[{record.sortOrder}]</span>)
              }, {
                title: 'Action',
                key: 'action',
                render: (text, record) => (<span>
                  <Button onClick={(e) => push(`/survey/fields/edit/${record.id}`)} shape="circle" icon="edit"/>
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
