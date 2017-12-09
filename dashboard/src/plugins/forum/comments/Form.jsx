import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {Form, Row, Col, message} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'

import Layout from '../../../layout'
import {post, get} from '../../../ajax'
import {Submit, Quill, formItemLayout} from '../../../components/form'

const FormItem = Form.Item

class Widget extends Component {
  state = {
    body: ''
  }
  componentDidMount() {
    const {id} = this.props.match.params
    if (id) {
      get(`/api/forum/comments/${id}`).then((rst) => {
        this.setState({body: rst.body})
      }).catch(message.error)
    }
  }
  handleBodyChange = (value) => {
    this.setState({body: value})
  }
  handleSubmit = (e) => {
    const {formatMessage} = this.props.intl
    const {push} = this.props
    const {id, articleId} = this.props.match.params
    e.preventDefault();
    this.props.form.validateFields((err, values) => {
      if (!err) {
        post(
          id
          ? `/api/forum/comments/${id}`
          : `/api/forum/comments?articleId=${articleId}`,
        Object.assign({}, values, {
          type: 'html',
          body: this.state.body
        })).then(() => {
          message.success(formatMessage({id: "messages.success"}))
          push('/forum/comments')
        }).catch(message.error);
      }
    });
  }
  render() {
    const {id} = this.props.match.params
    return (<Layout breads={[
        {
          href: '/forum/comments',
          label: <FormattedMessage id='forum.comments.index.title'/>
        },
        id
          ? {
            href: `/forum/comments/edit/${id}`,
            label: (<FormattedMessage id={"buttons.edit"} values={{
                id: id
              }}/>)
          }
          : {
            href: "/forum/comments/new",
            label: <FormattedMessage id={"buttons.new"}/>
          }
      ]}>
      <Row>
        <Col md={{
            span: 18
          }}>
          <Form onSubmit={this.handleSubmit}>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "attributes.name" />}>
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
