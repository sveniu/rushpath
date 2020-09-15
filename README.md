# Rushpath – a simple Dashlane CLI

Rushpath lets you manage multi-factor authentication for the Dashlane
password manager.

WARNING! The implementation is incomplete and may cause issues, particularly
around backup/setting synchronization. See the issue list for a better idea
of status.

## Features

- Enable and disable TOTP-based MFA, like Google Authenticator or Authy.

- Add or remove FIDO U2F devices, like Yubikeys.

## Installing

Installation is manual for now:

```sh
git clone https://github.com/sveniu/rushpath.git
cd rushpath/cmd/rushpath/
go build
./rushpath
```
