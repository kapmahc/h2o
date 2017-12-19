import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {Form, Row, Col, Input, message} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'

import Layout from '../../../layout'
import {post, get} from '../../../ajax'
import {Submit, formItemLayout} from '../../../components/form'

const FormItem = Form.Item

class Widget extends Component {
  componentDidMount() {
    const {setFieldsValue} = this.props.form
    const {user} = this.props
    if (user.uid) {
      get('/api/users/profile').then((rst) => setFieldsValue({name: rst.name, email: rst.email})).catch(message.error)
    }
  }
  handleSubmit = (e) => {
    const {formatMessage} = this.props.intl
    e.preventDefault();
    this.props.form.validateFields((err, values) => {
      if (!err) {
        post('/api/users/profile', values).then(() => {
          message.success(formatMessage({id: "helpers.success"}))
        }).catch(message.error);
      }
    });
  }
  render() {
    const {formatMessage} = this.props.intl
    const {getFieldDecorator} = this.props.form
    return (<Layout breads={[{
          href: "/users/profile",
          label: <FormattedMessage id={"nut.users.profile.title"}/>
        }
      ]}>
      <Row>
        <Col md={{
            span: 12,
            offset: 2
          }}>
          <Form onSubmit={this.handleSubmit}>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "attributes.email" />}>
              {getFieldDecorator('email', {})(<Input disabled={true}/>)}
            </FormItem>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "attributes.username" />} hasFeedback={true}>
              {
                getFieldDecorator('name', {
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
  user: PropTypes.object.isRequired,
  push: PropTypes.func.isRequired
}

const WidgetF = Form.create()(injectIntl(Widget))

export default connect(state => ({user: state.currentUser}), {
  push
},)(WidgetF)
