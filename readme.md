# Compare two composer files and merge them

For proper lock files comparison, you shoul create composer_l.lock and composer_r.lock

Result composer file will be written in composer_o.lock

It will contain all packages references from _l and _r files, but more fresh versions will be used

For packages version comparison github.com/mcuadros/go-version package is used