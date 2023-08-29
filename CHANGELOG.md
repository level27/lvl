## Current (main branch)

* `lvl app create`: if no organisation is given, your current organisation is used as default.
* Basic support for app component attachments.
* Added the ability to specify limit groups when adding components to Agency applications.

## 1.7.2
* Fixed `lvl app delete` not accepting `-y` to confirm deletion.

## 1.7.1
* Fixed `lvl update` not working on Linux/macOS.

## 1.7.0

* New `app component cron` commands for the new crons feature.
* `lvl login` now accepts username from command via `-u`, and can correctly read the password when piped via stdin. This allows automating it.
* `lvl update` allows you to easily update `lvl` to the latest version.
* Fix API marshalling errors with system checks and app component types.
* You can now refer to system checks by type in commands, e.g. `system check delete my.cool.system disk`
* You can now enable or disable DKIM on mail domains with `lvl mail domain dkim`.
* `lvl login` now supports `--trace` and passes the correct `User-Agent` to the API.
* Logging in with 2FA is now supported.

## 1.6.0

* `lvl system update` improvements:
    * Can now be used on imported systems (previously errored).
    * Fixed hostname getting cleared when used.
    * Can be used to set operating system info for imported systems.
* `lvl network zone add` command.
* `lvl organisation user sshkey create` command.
* `lvl domain zoneimport` command allows import DNS zone files into a domain.
* Fixed output of `lvl domain record create`.
* `lvl system sshconfig` stored IP addresses instead of FQDN.
* Added command help for `lvl sshkey favorite`.

## 1.5.1

* Invalid cookbook parameter names in `lvl system cookbooks update` and such now give an error again.

## 1.5.0

* Don't emit comments with `lvl system sshconfig` to avoid parsing issues on certain OpenSSH versions.
* Fix `lvl system describe` failing due to changes in API.
* New "id" output mode to print only IDs of values returned. This is for ease of use in scripts.

## 1.4.0

* Fixed updating linked app components.
* Fixed errors during app component update not being reported.
* New `lvl system sshconfig` command to add system names to your SSH config. This allows easier access via non-lvl commands such as `rsync`.
* Fixed `lvl organisation get` sometimes giving an error due to API schema changes.
* `lvl domain create` can now be used without licensee for certain actions like `none`.
* `lvl domain create` now accepts organisation names for `--organisation`, instead of solely IDs.
* Improved help for `lvl domain create` somewhat.

## 1.3.1

* Fixed the `sshkey` parameter on app components.

## 1.3.0

* Added `job delete <id>` command.

## 1.2.0

* Started tracking this changelog.
* Fix first run not being able to log in.
* `--version` flag and embed lvl version in the `User-Agent` headers sent.
* Allow waiting for commands like creates to complete (e.g. wait for status to change to "ok").
* Confirmations like "app created!" can now be shown as JSON for scripting purposes.
* System settings can now be managed as cookbooks.
* Can now create app components with required received parameters (e.g. Docker components).
* Can create linked app components.
* Added `job retry <id>` command.
