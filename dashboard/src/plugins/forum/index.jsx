const FormTag = import ('./tags/Form')
const FormArticle = import ('./articles/Form')
const FormComment = import ('./comments/Form')

export default {
  menus: [],
  routes: [
    {
      path: "/forum/tags/edit/:id",
      component: FormTag
    }, {
      path: "/forum/tags/new",
      component: FormTag
    }, {
      path: "/forum/tags",
      component: import ('./tags/Index')
    }, {
      path: "/forum/articles/edit/:id",
      component: FormArticle
    }, {
      path: "/forum/articles/new",
      component: FormArticle
    }, {
      path: "/forum/articles",
      component: import ('./articles/Index')
    }, {
      path: "/forum/comments/edit/:id",
      component: FormComment
    }, {
      path: "/forum/comments/new/:articleId",
      component: FormComment
    }, {
      path: "/forum/comments",
      component: import ('./comments/Index')
    }
  ]
}
