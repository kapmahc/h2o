import React, {Component} from 'react'
import {Form, Row, Col, Input, message} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'

import Layout from '../../../layout'
import {post} from '../../../ajax'
import {Submit, formItemLayout} from '../../../components/form'

const FormItem = Form.Item

class Widget extends Component {
  handleSubmit = (e) => {
    const {formatMessage} = this.props.intl
    const {setFieldsValue} = this.props.form
    e.preventDefault();
    this.props.form.validateFields((err, values) => {
      if (!err) {
        post('/api/users/change-password', values).then(() => {
          message.success(formatMessage({id: "helpers.success"}))
          setFieldsValue({currentPassword: "", newPassword: "", passwordConfirmation: ""})
        }).catch(message.error);
      }
    });
  }
  checkPassword = (rule, value, callback) => {
    const {formatMessage} = this.props.intl
    const {getFieldValue} = this.props.form
    if (value && value !== getFieldValue('newPassword')) {
      callback(formatMessage({id: "errors.passwords-not-match"}));
    } else {
      callback();
    }
  }
  render() {
    const {formatMessage} = this.props.intl
    const {getFieldDecorator} = this.props.form
    return (<Layout breads={[{
          href: "/users/change-password",
          label: <FormattedMessage id={"nut.users.change-password.title"}/>
        }
      ]}>
      <Row>
        <Col md={{
            span: 12,
            offset: 2
          }}>
          <Form onSubmit={this.handleSubmit}>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "attributes.currentPassword" />} hasFeedback={true}>
              {
                getFieldDecorator('currentPassword', {
                  rules: [
                    {
                      required: true,
                      message: formatMessage({id: "errors.empty-password"})
                    }
                  ]
                })(<Input type="password"/>)
              }
            </FormItem>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "attributes.newPassword" />} hasFeedback={true}>
              {
                getFieldDecorator('newPassword', {
                  rules: [
                    {
                      required: true,
                      message: formatMessage({id: "errors.empty-password"})
                    }
                  ]
                })(<Input type="password"/>)
              }
            </FormItem>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "attributes.passwordConfirmation" />} hasFeedback={true}>
              {
                getFieldDecorator('passwordConfirmation', {
                  rules: [
                    {
                      required: true,
                      message: formatMessage({id: "errors.empty"})
                    }, {
                      validator: this.checkPassword
                    }
                  ]
                })(<Input type="password"/>)
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
  intl: intlShape.isRequired
}

const WidgetF = Form.create()(injectIntl(Widget))

export default connect(state => ({}), {},)(WidgetF)
