import Link from 'next/link'
import Head from 'next/head'

import bootstrap_css from 'bootstrap/dist/css/bootstrap.min.css';
import quill_snow_css from 'react-quill/dist/quill.snow.css';

import Header from './Header';
import Footer from './Footer';

export default({
  children,
  title = 'This is the default title'
}) => (<div>
  <Head>
    <meta charSet='utf-8'/>
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no"/>
    <title>{title}</title>
  </Head>
  <style jsx="jsx" global="global">
    {
      bootstrap_css
    }</style>
  <style jsx="jsx" global="global">
    {
      quill_snow_css
    }</style>
  <style jsx="jsx" global="global">
    {
      "body { padding-top: 3rem; padding-bottom: 3rem; color: #5a5a5a; }"
    }</style>
  <Header/> {children}
  <Footer/>
</div>);
