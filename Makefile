version := 0.2.0
#path := $$HOME/.terraform.d/plugins/github.com/murtll/regru/${version}/linux_amd64
#path := $$HOME/.terraform.d/plugins/github.com/murtll/regru/${version}/darwin_arm64
path := $$HOME/.terraform.d/plugins/registry.terraform.io/murtll/regru/${version}/darwin_arm64
#path := $$HOME/.terraform.d/plugins/darwin_arm64

build:
	mkdir -p ${path}
	go build -o ${path}/terraform-provider-regru_${version}

get-go-version:
	@grep ^go go.mod | awk '{ print $$2 }'
