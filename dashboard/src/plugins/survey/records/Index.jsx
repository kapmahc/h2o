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

import Layout from '../../../layout'
import {get, _delete} from '../../../ajax'

class Widget extends Component {
  state = {
    items: []
  }
  componentDidMount() {
    const {formId} = this.props.match.params
    get(`/survey/records?formId=${formId}`).then((rst) => {
      this.setState({items: rst})
    }).catch(message.error);
  }
  handleRemove = (id) => {
    const {formatMessage} = this.props.intl
    _delete(`/survey/records/${id}`).then((rst) => {
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
          href: `/survey/records/${formId}`,
          label: <FormattedMessage id={"survey.records.index.title"}/>
        }
      ]}>
      <Row>
        <Col>
          <Table bordered={true} rowKey="id" dataSource={this.state.items} columns={[
              {
                title: <FormattedMessage id="attributes.email"/>,
                key: 'email',
                dataIndex: 'email'
              }, {
                title: <FormattedMessage id="attributes.value"/>,
                key: 'value',
                render: (text, record) => (<span>{record.value}</span>)
              }, {
                title: 'Action',
                key: 'action',
                render: (text, record) => (<span>
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
