import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {
  Form,
  Row,
  Col,
  Input,
  Select,
  message
} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'

import Layout from '../../../../layout'
import {post, get} from '../../../../ajax'
import {Submit, orders, formItemLayout} from '../../../../components/form'

const FormItem = Form.Item
const Option = Select.Option

class Widget extends Component {
  componentDidMount() {
    const {setFieldsValue} = this.props.form
    const {id} = this.props.match.params
    if (id) {
      get(`/api/admin/friend-links/${id}`).then((rst) => setFieldsValue({title: rst.title, home: rst.home, logo: rst.logo, sortOrder: rst.sortOrder.toString()})).catch(message.error)
    } else {
      setFieldsValue({sortOrder: '0'})
    }
  }
  handleSubmit = (e) => {
    const {formatMessage} = this.props.intl
    const {push} = this.props
    const {id} = this.props.match.params
    e.preventDefault();
    this.props.form.validateFields((err, values) => {
      if (!err) {
        post(
          id
          ? `/api/admin/friend-links/${id}`
          : '/api/admin/friend-links',
        Object.assign({}, values, {
          sortOrder: parseInt(values.sortOrder, 10)
        })).then(() => {
          message.success(formatMessage({id: "messages.success"}))
          push('/admin/friend-links')
        }).catch(message.error);
      }
    });
  }
  render() {
    const {formatMessage} = this.props.intl
    const {getFieldDecorator} = this.props.form
    const {id} = this.props.match.params
    return (<Layout breads={[
        {
          href: '/admin/friend-links',
          label: <FormattedMessage id='nut.admin.friend-links.index.title'/>
        },
        id
          ? {
            href: `/admin/friend-links/edit/${id}`,
            label: (<FormattedMessage id={"buttons.edit"} values={{
                id: id
              }}/>)
          }
          : {
            href: "/admin/friend-links/new",
            label: <FormattedMessage id={"buttons.new"}/>
          }
      ]}>
      <Row>
        <Col md={{
            span: 12,
            offset: 2
          }}>
          <Form onSubmit={this.handleSubmit}>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "attributes.title" />} hasFeedback={true}>
              {
                getFieldDecorator('title', {
                  rules: [
                    {
                      required: true,
                      message: formatMessage({id: "errors.empty"})
                    }
                  ]
                })(<Input/>)
              }
            </FormItem>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "nut.attributes.friend-link.home" />} hasFeedback={true}>
              {
                getFieldDecorator('home', {
                  rules: [
                    {
                      required: true,
                      message: formatMessage({id: "errors.empty"})
                    }
                  ]
                })(<Input/>)
              }
            </FormItem>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "nut.attributes.friend-link.logo" />} hasFeedback={true}>
              {
                getFieldDecorator('logo', {
                  rules: [
                    {
                      required: true,
                      message: formatMessage({id: "errors.empty"})
                    }
                  ]
                })(<Input/>)
              }
            </FormItem>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "attributes.sortOrder" />}>
              {
                getFieldDecorator('sortOrder')(<Select>
                  {orders(10).map((p) => (<Option key={p} value={p}>{p}</Option>))}
                </Select>)
              }
            </FormItem>

            <Submit/>
          </Form>
        </Col>
      </Row>
    </Layout>);
  }
}

Widget.propTypes = {
  intl: intlShape.isRequired,
  push: PropTypes.func.isRequired
}

const WidgetF = Form.create()(injectIntl(Widget))

export default connect(state => ({}), {
  push
},)(WidgetF)
