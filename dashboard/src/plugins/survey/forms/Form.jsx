import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {
  Form,
  Row,
  Col,
  Input,
  message,
  DatePicker
} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'
import moment from 'moment'

import Layout from '../../../layout'
import {post, get} from '../../../ajax'
import {Submit, Quill, formItemLayout, DATE_FORMAT} from '../../../components/form'

const FormItem = Form.Item
const {RangePicker} = DatePicker;

class Widget extends Component {
  state = {
    body: '',
    title: '',
    startUp: new Date().toISOString(),
    shutDown: new Date().toISOString()
  }
  componentDidMount() {
    const {setFieldsValue} = this.props.form
    const {id} = this.props.match.params
    if (id) {
      get(`/survey/forms/${id}`).then((rst) => {
        setFieldsValue({title: rst.title})
        this.setState({
          body: rst.body, startUp: rst.startUp.split('T')[0],
          shutDown: rst.shutDown.split('T')[0]
        })
      }).catch(message.error)
    }
  }
  handleBodyChange = (value) => {
    this.setState({body: value})
  }
  onChangeDateRange = (dt, ds) => {
    this.setState({startUp: ds[0], shutDown: ds[1]})
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
          ? `/survey/forms/${id}`
          : '/survey/forms',
        Object.assign({}, values, {
          type: 'html',
          mode: 'register',
          body: this.state.body,
          startUp: this.state.startUp,
          shutDown: this.state.shutDown
        })).then(() => {
          message.success(formatMessage({id: "helpers.success"}))
          push('/survey/forms')
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
          href: '/survey/forms',
          label: <FormattedMessage id='survey.forms.index.title'/>
        },
        id
          ? {
            href: `/survey/forms/edit/${id}`,
            label: (<FormattedMessage id={"buttons.edit"} values={{
                id: id
              }}/>)
          }
          : {
            href: "/survey/forms/new",
            label: <FormattedMessage id={"buttons.new"}/>
          }
      ]}>
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
            <FormItem {...formItemLayout} label={<FormattedMessage id = "attributes.dateRange" />}>
              <RangePicker value={[
                  moment(this.state.startUp, DATE_FORMAT),
                  moment(this.state.shutDown, DATE_FORMAT)
                ]} onChange={this.onChangeDateRange}/>
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
