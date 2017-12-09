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
import {Submit, Quill, orders, formItemLayout} from '../../../../components/form'

const FormItem = Form.Item
const Option = Select.Option

class Widget extends Component {
  state = {
    summary: ''
  }
  componentDidMount() {
    const {setFieldsValue} = this.props.form
    const {id} = this.props.match.params
    if (id) {
      get(`/api/admin/cards/${id}`).then((rst) => {
        setFieldsValue({
          title: rst.title,
          action: rst.action,
          logo: rst.logo,
          href: rst.href,
          sortOrder: rst.sortOrder.toString(),
          loc: rst.loc
        })
        this.setState({summary: rst.summary})
      }).catch(message.error)
    } else {
      setFieldsValue({sortOrder: '0'})
    }
  }
  handleChange = (value) => {
    this.setState({summary: value})
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
          ? `/api/admin/cards/${id}`
          : '/api/admin/cards',
        Object.assign({}, values, {
          sortOrder: parseInt(values.sortOrder, 10),
          type: 'html',
          summary: this.state.summary
        })).then(() => {
          message.success(formatMessage({id: "messages.success"}))
          push('/admin/cards')
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
          href: '/admin/cards',
          label: <FormattedMessage id='nut.admin.cards.index.title'/>
        },
        id
          ? {
            href: `/admin/cards/edit/${id}`,
            label: (<FormattedMessage id={"buttons.edit"} values={{
                id: id
              }}/>)
          }
          : {
            href: "/admin/cards/new",
            label: <FormattedMessage id={"buttons.new"}/>
          }
      ]}>
      <Row>
        <Col md={{
            span: 18
          }}>
          <Form onSubmit={this.handleSubmit}>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "attributes.loc" />} hasFeedback={true}>
              {
                getFieldDecorator('loc', {
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
            <FormItem {...formItemLayout} label={<FormattedMessage id = "attributes.logo" />} hasFeedback={true}>
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
            <FormItem {...formItemLayout} label={<FormattedMessage id = "attributes.summary" />}>
              <Quill value={this.state.summary} onChange={this.handleChange}/>
            </FormItem>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "attributes.action" />} hasFeedback={true}>
              {
                getFieldDecorator('action', {
                  rules: [
                    {
                      required: true,
                      message: formatMessage({id: "errors.empty"})
                    }
                  ]
                })(<Input/>)
              }
            </FormItem>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "attributes.href" />} hasFeedback={true}>
              {
                getFieldDecorator('href', {
                  rules: [
                    {
                      required: true,
                      message: formatMessage({id: "errors.empty"})
                    }
                  ]
                })(<Input/>)
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
