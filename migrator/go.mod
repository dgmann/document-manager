module github.com/dgmann/document-manager/migrator

go 1.14

replace github.com/unidoc/unidoc => github.com/dgmann/unidoc v2.2.0+incompatible

require (
	github.com/denisenkom/go-mssqldb v0.0.0-20200206145737-bbfc9a55622e
	github.com/dgmann/document-manager/api v0.0.0-20200417172829-90ed1db6b821
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a // indirect
	github.com/jmoiron/sqlx v1.2.0
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/mongodb/mongo-go-driver v0.1.1-0.20181220233027-8051092034cf // indirect
	github.com/namsral/flag v1.7.4-pre
	github.com/onsi/gomega v1.4.2 // indirect
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.5.0
	github.com/stretchr/testify v1.5.1
	github.com/unidoc/unidoc v2.2.0+incompatible
	golang.org/x/image v0.0.0-20200119044424-58c23975cae1 // indirect
	golang.org/x/sync v0.0.0-20200317015054-43a5402ce75a
	google.golang.org/appengine v1.6.5 // indirect
	gopkg.in/airbrake/gobrake.v2 v2.0.9 // indirect
	gopkg.in/gemnasium/logrus-airbrake-hook.v2 v2.1.2 // indirect
)
