import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {Form, Row, Col, Input, message} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'

import Layout from '../../../layout'
import {post} from '../../../ajax'
import {Submit, formItemLayout} from '../../../components/form'

const FormItem = Form.Item

class Widget extends Component {
  handleSubmit = (e) => {
    const {push, action} = this.props
    const {formatMessage} = this.props.intl
    e.preventDefault();
    this.props.form.validateFields((err, values) => {
      if (!err) {
        post(`/users/${action}`, values).then(() => {
          message.info(formatMessage({id: `nut.users.${action}.success`}))
          push('/users/sign-in')
        }).catch(message.error);
      }
    });
  }
  render() {
    const {action} = this.props
    const {formatMessage} = this.props.intl
    const {getFieldDecorator} = this.props.form
    return (<Layout breads={[{
          href: `/users/${action}`,
          label: <FormattedMessage id={`nut.users.${action}.title`}/>
        }
      ]}>
      <Row>
        <Col md={{
            span: 12,
            offset: 2
          }}>
          <Form onSubmit={this.handleSubmit}>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "attributes.email" />} hasFeedback={true}>
              {
                getFieldDecorator('email', {
                  rules: [
                    {
                      type: 'email',
                      message: formatMessage({id: "errors.not-valid-email"})
                    }, {
                      required: true,
                      message: formatMessage({id: "errors.empty-email"})
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
  action: PropTypes.string.isRequired,
  push: PropTypes.func.isRequired
}

const WidgetF = Form.create()(injectIntl(Widget))

export default connect(state => ({}), {
  push
},)(WidgetF)
