## Current (main branch)

* Fixed updating linked app components.
* Fixed errors during app component update not being reported.
* New `lvl system sshconfig` command to add system names to your SSH config. This allows easier access via non-lvl commands such as `rsync`.
* Fixed `lvl organisation get` sometimes giving an error due to API schema changes.
* `lvl domain create` can now be used without licensee for certain actions like `none`.

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
