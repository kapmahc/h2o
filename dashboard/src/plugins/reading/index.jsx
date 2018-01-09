import React from 'react'
import {Route} from 'react-router'

import IndexBooks from './books/Index'

// import IndexComments from './comments/Index'
// import FormComment from './comments/Form'

const routes = [
  (< Route key = "reading.books.index" path = "/reading/books" component = {
    IndexBooks
  } />)

  // (< Route key = "forum.comments.edit" path = "/forum/comments/edit/:id" component = {
  //   FormComment
  // } />),
  // (< Route key = "forum.comments.new" path = "/forum/comments/new/:articleId" component = {
  //   FormComment
  // } />),
  // (< Route key = "forum.comments.index" path = "/forum/comments" component = {
  //   IndexComments
  // } />)
]

export default routes
