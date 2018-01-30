import React from 'react';
import {Button} from 'antd';

import Layout from '../layouts/application';

export default class Error extends React.Component {
  static getInitialProps({res, err}) {
    const statusCode = res
      ? res.statusCode
      : err
        ? err.statusCode
        : null;
    return {statusCode}
  }

  render() {
    return (<Layout>
      <Button type="primary">Primary{this.props.statusCode}</Button>
    </Layout>)
  }
}
