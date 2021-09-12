module github.com/crawlab-team/plugin-notification

go 1.15

replace (
	github.com/crawlab-team/crawlab-core => /Users/marvzhang/projects/crawlab-team/crawlab-core
	github.com/crawlab-team/crawlab-grpc => /Users/marvzhang/projects/crawlab-team/crawlab-grpc
	github.com/crawlab-team/crawlab-plugin => /Users/marvzhang/projects/crawlab-team/crawlab-plugin
	gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.1.0
)

require (
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible // indirect
	github.com/apex/log v1.9.0
	github.com/crawlab-team/crawlab-core v0.6.0-beta.20210802.1344
	github.com/crawlab-team/crawlab-db v0.1.1
	github.com/crawlab-team/crawlab-plugin v0.0.0-20210604093326-57f35f02daf5
	github.com/gavv/httpexpect/v2 v2.2.0
	github.com/gin-gonic/gin v1.6.3
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imroc/req v0.3.0
	github.com/jaytaylor/html2text v0.0.0-20200412013138-3577fbdbcff7 // indirect
	github.com/matcornic/hermes v1.2.0
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/spf13/viper v1.7.1
	github.com/ssor/bom v0.0.0-20170718123548-6386211fdfcf // indirect
	go.mongodb.org/mongo-driver v1.4.5
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/russross/blackfriday.v2 v2.0.0-00010101000000-000000000000 // indirect
)
