module github.com/dgmann/document-manager/migrator

replace github.com/unidoc/unidoc => github.com/dgmann/unidoc v2.2.0+incompatible

require (
	cloud.google.com/go v0.33.0 // indirect
	github.com/denisenkom/go-mssqldb v0.0.0-20181014144952-4e0d7dc8888f
	github.com/dgmann/document-manager/api v0.0.0-20181113164830-0eeea4ba9b30 // indirect
	github.com/dgmann/document-manager/api-client v0.0.0-20181113164830-0eeea4ba9b30
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8 // indirect
	github.com/google/go-cmp v0.2.0 // indirect
	github.com/gosuri/uilive v0.0.0-20170323041506-ac356e6e42cd // indirect
	github.com/gosuri/uiprogress v0.0.0-20170224063937-d0567a9d84a1
	github.com/jmoiron/sqlx v1.2.0
	github.com/mattn/go-isatty v0.0.4 // indirect
	github.com/mattn/go-sqlite3 v1.10.0 // indirect
	github.com/namsral/flag v1.7.4-pre
	github.com/pkg/errors v0.8.0
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.2.0
	github.com/stretchr/testify v1.2.2
	github.com/unidoc/unidoc v2.1.1+incompatible
	golang.org/x/crypto v0.0.0-20181112202954-3d3f9f413869 // indirect
	golang.org/x/image v0.0.0-20181109232246-249dc8530c0e // indirect
	golang.org/x/sys v0.0.0-20181107165924-66b7b1311ac8 // indirect
	google.golang.org/appengine v1.3.0 // indirect
)

go 1.13
