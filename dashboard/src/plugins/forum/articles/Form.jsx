import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {
  Form,
  Row,
  Col,
  Input,
  message,
  Checkbox
} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'
// import TagSelect from 'ant-design-pro/lib/TagSelect'
// <TagSelect checkedTags={[1, 3]} onChange={this.handleTagsChange} expandable={true}>
//   {this.state.tagOptions.map(it => (<TagSelect.Option key={it.id} value={it.id}>{it.name}</TagSelect.Option>))}
// </TagSelect>

import Layout from '../../../layouts/dashboard'
import {post, get} from '../../../ajax'
import {USER, ADMIN} from '../../../auth'
import {Submit, Quill, formItemLayout} from '../../../components/form'

const FormItem = Form.Item
const CheckboxGroup = Checkbox.Group

class Widget extends Component {
  state = {
    body: '',
    tagOptions: [],
    tagValues: []
  }
  componentDidMount() {
    const {setFieldsValue} = this.props.form
    const {id} = this.props.match.params
    if (id) {
      get(`/forum/articles/${id}`).then((rst) => {
        setFieldsValue({title: rst.title})
        this.setState({
          body: rst.body,
          tagValues: rst.tags
            ? rst.tags.map((t) => t.id)
            : []
        })
      }).catch(message.error)
    }
    get('/forum/tags').then((rst) => this.setState({
      // tagOptions: rst,
      tagOptions: rst.map((t) => {
        return {label: t.name, value: t.id}
      })
    })).catch(message.error)
  }
  handleBodyChange = (value) => {
    this.setState({body: value})
  }
  handleTagsChange = (values) => {
    this.setState({tagValues: values})
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
          ? `/forum/articles/${id}`
          : '/forum/articles',
        Object.assign({}, values, {
          type: 'html',
          body: this.state.body,
          tags: this.state.tagValues
        })).then(() => {
          message.success(formatMessage({id: "helpers.success"}))
          push('/forum/articles')
        }).catch(message.error);
      }
    });
  }
  render() {
    const {formatMessage} = this.props.intl
    const {getFieldDecorator} = this.props.form
    const {id} = this.props.match.params
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
          href: '/forum/articles',
          label: {
            id: 'forum.articles.index.title'
          }
        },
        id
          ? {
            href: `/forum/articles/edit/${id}`,
            label: title
          }
          : {
            href: "/forum/articles/new",
            label: title
          }
      ]} title={title} roles={[USER, ADMIN]}>
      <Row>
        <Col md={{
            span: 18
          }}>
          <Form onSubmit={this.handleSubmit}>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "attributes.title" />} hasFeedback={true}>
              {
                getFieldDecorator('title', {
                  rules: [
                    {
                      required: true,
                      message: formatMessage({id: "errors.empty"})
                    }
                  ]
                })(<Input/>)
              }
            </FormItem>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "forum.attributes.article.tags" />}>
              <CheckboxGroup options={this.state.tagOptions} value={this.state.tagValues} onChange={this.handleTagsChange}/>
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
