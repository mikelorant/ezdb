module github.com/mikelorant/ezdb2

go 1.18

require (
	github.com/AlecAivazis/survey/v2 v2.3.5
	github.com/alecthomas/chroma v0.10.0
	github.com/aws/aws-sdk-go-v2 v1.16.7
	github.com/aws/aws-sdk-go-v2/config v1.15.14
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.11.20
	github.com/aws/aws-sdk-go-v2/service/s3 v1.27.1
	github.com/go-sql-driver/mysql v1.6.0
	github.com/goccy/go-yaml v1.9.5
	github.com/jamf/go-mysqldump v0.7.1
	github.com/rodaine/table v1.0.2-0.20210416185537-a3154d83485f
	github.com/schollz/progressbar/v3 v3.8.6
	github.com/spf13/cobra v1.5.0
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d
	golang.org/x/term v0.0.0-20220526004731-065cf7ba2467
	gopkg.in/yaml.v3 v3.0.1
)

replace github.com/jamf/go-mysqldump => github.com/mikelorant/go-mysqldump v0.7.11

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.3 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.12.9 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.12.8 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.14 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.8 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.15 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.13.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.11.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.16.9 // indirect
	github.com/aws/smithy-go v1.12.0 // indirect
	github.com/dlclark/regexp2 v1.7.0 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/xerrors v0.0.0-20220609144429-65e65417b02f // indirect
)
