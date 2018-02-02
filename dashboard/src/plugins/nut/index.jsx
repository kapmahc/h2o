const AdminFormLocale = import ('./admin/locales/Form')
const AdminFormLink = import ('./admin/links/Form')
const AdminFormCard = import ('./admin/cards/Form')
const AdminFormFriendLink = import ('./admin/friend-links/Form')

export default {
  routes: [
    {
      path: "/users/sign-in",
      component: import ('./users/SignIn')
    }, {
      path: "/users/sign-up",
      component: import ('./users/SignUp')
    }, {
      path: "/users/confirm",
      component: import ('./users/Confirm')
    }, {
      path: "/users/unlock",
      component: import ('./users/Unlock')
    }, {
      path: "/users/forgot-password",
      component: import ('./users/ForgotPassword')
    }, {
      path: "/users/reset-password/:token",
      component: import ('./users/ResetPassword')
    }, {
      path: "/users/logs",
      component: import ('./users/Logs')
    }, {
      path: "/users/profile",
      component: import ('./users/Profile')
    }, {
      path: "/users/change-password",
      component: import ('./users/ChangePassword')
    }, {
      path: "/leave-words/new",
      component: import ('./leave-words/New')
    }, {
      path: "/admin/site/status",
      component: import ('./admin/site/Status')
    }, {
      path: "/admin/site/info",
      component: import ('./admin/site/Info')
    }, {
      path: "/admin/site/author",
      component: import ('./admin/site/Author')
    }, {
      path: "/admin/site/seo",
      component: import ('./admin/site/Seo')
    }, {
      path: "/admin/site/smtp",
      component: import ('./admin/site/Smtp')
    }, {
      path: "/admin/site/donate",
      component: import ('./admin/site/Donate')
    }, {
      path: "/admin/site/home",
      component: import ('./admin/site/Home')
    }, {
      path: "/admin/users",
      component: import ('./admin/users/Index')
    }, {
      path: "/admin/leave-words",
      component: import ('./admin/leave-words/Index')
    }, {
      path: "/admin/locales/edit/:id",
      component: AdminFormLocale
    }, {
      path: "/admin/locales/new",
      component: AdminFormLocale
    }, {
      path: "/admin/locales",
      component: import ('./admin/locales/Index')
    }, {
      path: "/admin/friend-links/edit/:id",
      component: AdminFormFriendLink
    }, {
      path: "/admin/friend-links/new",
      component: AdminFormFriendLink
    }, {
      path: "/admin/friend-links",
      component: import ('./admin/friend-links/Index')
    }, {
      path: "/admin/links/edit/:id",
      component: AdminFormLink
    }, {
      path: "/admin/links/new",
      component: AdminFormLink
    }, {
      path: "/admin/links",
      component: import ('./admin/links/Index')
    }, {
      path: "/admin/cards/edit/:id",
      component: AdminFormCard
    }, {
      path: "/admin/cards/new",
      component: AdminFormCard
    }, {
      path: "/admin/cards",
      component: import ('./admin/cards/Index')
    }, {
      path: "/attachments",
      component: import ('./attachments/Index')
    }
  ],
  menus: []
}
