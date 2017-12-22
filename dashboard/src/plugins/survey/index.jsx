import React from 'react'
import {Route} from 'react-router'

import IndexForms from './forms/Index'
import FormForm from './forms/Form'
import IndexFields from './fields/Index'
import FormField from './fields/Form'
import IndexRecords from './records/Index'

const routes = [
  (< Route key = "survey.forms.edit" path = "/survey/forms/edit/:id" component = {
    FormForm
  } />),
  (< Route key = "survey.forms.new" path = "/survey/forms/new" component = {
    FormForm
  } />),
  (< Route key = "survey.forms.index" path = "/survey/forms" component = {
    IndexForms
  } />),

  (< Route key = "survey.fields.edit" path = "/survey/fields/edit/:id" component = {
    FormField
  } />),
  (< Route key = "survey.fields.new" path = "/survey/fields/new/:formId" component = {
    FormField
  } />),
  (< Route key = "survey.fields.index" path = "/survey/fields/:formId" component = {
    IndexFields
  } />),

  (< Route key = "survey.fields.index" path = "/survey/records/:formId" component = {
    IndexRecords
  } />)
]

export default routes
