keyfob is a two-factor authentication agent suitable for AWS and Github. Works pretty much the same as Google Authenticator, but uses your laptop's keychain.

## Installation

If you're on a mac, you can just do this:

    wget -O - https://raw.githubusercontent.com/StevenACoffman/keyfob/master/install.sh | bash


This will download the github 0.1.0 binary release for mac, and move any of your MFA secrets from `2fa` over to your keychain.

## Usage

    keyfob add [name] [key]
    keyfob otp [name]
    keyfob help

`keyfob add name` adds a new key to the keyfob keychain with the given name. It
prints a prompt to standard error and reads a two-factor key from standard
input. Two-factor keys are short case-insensitive strings of letters A-Z and
digits 2-7.

The new key generates time-based (TOTP) authentication codes.

`keyfob opt [name]` prints a One Time Password (aka two-factor authentication) code from the key with the
given name. If `--clip` is specified, `keyfob` also copies to the code to the system
clipboard.

The time-based authentication codes are derived from a hash of the
key and the current time, so it is important that the system clock have at
least one-minute accuracy.

The keychain is stored unencrypted in the text file `$HOME/.keyfob`.

## Example

During GitHub 2FA setup, at the “Scan this barcode with your app” step,
click the “enter this text code instead” link. A window pops up showing
“your two-factor secret,” a short string of letters and digits.

Add it to keyfob under the name github, typing the secret at the prompt:

    $ keyfob add github
    keyfob key for github: nzxxiidbebvwk6jb
    $

Then whenever GitHub prompts for a 2FA code, run keyfob to obtain one:

    $ keyfob otp github
    268346
    $

## Derivation

This is just a little toy cobbled together from [2fa](https://github.com/rsc/2fa/), [cobra](https://github.com/spf13/cobra), and [go-keyring](https://github.com/zalando/go-keyring) and using [goreleaser](https://github.com/goreleaser/goreleaser).

Unlike 2fa, this doesn't support listing all the stored codes, or adding 7 or 8 character long TOTP, or counter-based (HOTP) codes. Pillaging ... ehrm... adapting the 2fa code to do that in here would be easy, but I don't need it.

## Really, does this make sense?

At least to me, it does. My laptop features encrypted storage, a stronger authentication mechanism, and I take good care of its physical integrity.

My phone also runs arbitrary apps, is constantly connected to the Internet, gets forgotten on tables.

Thanks to the convenience of a command line utility, I'm more likely to enable MFA in more places.

Clearly a win for security.

## Dependencies

#### OS X

The OS X implementation depends on the `/usr/bin/security` binary for
interfacing with the OS X keychain. It should be available by default.

#### Linux

The Linux implementation depends on the [Secret Service][SecretService] dbus
interface, which is provided by [GNOME Keyring](https://wiki.gnome.org/Projects/GnomeKeyring).

It's expected that the default collection `login` exists in the keyring, because
it's the default in most distros. If it doesn't exist, you can create it through the
keyring frontend program [Seahorse](https://wiki.gnome.org/Apps/Seahorse):

 * Open `seahorse`
 * Go to **File > New > Password Keyring**
 * Click **Continue**
 * When asked for a name, use: **login**
 
 
## Usage with aws-vault

This assumes you have installed `keyfob` but need to set up your secrets.

Your own organization __*might*__ have a different preferred `source_profile` name from `sosourcey` below.

1. Skip to **[2](#2)** if you already added your AWS access key and secret access key to aws vault. Otherwise do this:
```
$ aws-vault add sosourcey --keychain login
```
2. <a name="2"></a>Go to AWS, and [make a new MFA token](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_mfa_enable_virtual.html#enable-virt-mfa-for-iam-user). Either take a screenshot of the QR Code (⌘⇧3 aka Command-Shift-3) and run `zbarimg` on it as below, or click the option to see the text version. Save that secret somewhere. Also add it to your Google Authenticator as normal.
```
brew cask install aws-vault
brew install go zbar awscli
# To get the text secret out of the QR Code if you didn't ask to see that
zbarimg AWS_IAM_Management_Console.png
```
3. Copy the `aws-credential-helper.sh` script in this repository to a place in your shell path and remember the absolute path to there. 

4. Add to your `.aws/config` file something like this:
```
[default]
credential_process = /Users/scoffman/bin/aws-credential-helper-engineer.sh
region = us-east-1
output = json
 
[profile sosourcey]
region = us-east-1
mfa_serial = arn:aws:iam::111111111111:mfa/scoffman
 
[profile engineer]
mfa_serial = arn:aws:iam::111111111111:mfa/scoffman
region = us-east-1
role_arn = arn:aws:iam::111111111111:role/put-power-role-here
source_profile = sosourcey
```
5. Ma

