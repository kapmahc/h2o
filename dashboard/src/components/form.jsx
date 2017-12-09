import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {Form, Button, message} from 'antd'
import {FormattedMessage} from 'react-intl'
import ReactQuill from 'react-quill'

import {post} from '../ajax'

const FormItem = Form.Item

export const formItemLayout = {
  labelCol: {
    xs: {
      span: 24
    },
    sm: {
      span: 8
    }
  },
  wrapperCol: {
    xs: {
      span: 24
    },
    sm: {
      span: 16
    }
  }
};

export const tailFormItemLayout = {
  wrapperCol: {
    xs: {
      span: 24,
      offset: 0
    },
    sm: {
      span: 16,
      offset: 8
    }
  }
};

export class UEditor extends Component {
  onClick = () => {
    const {action, target} = this.props
    post('/api/token', {
      act: action,
      tid: target
    }).then((rst) => window.open(`${action}/${rst.token}`, '_blank')).catch(message.error)
  }
  render() {
    return (<Button onClick={this.onClick} shape="circle" icon="chrome"/>)
  }
}

UEditor.propTypes = {
  action: PropTypes.string.isRequired,
  target: PropTypes.number.isRequired
}

export const orders = (size) => Array(size * 2 + 1).fill().map((_, id) => (id - size).toString())

export class Quill extends Component {
  render() {
    const {value, onChange} = this.props
    const modules = {
      toolbar: [
        [
          {
            'font': []
          }
        ],
        [
          {
            size: []
          }
        ],
        [
          'bold', 'italic', 'underline', 'strike', 'blockquote'
        ],
        [
          {
            'list': 'ordered'
          }, {
            'list': 'bullet'
          }, {
            'indent': '-1'
          }, {
            'indent': '+1'
          }
        ],
        [
          'link', 'video'
        ],
        ['clean']
      ],
      clipboard: {
        // toggle to add extra line breaks when pasting HTML:
        matchVisual: false
      }
    }

    const formats = [
      'header',
      'font',
      'size',
      'bold',
      'italic',
      'underline',
      'strike',
      'blockquote',
      'list',
      'bullet',
      'indent',
      'link',
      'image',
      'video'
    ]
    return (<ReactQuill modules={modules} formats={formats} value={value} onChange={onChange} theme="snow"/>)
  }
}

Quill.propTypes = {
  value: PropTypes.string.isRequired,
  onChange: PropTypes.func.isRequired
}

export class Submit extends Component {
  render() {
    return (<FormItem {...tailFormItemLayout}>
      <Button type="primary" htmlType="submit">
        <FormattedMessage id="buttons.submit"/>
      </Button>
    </FormItem>);
  }
}
