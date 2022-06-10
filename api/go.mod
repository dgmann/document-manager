module github.com/dgmann/document-manager/api

require (
	github.com/cskr/pubsub v1.0.2
	github.com/dgmann/document-manager/pdf-processor/pkg/processor v0.0.0-20201219172655-8592d79ff120
	github.com/go-chi/chi v4.1.1+incompatible
	github.com/go-chi/cors v1.1.1
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/klauspost/compress v1.10.4 // indirect
	github.com/sirupsen/logrus v1.5.0
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/stretchr/testify v1.6.1
	go.mongodb.org/mongo-driver v1.5.1
	golang.org/x/crypto v0.0.0-20200414173820-0848c9571904 // indirect
	golang.org/x/net v0.0.0-20201216054612-986b41b23924
	golang.org/x/sync v0.0.0-20200317015054-43a5402ce75a // indirect
	google.golang.org/grpc v1.34.0
	gopkg.in/olahol/melody.v1 v1.0.0-20170518105555-d52139073376
)

go 1.13
