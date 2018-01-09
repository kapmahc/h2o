import React from 'react'
import {Route} from 'react-router'

import IndexBooks from './books/Index'

import IndexNotes from './notes/Index'
import FormNote from './notes/Form'

const routes = [
  (< Route key = "reading.books.index" path = "/reading/books" component = {
    IndexBooks
  } />),

  (< Route key = "reading.notes.edit" path = "/reading/notes/edit/:id" component = {
    FormNote
  } />),
  (< Route key = "reading.notes.new" path = "/reading/notes/new/:bookId" component = {
    FormNote
  } />),
  (< Route key = "reading.notes.index" path = "/reading/notes" component = {
    IndexNotes
  } />)
]

export default routes
