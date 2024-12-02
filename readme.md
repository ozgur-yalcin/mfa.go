# mfa

mfa *("multi-factor authentication")* is a command line tool for generating and validating one-time password.

**Description**:

* An easy-to-use substitute for 2FA apps like Google authenticator.
* Supports the OATH algorithms, such as TOTP and HOTP.
* No need for network connection.
* No need for phone.

## Installation

```
go mod tidy
go build -v
```

## Usage

```
mfa qr [flags] <image-path>
mfa gen [flags] <secret-key>
mfa add [flags] <issuer> <secret-key>
mfa set [flags] <issuer> <secret-key>
mfa del <issuer>
mfa list <issuer>
mfa version
```

```
Commands:
  gen         Generate otp by secret key
  qr          Create account by qr code
  add         Create account by secret key
  set         Update account by secret key
  del         Delete accounts by secret key
  list        List accounts
  version     show version
```

```
Flags:
  -m, --mode string    time-variant TOTP or event-based HOTP (default "totp")
  -H, --hash string    hash method (SHA1, SHA256, SHA512) (default "SHA1")
  -i, --period int     period of calculate otp for TOTP (default 30)
  -l, --digits int     otp length for HOTP (default 6)
  -c, --counter int    number of iterations count for HOTP
```

## Examples

### Generate code

Generate a **time-based** otp but do not save the secret key

```
mfa gen ADOO3MCCCVO5AVD6
```

Generate a **counter-based** otp with counter 1

```
mfa gen -m hotp -c 1 ADOO3MCCCVO5AVD6
```

### Create account

Create an issuerd GitHub

```
mfa add GitHub ADOO3MCCCVO5AVD6
```

Create an account, the issuer is GitHub, the user is ozgur-yalcin

```
mfa add GitHub:ozgur-yalcin ADOO3MCCCVO5AVD6
```

### List account

List all accounts

```shell
mfa list 
```

List all accounts named GitHub

```
mfa list GitHub
```

List accounts whose issuer is GitHub and whose user is ozgur-yalcin

```
mfa list GitHub:ozgur-yalcin
```

List accounts whose issuer is GitHub and whose user is ozgur-yalcin

```
mfa list GitHub ozgur-yalcin
```

### Delete accounts

Delete all accounts named GitHub

```
mfa del GitHub
```

Delete accounts  whose issuer is GitHub and whose user is ozgur-yalcin

```
mfa del GitHub ozgur-yalcin
```

### Update account

Update the secret key of accounts which issuer is GitHub

```
mfa set GitHub 5BRSSSBJUWBQBOXE
```

Update the secret key of accounts which issuer is GitHub and the user is ozgur-yalcin

```
mfa set GitHub:ozgur-yalcin 5BRSSSBJUWBQBOXE
```

## License

MIT License, see [license.md](license.md).
