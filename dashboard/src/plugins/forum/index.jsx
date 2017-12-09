import React from 'react'
import {Route} from 'react-router'

import IndexTags from './tags/Index'
import FormTag from './tags/Form'
import IndexArticles from './articles/Index'
import FormArticle from './articles/Form'
import IndexComments from './comments/Index'
import FormComment from './comments/Form'

const routes = [
  (< Route key = "forum.tags.edit" path = "/forum/tags/edit/:id" component = {
    FormTag
  } />),
  (< Route key = "forum.tags.new" path = "/forum/tags/new" component = {
    FormTag
  } />),
  (< Route key = "forum.tags.index" path = "/forum/tags" component = {
    IndexTags
  } />),

  (< Route key = "forum.articles.edit" path = "/forum/articles/edit/:id" component = {
    FormArticle
  } />),
  (< Route key = "forum.articles.new" path = "/forum/articles/new" component = {
    FormArticle
  } />),
  (< Route key = "forum.articles.index" path = "/forum/articles" component = {
    IndexArticles
  } />),

  (< Route key = "forum.comments.edit" path = "/forum/comments/edit/:id" component = {
    FormComment
  } />),
  (< Route key = "forum.comments.new" path = "/forum/comments/new/:articleId" component = {
    FormComment
  } />),
  (< Route key = "forum.comments.index" path = "/forum/comments" component = {
    IndexComments
  } />)
]

export default routes
