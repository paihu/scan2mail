# mail's attachmentfile to Slack

Email's domain section and `const domain` must be equal.

Email's local section and Slack Display name must be equal.

## usage
* Put config file(config.toml) to binary exist dir.
* /path/to/binary < raw_email_data

## docker
* docker build -t scan_mail -f docker/Dockerfile .
* docker run --rm -d -e DOMAIN=<accept domain> -e TOKEN=<slack bot token> -e USER=<user:password[,user:password,...] for smtp auth> scan_mail
* set mailer smtpauth user@domain, password  and send attached email

