import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {
  Form,
  Row,
  Col,
  Input,
  Select,
  Checkbox,
  message
} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'

import Layout from '../../../layouts/dashboard'
import {post, get} from '../../../ajax'
import {USER, ADMIN} from '../../../auth'
import {Submit, orders, formItemLayout, tailFormItemLayout} from '../../../components/form'

const FormItem = Form.Item
const Option = Select.Option

class Widget extends Component {
  state = {
    required: true
  }
  componentDidMount() {
    const {setFieldsValue} = this.props.form
    const {id} = this.props.match.params
    if (id) {
      get(`/survey/fields/${id}`).then((rst) => {
        setFieldsValue({
          text: rst.text,
          body: JSON.parse(rst.body).join('\n'),
          sortOrder: rst.sortOrder,
          value: JSON.parse(rst.value).join(';'),
          name: rst.name,
          label: rst.label,
          type: rst.type
        })
        this.setState({required: rst.required})
      }).catch(message.error)
    } else {
      setFieldsValue({sortOrder: '0', type: 'text'})
    }
  }
  handleRequired = (e) => {
    this.setState({required: e.target.checked})
  }
  handleSubmit = (e) => {
    const {formatMessage} = this.props.intl
    const {push} = this.props
    const {id, formId} = this.props.match.params
    e.preventDefault();
    this.props.form.validateFields((err, values) => {
      if (!err) {
        post(
          id
          ? `/survey/fields/${id}`
          : `/survey/fields?formId=${formId}`,
        Object.assign({}, values, {
          sortOrder: parseInt(values.sortOrder, 10),
          required: this.state.required
        })).then(() => {
          message.success(formatMessage({id: "helpers.success"}))
          push(`/survey/forms`)
        }).catch(message.error);
      }
    });
  }
  render() {
    const {formatMessage} = this.props.intl
    const {getFieldDecorator} = this.props.form
    const {id, formId} = this.props.match.params
    const title = id
      ? {
        id: "buttons.edit",
        values: {
          id: id
        }
      }
      : {
        id: "buttons.new"
      }
    return (<Layout breads={[
        {
          href: '/survey/forms',
          label: {
            id: 'survey.forms.index.title'
          }
        }, {
          href: `/survey/fields/${formId}`,
          label: {
            id: "survey.fields.index.title"
          }
        },
        id
          ? {
            href: `/survey/fields/edit/${id}`,
            label: title
          }
          : {
            href: `/survey/fields/new/${formId}`,
            label: title
          }
      ]} title={title} roles={[USER, ADMIN]}>
      <Row>
        <Col md={{
            span: 8,
            offset: 2
          }}>
          <Form onSubmit={this.handleSubmit}>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "attributes.type" />}>
              {
                getFieldDecorator('type')(<Select>
                  {["text", "textarea", "select", "radios", "checkboxes"].map((p) => (<Option key={p} value={p}>{p}</Option>))}
                </Select>)
              }
            </FormItem>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "attributes.name" />} hasFeedback={true}>
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
            <FormItem {...formItemLayout} label={<FormattedMessage id = "attributes.label" />} hasFeedback={true}>
              {
                getFieldDecorator('label', {
                  rules: [
                    {
                      required: true,
                      message: formatMessage({id: "errors.empty"})
                    }
                  ]
                })(<Input/>)
              }
            </FormItem>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "attributes.value" />} hasFeedback={true}>
              {getFieldDecorator('value', {rules: []})(<Input/>)}
            </FormItem>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "attributes.body" />} hasFeedback={true}>
              {getFieldDecorator('body', {rules: []})(<Input.TextArea rows={8}/>)}
            </FormItem>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "attributes.sortOrder" />}>
              {
                getFieldDecorator('sortOrder')(<Select>
                  {orders(10).map((p) => (<Option key={p} value={p}>{p}</Option>))}
                </Select>)
              }
            </FormItem>
            <FormItem {...tailFormItemLayout}>
              <Checkbox checked={this.state.required} onChange={this.handleRequired}>
                <FormattedMessage id="attributes.required"/>
              </Checkbox>
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
