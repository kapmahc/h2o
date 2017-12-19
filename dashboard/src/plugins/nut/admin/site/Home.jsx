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
import {Submit, formItemLayout} from '../../../../components/form'

const FormItem = Form.Item
const Option = Select.Option

class Widget extends Component {
  componentDidMount() {
    const {setFieldsValue} = this.props.form
    get('/api/admin/site/home').then((rst) => setFieldsValue(rst)).catch(message.error)
  }
  handleSubmit = (e) => {
    const {formatMessage} = this.props.intl
    e.preventDefault();
    this.props.form.validateFields((err, values) => {
      if (!err) {
        post('/api/admin/site/home', values).then(() => {
          message.success(formatMessage({id: "helpers.success"}))
        }).catch(message.error);
      }
    });
  }
  render() {
    const {formatMessage} = this.props.intl
    const {getFieldDecorator} = this.props.form
    return (<Layout breads={[{
          href: "/admin/site/home",
          label: <FormattedMessage id={"nut.admin.site.home.title"}/>
        }
      ]}>
      <Row>
        <Col md={{
            span: 16,
            offset: 2
          }}>
          <Form onSubmit={this.handleSubmit}>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "nut.admin.site.home.favicon" />} hasFeedback={true}>
              {
                getFieldDecorator('favicon', {
                  rules: [
                    {
                      required: true,
                      message: formatMessage({id: "errors.empty"})
                    }
                  ]
                })(<Input/>)
              }
            </FormItem>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "nut.admin.site.home.theme" />}>
              {
                getFieldDecorator('theme')(<Select>
                  {["off-canvas", "carousel", "album"].map((p) => (<Option key={p} value={p}>{p}</Option>))}
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
