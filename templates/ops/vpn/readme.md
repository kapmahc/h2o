# OpenVPN System Based On User/Password Authentication with Day Control (shell script) - Ubuntu

- Install openvpn

  ```bash
  apt-get install openvpn
  ```

- Generate the certificate Copy the certificate generator scripts from OpenVPN docs

  ```bash
  cp -R /usr/share/doc/openvpn/examples/easy-rsa /etc/openvpn/.
  cd /etc/openvpn/easy-rsa/2.0/
  ```

- Modify certificate variables

  ```bash
  vi vars
  ```

  Edit this file and change the following lines into your case

  ```bash
  export KEY_COUNTRY="US"
  export KEY_PROVINCE="California"
  export KEY_CITY="Los Angeles"
  export KEY_ORG="{{.name}}"
  export KEY_EMAIL="{{.currentUser.Email}}"
  ```

- Save and exit. Run the variable script and clean

  ```bash
  source ./vars
  ./clean-all
  ```

- Generate the public and private certificates. Just press ENTER or YES by default

  ```bash
  ./build-ca
  ./build-key-server server
  ./build-key client
  ./build-dh
  mv keys /etc/openvpn/.
  ```

- Create directory for script '/etc/openvpn/script'

  ```bash
  mkdir /etc/openvpn/script
  cd /etc/openvpn/script
  ```

  Create file '/etc/openvpn/script/config.sh'

  ```bash
  #!/bin/bash
  HOST="{{.home}}/vpn/api"
  TOKEN="{{.token}}"
  ```

- Create file '/etc/openvpn/script/test_connect.sh'

  ```bash
  #!/bin/bash
  . /etc/openvpn/script/config.sh
  ##Test Authentication
  username=$1
  password=$2
  user_id=$()
  ##Check user
  [ "$user_id" != '' ] && [ "$user_id" = "$username" ] && echo "user : $username" && echo 'authentication ok.' && exit 0 || echo 'authentication failed.'; exit 1
  ```

- Create file '/etc/openvpn/script/login.sh'

  ```bash
  #!/bin/bash
  . /etc/openvpn/script/config.sh
  ##Authentication
  user_id=$()
  ##Check user
  [ "$user_id" != '' ] && [ "$user_id" = "$username" ] && echo "user : $username" && echo 'authentication ok.' && exit 0 || echo 'authentication failed.'; exit 1
  ```

  Create file '/etc/openvpn/script/connect.sh'

  ```bash
  #!/bin/bash
  . /etc/openvpn/script/config.sh
  ##insert data connection to table log
  echo '$common_name','$trusted_ip','$trusted_port','$ifconfig_pool_remote_ip','$remote_port_1','$bytes_received','$bytes_sent'
  ##set status online to user connected
  echo '$common_name'
  ```

- Create file '/etc/openvpn/script/disconnect.sh'

  ```bash
  #!/bin/bash
  . /etc/openvpn/script/config.sh
  ##set status offline to user disconnected
  echo '$common_name'
  ##insert data disconnected to table log
  echo '$bytes_received' '$bytes_sent' '$trusted_ip' '$trusted_port' '$common_name'
  ```

- Edit file /etc/sysctl.conf

  ```bash
  net.ipv4.ip_forward=1
  ```

- Iptables to share internet

  ```bash
  echo "1" > /proc/sys/net/ipv4/ip_forward
  echo "1" > /proc/sys/net/ipv4/ip_dynaddr

  iptables -A INPUT -i tun0 -j ACCEPT
  iptables -A FORWARD -i tun0 -j ACCEPT

  iptables -A INPUT -i tun1 -j ACCEPT
  iptables -A FORWARD -i tun1 -j ACCEPT

  iptables -A INPUT -p udp --dport {{.port}} -j ACCEPT

  iptables -t nat -A POSTROUTING -s {{.network}}.0/24 -o eth0 -j MASQUERADE
  ```

- Add to /etc/openvpn/server.conf

  ```
  # protocol port
  port {{.port}}
  proto udp
  dev tun

  # ip server client
  server {{.network}}.0 255.255.255.0

  # key
  ca /etc/openvpn/keys/ca.crt
  cert /etc/openvpn/keys/server.crt
  key /etc/openvpn/keys/server.key
  dh /etc/openvpn/keys/dh1024.pem

  # option
  persist-key
  persist-tun
  keepalive 5 60
  reneg-sec 432000

  # option authen.
  comp-lzo
  user nobody
  # group nogroup
  client-to-client
  username-as-common-name
  client-cert-not-required
  auth-user-pass-verify /etc/openvpn/script/login.sh via-env

  # push to client
  max-clients 50
  push "persist-key"
  push "persist-tun"
  push "redirect-gateway def1"
  push "explicit-exit-notify 1"

  # DNS-Server
  push "dhcp-option DNS 8.8.8.8"
  push "dhcp-option DNS 8.8.4.4"

  # script connect-disconnect
  script-security 3 system
  client-connect /etc/openvpn/script/connect.sh
  client-disconnect /etc/openvpn/script/disconnect.sh

  # log-status
  status /etc/openvpn/log/udp-{{.port}}.log
  log-append /etc/openvpn/log/openvpn.log
  verb 3
  ```

- Client config

  ```bash
  client
  dev tun

  proto udp
  remote x.x.x.x {{.port}}

  nobind
  auth-user-pass
  reneg-sec 432000
  resolv-retry infinite

  ca ca.crt
  comp-lzo
  verb 1
  ```
