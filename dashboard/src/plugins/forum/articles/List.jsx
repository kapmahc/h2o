import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {Row, Col, Card} from 'antd'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'
import {Link} from 'react-router-dom'
import {FormattedMessage} from 'react-intl'

import Layout from '../../../layouts/application'
import Quill from '../../../components/Quill'

class Widget extends Component {
  render() {
    const {title, items} = this.props
    return (<Layout breads={[{
          href: "/forum/articles/latest",
          label: title
        }
      ]} title={title}>
      <Row gutter={16}>
        {
          items.map((it, id) => (<Col md={{
              span: 6
            }} key={id}>
            <Card title={it.title} extra={<Link to = {
                `/forum/articles/show/${it.id}`
              } > <FormattedMessage id="buttons.more"/></Link>}>
              <Quill body={it.body}/>
            </Card>
          </Col>))
        }
      </Row>
    </Layout>);
  }
}

Widget.propTypes = {
  title: PropTypes.object.isRequired,
  items: PropTypes.array.isRequired,
  push: PropTypes.func.isRequired
}

export default connect(state => ({}), {
  push
},)(Widget)
