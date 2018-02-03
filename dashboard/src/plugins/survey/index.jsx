import {USER, ADMIN} from '../../auth'

const FormForm = import ('./forms/Form')
const FormField = import ('./fields/Form')

export default {
  menus: [
    {
      icon: "notification",
      label: "survey.dashboard.title",
      href: "survey",
      roles: [
        ADMIN, USER
      ],
      items: [
        {
          label: "survey.forms.index.title",
          href: "/survey/forms"
        }
      ]
    }
  ],
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
