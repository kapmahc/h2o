import React from 'react'
import {Route} from 'react-router'
import Exception from 'ant-design-pro/lib/Exception'

import Home from './Home'

import UsersSignIn from './users/SignIn'
import UsersSignUp from './users/SignUp'
import UsersEmailForm from './users/EmailForm'
import UsersResetPassword from './users/ResetPassword'
import UsersLogs from './users/Logs'
import UsersProfile from './users/Profile'
import UsersChangePassword from './users/ChangePassword'

import NewLeaveWord from './leave-words/New'

import AdminSiteStatus from './admin/site/Status'
import AdminSiteInfo from './admin/site/Info'
import AdminSiteAuthor from './admin/site/Author'
import AdminSiteSeo from './admin/site/Seo'
import AdminSiteSmtp from './admin/site/Smtp'
import AdminSiteHome from './admin/site/Home'
import AdminIndexUsers from './admin/users/Index'
import AdminIndexLeaveWords from './admin/leave-words/Index'
import AdminIndexLocales from './admin/locales/Index'
import AdminFormLocale from './admin/locales/Form'
import AdminIndexLinks from './admin/links/Index'
import AdminFormLink from './admin/links/Form'
import AdminIndexCards from './admin/cards/Index'
import AdminFormCard from './admin/cards/Form'
import AdminIndexFriendLinks from './admin/friend-links/Index'
import AdminFormFriendLink from './admin/friend-links/Form'

import AttachmentsIndex from './attachments/Index'

const UsersConfirm = () => (<UsersEmailForm action="confirm"/>)
const UsersUnlock = () => (<UsersEmailForm action="unlock"/>)
const UsersForgotPassword = () => (<UsersEmailForm action="forgot-password"/>)
const NotFound = () => (<Exception type="404"/>)

const routes = [
  (< Route key = "nut.home" exact path = "/" component = {
    Home
  } />),

  (< Route key = "nut.users.sign-in" path = "/users/sign-in" component = {
    UsersSignIn
  } />),
  (< Route key = "nut.users.sign-up" path = "/users/sign-up" component = {
    UsersSignUp
  } />),
  (< Route key = "nut.users.confirm" path = "/users/confirm" component = {
    UsersConfirm
  } />),
  (< Route key = "nut.users.unlock" path = "/users/unlock" component = {
    UsersUnlock
  } />),
  (< Route key = "nut.users.forgot-password" path = "/users/forgot-password" component = {
    UsersForgotPassword
  } />),
  (< Route key = "nut.users.reset-password" path = "/users/reset-password/:token" component = {
    UsersResetPassword
  } />),
  (< Route key = "nut.users.logs" path = "/users/logs" component = {
    UsersLogs
  } />),
  (< Route key = "nut.users.profile" path = "/users/profile" component = {
    UsersProfile
  } />),
  (< Route key = "nut.users.change-password" path = "/users/change-password" component = {
    UsersChangePassword
  } />),

  (< Route key = "nut.leave-words.new" path = "/leave-words/new" component = {
    NewLeaveWord
  } />),

  (< Route key = "nut.admin.site.status" path = "/admin/site/status" component = {
    AdminSiteStatus
  } />),
  (< Route key = "nut.admin.site.info" path = "/admin/site/info" component = {
    AdminSiteInfo
  } />),
  (< Route key = "nut.admin.site.author" path = "/admin/site/author" component = {
    AdminSiteAuthor
  } />),
  (< Route key = "nut.admin.site.seo" path = "/admin/site/seo" component = {
    AdminSiteSeo
  } />),
  (< Route key = "nut.admin.site.smtp" path = "/admin/site/smtp" component = {
    AdminSiteSmtp
  } />),
  (< Route key = "nut.admin.site.home" path = "/admin/site/home" component = {
    AdminSiteHome
  } />),

  (< Route key = "nut.admin.users.index" path = "/admin/users" component = {
    AdminIndexUsers
  } />),
  (< Route key = "nut.admin.leave-words.index" path = "/admin/leave-words" component = {
    AdminIndexLeaveWords
  } />),

  (< Route key = "nut.admin.locales.edit" path = "/admin/locales/edit/:code" component = {
    AdminFormLocale
  } />),
  (< Route key = "nut.admin.locales.new" path = "/admin/locales/new" component = {
    AdminFormLocale
  } />),
  (< Route key = "nut.admin.locales.index" path = "/admin/locales" component = {
    AdminIndexLocales
  } />),

  (< Route key = "nut.admin.friend-links.edit" path = "/admin/friend-links/edit/:id" component = {
    AdminFormFriendLink
  } />),
  (< Route key = "nut.admin.friend-links.new" path = "/admin/friend-links/new" component = {
    AdminFormFriendLink
  } />),
  (< Route key = "nut.admin.friend-links.index" path = "/admin/friend-links" component = {
    AdminIndexFriendLinks
  } />),

  (< Route key = "nut.admin.links.edit" path = "/admin/links/edit/:id" component = {
    AdminFormLink
  } />),
  (< Route key = "nut.admin.links.new" path = "/admin/links/new" component = {
    AdminFormLink
  } />),
  (< Route key = "nut.admin.links.index" path = "/admin/links" component = {
    AdminIndexLinks
  } />),

  (< Route key = "nut.admin.cards.edit" path = "/admin/cards/edit/:id" component = {
    AdminFormCard
  } />),
  (< Route key = "nut.admin.cards.new" path = "/admin/cards/new" component = {
    AdminFormCard
  } />),
  (< Route key = "nut.admin.cards.index" path = "/admin/cards" component = {
    AdminIndexCards
  } />),

  (< Route key = "nut.attachments.index" path = "/attachments" component = {
    AttachmentsIndex
  } />),

  (<Route key="nut.no-match" component={NotFound}/>)
]

export default routes
