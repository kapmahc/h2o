import {USER, ADMIN} from '../../auth'

const FormTag = import ('./tags/Form')
const FormArticle = import ('./articles/Form')
const FormComment = import ('./comments/Form')

export default {
  menus: [
    {
      icon: "tablet",
      label: "forum.dashboard.title",
      href: "forum",
      roles: [
        USER, ADMIN
      ],
      items: [
        {
          label: "forum.articles.index.title",
          href: "/forum/articles"
        }, {
          label: "forum.comments.index.title",
          href: "/forum/comments"
        }, {
          label: "forum.tags.index.title",
          href: "/forum/tags",
          roles: [ADMIN]
        }
      ]
    }
  ],
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
