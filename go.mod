module github.com/xitongsys/parquet-go-source

go 1.14

require (
	cloud.google.com/go/storage v1.15.0
	github.com/Azure/azure-pipeline-go v0.2.3
	github.com/Azure/azure-storage-blob-go v0.13.0
	github.com/aws/aws-sdk-go v1.38.45
	github.com/colinmarc/hdfs/v2 v2.2.0
	github.com/golang/mock v1.5.0
	github.com/ncw/swift v1.0.53
	github.com/spf13/afero v1.6.0
	github.com/xitongsys/parquet-go v1.6.0
)

replace github.com/Azure/azure-storage-blob-go => github.com/yangp18/azure-storage-blob-go v0.13.1-0.20210524163401-c2ea9947082d
