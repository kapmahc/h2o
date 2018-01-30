import Document, {Head, Main, NextScript} from 'next/document';
import flush from 'styled-jsx/server';
import fetch from 'isomorphic-unfetch';
import {connect} from 'react-redux';
import PropTypes from 'prop-types';

import {refresh, signIn} from '../actions';
import {backend, get_locale, set_locale} from '../utils';

export default class Widget extends Document {
  static async getInitialProps(context) {
    const props = await super.getInitialProps(context);
    const locale = get_locale();
    // const {
    //   req: {
    //     locale,
    //     localeDataScript
    //   }
    // } = context
    //
    // const res = await fetch(backend('/layout'));
    // const siteInfo = await res.json();
    // refresh(siteInfo);
    return {
      ...props,
      locale
      // locale,
      // localeDataScript
    }
  }
  // static async getInitialProps({renderPage}) {
  //   const {html, head, errorHtml, chunks} = renderPage()
  //   const styles = flush()
  //   const res = await fetch(backend('/layout'));
  //   const siteInfo = await res.json();
  //    refresh(siteInfo);
  //   console.log(siteInfo);
  //   return {
  //     html,
  //     head,
  //     errorHtml,
  //     chunks,
  //     styles,
  //     siteInfo
  //   };
  // }

  render() {
    const {locale} = this.props;
    return (<html lang={locale}>
      <Head></Head>
      <body>
        <Main/>
        <NextScript/>
      </body>
    </html>);
  }
};

// Widget.propTypes = {
//   refresh: PropTypes.func.isRequired,
//   signIn: PropTypes.func.isRequired,
//   user: PropTypes.object.isRequired,
//   site: PropTypes.object.isRequired
// }
//
// export default connect(state => ({user: state.currentUser, site: state.siteInfo}), {refresh, signIn})(Widget);