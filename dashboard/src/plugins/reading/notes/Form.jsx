import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {Form, Row, Col, message} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'

import Layout from '../../../layouts/dashboard'
import {post, get} from '../../../ajax'
import {USER, ADMIN} from '../../../auth'
import {Submit, Quill, formItemLayout} from '../../../components/form'

const FormItem = Form.Item

class Widget extends Component {
  state = {
    body: ''
  }
  componentDidMount() {
    const {id} = this.props.match.params
    if (id) {
      get(`/reading/notes/${id}`).then((rst) => {
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
    const {id, bookId} = this.props.match.params
    e.preventDefault();
    this.props.form.validateFields((err, values) => {
      if (!err) {
        post(
          id
          ? `/reading/notes/${id}`
          : `/reading/notes?bookId=${bookId}`,
        Object.assign({}, values, {
          type: 'html',
          body: this.state.body
        })).then(() => {
          message.success(formatMessage({id: "helpers.success"}))
          push('/reading/notes')
        }).catch(message.error);
      }
    });
  }
  render() {
    const {id, bookId} = this.props.match.params
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
          href: '/reading/notes',
          label: {
            id: 'reading.notes.index.title'
          }
        },
        id
          ? {
            href: `/reading/notes/edit/${id}`,
            label: title
          }
          : {
            href: `/reading/notes/new/${bookId}`,
            label: title
          }
      ]} title={title} roles={[USER, ADMIN]}>
      <Row>
        <Col md={{
            span: 18
          }}>
          <Form onSubmit={this.handleSubmit}>
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
