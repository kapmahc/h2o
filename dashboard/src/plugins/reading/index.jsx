const FormNote = import ('./notes/Form')

export default {
  menus: [],
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
