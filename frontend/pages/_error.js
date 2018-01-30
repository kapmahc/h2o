import React from 'react';
import {Button} from 'antd';
import withRedux from 'next-redux-wrapper';

import Layout from '../layouts/application';
import makeStore from '../store';

class Widget extends React.Component {
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

export default withRedux(makeStore, (state) => ({}))(Widget);
