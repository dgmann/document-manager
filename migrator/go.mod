module github.com/dgmann/document-manager/migrator

go 1.14

replace github.com/unidoc/unidoc => github.com/dgmann/unidoc v2.2.0+incompatible

require (
	github.com/denisenkom/go-mssqldb v0.0.0-20200206145737-bbfc9a55622e
	github.com/dgmann/document-manager/apiclient v0.0.0-20201220102232-6598e40d4778
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/jmoiron/sqlx v1.2.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/namsral/flag v1.7.4-pre
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/unidoc/unidoc v2.2.0+incompatible
	golang.org/x/crypto v0.0.0-20201002170205-7f63de1d35b0 // indirect
	golang.org/x/image v0.0.0-20200119044424-58c23975cae1 // indirect
	golang.org/x/sync v0.0.0-20201008141435-b3e1573b7520
	golang.org/x/sys v0.0.0-20201009025420-dfb3f7c4e634 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)
