import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {Form, Row, Col, Input, message} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'

import Layout from '../../../../layouts/dashboard'
import {post, get, backend} from '../../../../ajax'
import {Submit, Quill, formItemLayout} from '../../../../components/form'
import {ADMIN} from '../../../../auth'

const FormItem = Form.Item

class Widget extends Component {
  state = {
    body: ''
  }
  componentDidMount() {
    const {setFieldsValue} = this.props.form
    get('/admin/site/donate').then((rst) => {
      setFieldsValue({paypal: rst.paypal})
      this.setState({body: rst.body})
    }).catch(message.error)
  }
  handleChange = (value) => {
    this.setState({body: value})
  }
  handleSubmit = (e) => {
    const {formatMessage} = this.props.intl
    e.preventDefault();
    this.props.form.validateFields((err, values) => {
      if (!err) {
        post('/admin/site/donate', Object.assign({}, values, {
          type: 'html',
          body: this.state.body
        })).then(() => {
          message.success(formatMessage({id: "helpers.success"}))
        }).catch(message.error);
      }
    });
  }
  render() {
    const {getFieldDecorator} = this.props.form
    const title = {
      id: "nut.admin.site.donate.title"
    }
    return (<Layout breads={[{
          href: "/admin/site/donate",
          label: title
        }
      ]} title={title} roles={[ADMIN]}>
      <Row>
        <Col md={{
            span: 12,
            offset: 2
          }}>
          <Form onSubmit={this.handleSubmit}>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "nut.admin.site.donate.paypal" />} hasFeedback={true}>
              {getFieldDecorator('paypal')(<Input/>)}
            </FormItem>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "attributes.body" />}>
              <Quill value={this.state.body} onChange={this.handleChange}/>
            </FormItem>
            <Submit>
              &nbsp; &nbsp;
              <a target="_blank" href={backend("/htdocs/donate")}><FormattedMessage id="buttons.view"/></a>
            </Submit>
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
