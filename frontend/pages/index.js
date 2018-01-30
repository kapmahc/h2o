import React from 'react';
import withRedux from 'next-redux-wrapper';
import Layout from '../layouts/application';

import makeStore from '../store';

class Widget extends React.Component {
  render() {
    return (<Layout>
      <div>Hello World.</div>
    </Layout>)
  }
};

export default withRedux(makeStore, (state) => ({}))(Widget);
