import Link from 'next/link'
import Head from 'next/head'

export default({
  children,
  title = 'This is the default title'
}) => (<div>
  <Head>
    <meta name='viewport' content='width=device-width, initial-scale=1'/>
    <meta charSet='utf-8'/>
    <link rel='stylesheet' href='/static/main.css'/>
    <title>{title}</title>
  </Head>
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