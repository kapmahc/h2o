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
    const {id} = this.props.match.params
    if (id) {
      get(`/api/forum/tags/${id}`).then((rst) => setFieldsValue({name: rst.name})).catch(message.error)
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
          ? `/api/forum/tags/${id}`
          : '/api/forum/tags',
        values).then(() => {
          message.success(formatMessage({id: "messages.success"}))
          push('/forum/tags')
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
          href: '/forum/tags',
          label: <FormattedMessage id='forum.tags.index.title'/>
        },
        id
          ? {
            href: `/forum/tags/edit/${id}`,
            label: (<FormattedMessage id={"buttons.edit"} values={{
                id: id
              }}/>)
          }
          : {
            href: "/forum/tags/new",
            label: <FormattedMessage id={"buttons.new"}/>
          }
      ]}>
      <Row>
        <Col md={{
            span: 8,
            offset: 2
          }}>
          <Form onSubmit={this.handleSubmit}>
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
