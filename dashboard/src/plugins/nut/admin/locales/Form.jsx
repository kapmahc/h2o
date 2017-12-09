import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {Form, Row, Col, Input, message} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'

import Layout from '../../../../layout'
import {post, get} from '../../../../ajax'
import {Submit, formItemLayout} from '../../../../components/form'

const FormItem = Form.Item

class Widget extends Component {
  componentDidMount() {
    const {setFieldsValue} = this.props.form
    const {code} = this.props.match.params
    if (code) {
      get(`/api/admin/locales/${code}`).then(setFieldsValue).catch(message.error)
    }
  }
  handleSubmit = (e) => {
    const {formatMessage} = this.props.intl
    const {push} = this.props
    e.preventDefault();
    this.props.form.validateFields((err, values) => {
      if (!err) {
        post('/api/admin/locales', values).then(() => {
          message.success(formatMessage({id: "messages.success"}))
          push('/admin/locales')
        }).catch(message.error);
      }
    });
  }
  render() {
    const {formatMessage} = this.props.intl
    const {getFieldDecorator} = this.props.form
    const {code} = this.props.match.params
    return (<Layout breads={[
        {
          href: '/admin/locales',
          label: <FormattedMessage id='nut.admin.locales.index.title'/>
        },
        code
          ? {
            href: `/admin/locales/edit/${code}`,
            label: (<FormattedMessage id={"buttons.edit"} values={{
                id: code
              }}/>)
          }
          : {
            href: "/admin/locales/new",
            label: <FormattedMessage id={"buttons.new"}/>
          }
      ]}>
      <Row>
        <Col md={{
            span: 12,
            offset: 2
          }}>
          <Form onSubmit={this.handleSubmit}>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "nut.attributes.locale.code" />} hasFeedback={true}>
              {
                getFieldDecorator('code', {
                  rules: [
                    {
                      required: true,
                      message: formatMessage({id: "errors.empty"})
                    }
                  ]
                })(<Input/>)
              }
            </FormItem>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "nut.attributes.locale.message" />} hasFeedback={true}>
              {
                getFieldDecorator('message', {
                  rules: [
                    {
                      required: true,
                      message: formatMessage({id: "errors.empty"})
                    }
                  ]
                })(<Input.TextArea rows={6}/>)
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
