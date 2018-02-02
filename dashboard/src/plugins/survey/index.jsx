const FormForm = import ('./forms/Form')
const FormField = import ('./fields/Form')

export default {
  menus: [],
  routes: [
    {
      path: "/survey/forms/edit/:id",
      component: FormForm
    }, {
      path: "/survey/forms/new",
      component: FormForm
    }, {
      path: "/survey/forms",
      component: import ('./forms/Index')
    }, {
      path: "/survey/fields/edit/:id",
      component: FormField
    }, {
      path: "/survey/fields/new/:formId",
      component: FormField
    }, {
      path: "/survey/fields/:formId",
      component: import ('./fields/Index')
    }, {
      path: "/survey/records/:formId",
      component: import ('./records/Index')
    }
  ]
}
