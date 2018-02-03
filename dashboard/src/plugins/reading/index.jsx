import {USER, ADMIN} from '../../auth'

const FormNote = import ('./notes/Form')

export default {
  menus: [
    {
      icon: "book",
      label: "reading.dashboard.title",
      href: "reading",
      roles: [
        ADMIN, USER
      ],
      items: [
        {
          label: "reading.books.index.title",
          href: "/reading/books",
          roles: [ADMIN]
        }, {
          label: "reading.notes.index.title",
          href: "/reading/notes"
        }
      ]
    }
  ],
  routes: [
    {
      path: "/reading/books",
      component: import ('./books/Index')
    }, {
      path: "/reading/notes/edit/:id",
      component: FormNote
    }, {
      path: "/reading/notes/new/:bookId",
      component: FormNote
    }, {
      path: "/reading/notes",
      component: import ('./notes/Index')
    }
  ]
}
