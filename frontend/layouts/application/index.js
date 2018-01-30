import Link from 'next/link'
import Head from 'next/head'

import bootstrap_css from 'bootstrap/dist/css/bootstrap.css';
import quill_snow_css from 'react-quill/dist/quill.snow.css';

export default({
  children,
  title = 'This is the default title'
}) => (<div>
  <Head>
    <meta name='viewport' content='width=device-width, initial-scale=1'/>
    <meta charSet='utf-8'/>
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
      quill_snow_css
    }</style>
  <header>
    <nav>
      <Link href='/'>
        <a>Home</a>
      </Link>
      |
      <Link href='/about'>
        <a>About</a>
      </Link>
      |
      <Link href='/contact'>
        <a>Contact</a>
      </Link>
    </nav>
  </header>

  {children}

  <footer>
    {'I`m here to stay'}
  </footer>
</div>)