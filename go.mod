module github.com/protolambda/zcli

go 1.16

require (
	github.com/protolambda/ask v0.0.5
	github.com/protolambda/messagediff v1.4.0
	github.com/protolambda/zrnt v0.16.0
	github.com/protolambda/ztyp v0.1.8
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)

replace github.com/protolambda/zrnt => ../zrnt
replace github.com/protolambda/ask => ../ask
