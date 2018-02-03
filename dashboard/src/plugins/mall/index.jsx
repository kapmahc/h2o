import {USER, ADMIN} from '../../auth'

export default {
  routes: [],
  menus: [
    {
      icon: "shopping-cart",
      label: "shop.dashboard.title",
      href: "shop",
      roles: [
        USER, ADMIN
      ],
      items: []
    }
  ]
}
