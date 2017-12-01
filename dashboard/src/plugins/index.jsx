import nut from './nut'
import forum from './forum'
import survey from './survey'
import reading from './reading'
import mall from './mall'
import pos from './pos'
import erp from './erp'
import ops_vpn from './ops/vpn'
import ops_mail from './ops/mail'

const routes = [].concat(forum).concat(survey).concat(reading).concat(mall).concat(pos).concat(erp).concat(ops_mail).concat(ops_vpn).concat(nut)

export default {routes}
