// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.21.12
// source: fraud_detection/fraud_detection.proto

package fraud_detection

import (
	common "github.com/JuanGQCadavid/ds-practice-2025/utils/pb/common"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_fraud_detection_fraud_detection_proto protoreflect.FileDescriptor

var file_fraud_detection_fraud_detection_proto_rawDesc = string([]byte{
	0x0a, 0x25, 0x66, 0x72, 0x61, 0x75, 0x64, 0x5f, 0x64, 0x65, 0x74, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x2f, 0x66, 0x72, 0x61, 0x75, 0x64, 0x5f, 0x64, 0x65, 0x74, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x66, 0x72, 0x61, 0x75, 0x64, 0x1a, 0x13,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x32, 0xfe, 0x01, 0x0a, 0x15, 0x46, 0x72, 0x61, 0x75, 0x64, 0x44, 0x65, 0x74,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x36, 0x0a,
	0x09, 0x69, 0x6e, 0x69, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x13, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x49, 0x6e, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x14, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x49, 0x6e, 0x69, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36, 0x0a, 0x09, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x55, 0x73,
	0x65, 0x72, 0x12, 0x13, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4e, 0x65, 0x78, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x4e, 0x65, 0x78, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3c, 0x0a,
	0x0f, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x43, 0x72, 0x65, 0x64, 0x69, 0x74, 0x43, 0x61, 0x72, 0x64,
	0x12, 0x13, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4e, 0x65, 0x78, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4e,
	0x65, 0x78, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x0a, 0x63,
	0x6c, 0x65, 0x61, 0x6e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x13, 0x2e, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x4e, 0x65, 0x78, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4e, 0x65, 0x78, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x42, 0x44, 0x5a, 0x42, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x4a, 0x75, 0x61, 0x6e, 0x47, 0x51, 0x43, 0x61, 0x64, 0x61, 0x76, 0x69, 0x64,
	0x2f, 0x64, 0x73, 0x2d, 0x70, 0x72, 0x61, 0x63, 0x74, 0x69, 0x63, 0x65, 0x2d, 0x32, 0x30, 0x32,
	0x35, 0x2f, 0x75, 0x74, 0x69, 0x6c, 0x73, 0x2f, 0x70, 0x62, 0x2f, 0x66, 0x72, 0x61, 0x75, 0x64,
	0x5f, 0x64, 0x65, 0x74, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
})

var file_fraud_detection_fraud_detection_proto_goTypes = []any{
	(*common.InitRequest)(nil),  // 0: common.InitRequest
	(*common.NextRequest)(nil),  // 1: common.NextRequest
	(*common.InitResponse)(nil), // 2: common.InitResponse
	(*common.NextResponse)(nil), // 3: common.NextResponse
}
var file_fraud_detection_fraud_detection_proto_depIdxs = []int32{
	0, // 0: fraud.FraudDetectionService.initOrder:input_type -> common.InitRequest
	1, // 1: fraud.FraudDetectionService.checkUser:input_type -> common.NextRequest
	1, // 2: fraud.FraudDetectionService.checkCreditCard:input_type -> common.NextRequest
	1, // 3: fraud.FraudDetectionService.cleanOrder:input_type -> common.NextRequest
	2, // 4: fraud.FraudDetectionService.initOrder:output_type -> common.InitResponse
	3, // 5: fraud.FraudDetectionService.checkUser:output_type -> common.NextResponse
	3, // 6: fraud.FraudDetectionService.checkCreditCard:output_type -> common.NextResponse
	3, // 7: fraud.FraudDetectionService.cleanOrder:output_type -> common.NextResponse
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_fraud_detection_fraud_detection_proto_init() }
func file_fraud_detection_fraud_detection_proto_init() {
	if File_fraud_detection_fraud_detection_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_fraud_detection_fraud_detection_proto_rawDesc), len(file_fraud_detection_fraud_detection_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_fraud_detection_fraud_detection_proto_goTypes,
		DependencyIndexes: file_fraud_detection_fraud_detection_proto_depIdxs,
	}.Build()
	File_fraud_detection_fraud_detection_proto = out.File
	file_fraud_detection_fraud_detection_proto_goTypes = nil
	file_fraud_detection_fraud_detection_proto_depIdxs = nil
}
