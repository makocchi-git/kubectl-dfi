module github.com/makocchi-git/kubectl-dfi

go 1.12

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/MakeNowJust/heredoc v0.0.0-20171113091838-e9091a26100e // indirect
	github.com/docker/docker v1.13.1 // indirect
	github.com/golangci/golangci-lint v1.17.1 // indirect
	github.com/gookit/color v1.1.7
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de // indirect
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/russross/blackfriday v2.0.0+incompatible // indirect
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/spf13/cobra v0.0.2
	github.com/spf13/pflag v1.0.1
	gotest.tools v2.2.0+incompatible // indirect
	k8s.io/api v0.0.0-20190531132109-d3f5f50bdd94
	k8s.io/apimachinery v0.0.0-20190531131812-859a0ba5e71a
	k8s.io/cli-runtime v0.0.0-20190531135611-d60f41fb4dc3
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/kubernetes v1.14.3
)

replace (
	github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.4.2
	github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190607191414-238f8eaa31aa
	github.com/russross/blackfriday => github.com/russross/blackfriday v1.5.2
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190531132438-d58e65e5f4b1
	k8s.io/kubernetes => k8s.io/kubernetes v1.14.2
)
