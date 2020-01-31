#!/bin/sh
function usage() {
  echo "docker run -e DOMAIN=domain -e TOKEN=slacktoken -e USER=user:password $0"
  exit 1
}
if [ -z "$DOMAIN" ]; then
  usage
fi
if [ -z "$TOKEN" ]; then
  usage
fi
if [ -z "$USER" ]; then
  usage
fi
mkdir /etc/supervisor
cat > /etc/supervisor/supervisord.conf <<EOF
[supervisord]
nodaemon=true

[program:rsyslog]
command=/usr/sbin/rsyslogd -n

[program:postfix]
command=/opt/postfix.sh
EOF

cat > /opt/postfix.sh <<EOF
#!/bin/sh
/usr/sbin/postfix start
trap '/usr/sbin/postfix stop' SIGTERM
tail -f /var/log/mail.log
EOF
chmod +x /opt/postfix.sh

postconf -e inet_interfaces=all
postconf -e myhostname=$DOMAIN
postconf -F '*/*/chroot = n'
postconf -e message_size_limit=102400000 
postconf -e mailbox_size_limit=102400000 

cat > /opt/config.toml <<EOF
valid_domain=["$DOMAIN"]
bot_token="$TOKEN"
log_dir="/var/log/"
EOF

postconf -e virtual_alias_maps=regexp:/etc/postfix/aliases_virtual
newaliases

cat > /etc/postfix/aliases_virtual <<EOF
/.*/ root
EOF

postconf -e mailbox_command=/opt/scan_mail

postconf -e smtpd_sasl_auth_enable=yes
postconf -e smtpd_sasl_path=smtpd
postconf -e broken_sasl_auth_clients=yes
postconf -e smtpd_recipient_restrictions=permit_sasl_authenticated,reject_unauth_destination
postconf -e cyrus_sasl_config_path=/etc/postfix/

cat > /etc/postfix/smtpd.conf <<EOF
pwcheck_method: auxprop
auxprop_plugin: sasldb
mech_list: PLAIN LOGIN CRAM-MD5 DIGEST-MD5 NTLM
EOF

function add_password() {
  while IFS=':' read -r _user _password; do
    echo $_password | saslpasswd2 -p -c -u $DOMAIN $_user
  done
  chown postfix -R /etc/sasl2
}

echo $USER | tr , \\n | add_password

exec /usr/bin/supervisord -c /etc/supervisor/supervisord.conf 
