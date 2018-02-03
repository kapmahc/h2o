import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {Form, Row, Col, Input, message} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'

import Layout from '../../../layouts/dashboard'
import {post, get} from '../../../ajax'
import {ADMIN} from '../../../auth'
import {Submit, Quill, formItemLayout} from '../../../components/form'

const FormItem = Form.Item

class Widget extends Component {
  state = {
    title: '',
    body: '',
    comments: []
  }
  componentDidMount() {
    const {setFieldsValue} = this.props.form
    const {id} = this.props.match.params
    get(`/forum/articles/${id}`).then((rst) => {
      setFieldsValue({favicon: rst.favicon})
      if (rst.home) {
        this.setState({body: rst.home.body})
      }

    }).catch(message.error)
  }
  handleSubmit = (e) => {
    const {formatMessage} = this.props.intl
    e.preventDefault();
    this.props.form.validateFields((err, values) => {
      if (!err) {
        post('/admin/site/home', Object.assign({}, values, {
          type: 'html',
          body: this.state.body
        })).then(() => {
          message.success(formatMessage({id: "helpers.success"}))
        }).catch(message.error);
      }
    });
  }
  handleBodyChange = (value) => {
    this.setState({body: value})
  }
  render() {
    const {formatMessage} = this.props.intl
    const {getFieldDecorator} = this.props.form
    const title = {
      id: "nut.admin.site.home.title"
    }
    return (<Layout breads={[{
          href: "/admin/site/home",
          label: title
        }
      ]} title={title} roles={[ADMIN]}>
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
            <FormItem {...formItemLayout} label={<FormattedMessage id = "attributes.body" />}>
              <Quill value={this.state.body} onChange={this.handleBodyChange}/>
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
