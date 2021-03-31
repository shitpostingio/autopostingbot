package config

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/shitpostingio/autopostingbot/config/structs"
)

type TestingStruct struct {
	Integer int
	Float   float32
	String  string
	Struct  EmbeddedStruct

	Optional int `type:"optional"`
	Webhook  int `type:"webhook"` //same as optional in this context

	Reloadable int `type:"reloadable"`

	SliceOfPrimitives []int
	SliceOfStructs    []EmbeddedStruct

	SliceOptional []int `type:"optional"`
}

type EmbeddedStruct struct {
	Primitive int
	Slice     []int
}

func Test_checkMandatoryFields(t *testing.T) {
	type args struct {
		isReload bool
		config   structs.Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{

		{
			name: "without reload - success",
			args: args{
				isReload: false,
				config: structs.Config{
					Autoposting: structs.AutopostingConfiguration{
						BotToken:      "token",
						ChannelID:     -123232,
						MediaPath:     "/media",
						Algorithm:     "debug",
						ChannelHandle: "durov",
					},
					Tdlib: structs.TdlibConfiguration{
						APIId:              1223,
						APIHash:            "hash",
						DatabaseDirectory:  "/db",
						FilesDirectory:     "/files",
						SystemLanguageCode: "en",
						DeviceModel:        "test",
						SystemVersion:      "1.2",
						ApplicationVersion: "1.4",
					},
					DocumentStore: structs.DocumentStoreConfiguration{
						DatabaseName:   "database_name",
						AuthMechanism:  "auth",
						Username:       "username",
						Password:       "password",
						AuthSource:     "authsource",
						ReplicaSetName: "",
						Hosts:          nil,
					},
					AnalysisAPI: structs.AnalysisAPIConfiguration{
						Address:                  "http://example.com/api",
						ImageEndpoint:            "http://example.com/api",
						VideoEndpoint:            "http://example.com/api",
						AuthorizationHeaderName:  "Bearer",
						AuthorizationHeaderValue: "token",
						CallerAPIKeyHeaderName:   "name",
					},
					Localization: structs.LocalizationConfiguration{
						Path:     "/path/to/locale",
						Language: "en",
					},
				},
			},
		},
		{
			name: "without reload - failure",
			args: args{
				isReload: false,
				config: structs.Config{
					Autoposting: structs.AutopostingConfiguration{
						BotToken:      "token",
						ChannelID:     -123232,
						MediaPath:     "/media",
						Algorithm:     "debug",
						ChannelHandle: "durov",
					},
					Tdlib: structs.TdlibConfiguration{
						APIId:              1223,
						APIHash:            "hash",
						DatabaseDirectory:  "/db",
						FilesDirectory:     "/files",
						SystemLanguageCode: "en",
						DeviceModel:        "test",
						SystemVersion:      "1.2",
						ApplicationVersion: "1.4",
					},
					DocumentStore: structs.DocumentStoreConfiguration{
						DatabaseName:   "database_name",
						AuthMechanism:  "auth",
						Username:       "username",
						Password:       "password",
						AuthSource:     "authsource",
						ReplicaSetName: "",
						Hosts:          nil,
					},
					AnalysisAPI: structs.AnalysisAPIConfiguration{
						Address:                  "http://example.com/api",
						ImageEndpoint:            "http://example.com/api",
						VideoEndpoint:            "http://example.com/api",
						AuthorizationHeaderName:  "",
						AuthorizationHeaderValue: "token",
						CallerAPIKeyHeaderName:   "name",
					},
					Localization: structs.LocalizationConfiguration{
						Path:     "/path/to/locale",
						Language: "en",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkMandatoryFields(tt.args.isReload, tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("checkMandatoryFields() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkStruct(t *testing.T) {
	type args struct {
		isReload     bool
		parent       string
		typeToCheck  reflect.Type
		valueToCheck reflect.Value
	}
	tests := []struct {
		name       string
		args       args
		wantErrors []error
	}{
		{
			name: "all fields are present",
			args: args{
				isReload:    false,
				parent:      "test",
				typeToCheck: reflect.TypeOf(TestingStruct{}),
				valueToCheck: reflect.ValueOf(TestingStruct{
					Integer: 1,
					Float:   1.2,
					String:  "str",
					Struct: EmbeddedStruct{
						Primitive: 1,
						Slice:     []int{1},
					},
					Optional:          1,
					Webhook:           1,
					Reloadable:        1,
					SliceOfPrimitives: []int{1},
					SliceOfStructs: []EmbeddedStruct{
						{
							Primitive: 1,
							Slice:     []int{1},
						},
					},
					SliceOptional: []int{1},
				}),
			},
			wantErrors: []error{},
		},
		{
			name: "primitive missing",
			args: args{
				isReload:    false,
				parent:      "test",
				typeToCheck: reflect.TypeOf(TestingStruct{}),
				valueToCheck: reflect.ValueOf(TestingStruct{
					Integer: 1,
					Float:   1.2,
					String:  "",
					Struct: EmbeddedStruct{
						Primitive: 1,
						Slice:     []int{1},
					},
					Optional:          1,
					Webhook:           1,
					Reloadable:        1,
					SliceOfPrimitives: []int{1},
					SliceOfStructs: []EmbeddedStruct{
						{
							Primitive: 1,
							Slice:     []int{1},
						},
					},
				}),
			},
			wantErrors: []error{
				fmt.Errorf("non optional field String in section test had zero value"),
			},
		},
		{
			name: "few primitive missing",
			args: args{
				isReload:    false,
				parent:      "test",
				typeToCheck: reflect.TypeOf(TestingStruct{}),
				valueToCheck: reflect.ValueOf(TestingStruct{
					Integer: 0,
					Float:   1.2,
					String:  "",
					Struct: EmbeddedStruct{
						Primitive: 1,
						Slice:     []int{1},
					},
					Optional:          1,
					Webhook:           1,
					Reloadable:        1,
					SliceOfPrimitives: []int{1},
					SliceOfStructs: []EmbeddedStruct{
						{
							Primitive: 1,
							Slice:     []int{1},
						},
					},
				}),
			},
			wantErrors: []error{
				fmt.Errorf("non optional field Integer in section test had zero value"),
				fmt.Errorf("non optional field String in section test had zero value"),
			},
		},
		{
			name: "optional field missing",
			args: args{
				isReload:    false,
				parent:      "test",
				typeToCheck: reflect.TypeOf(TestingStruct{}),
				valueToCheck: reflect.ValueOf(TestingStruct{
					Integer: 1,
					Float:   1.2,
					String:  "str",
					Struct: EmbeddedStruct{
						Primitive: 1,
						Slice:     []int{1},
					},
					Webhook:           1,
					Reloadable:        1,
					SliceOfPrimitives: []int{1},
					SliceOfStructs: []EmbeddedStruct{
						{
							Primitive: 1,
							Slice:     []int{1},
						},
					},
				}),
			},
			wantErrors: []error{},
		},
		{
			name: "webhook field missing",
			args: args{
				isReload:    false,
				parent:      "test",
				typeToCheck: reflect.TypeOf(TestingStruct{}),
				valueToCheck: reflect.ValueOf(TestingStruct{
					Integer: 1,
					Float:   1.2,
					String:  "str",
					Struct: EmbeddedStruct{
						Primitive: 1,
						Slice:     []int{1},
					},
					Webhook:           0,
					Reloadable:        1,
					SliceOfPrimitives: []int{1},
					SliceOfStructs: []EmbeddedStruct{
						{
							Primitive: 1,
							Slice:     []int{1},
						},
					},
				}),
			},
			wantErrors: []error{},
		},
		{
			name: "field missing in embedded struct",
			args: args{
				isReload:    false,
				parent:      "test",
				typeToCheck: reflect.TypeOf(TestingStruct{}),
				valueToCheck: reflect.ValueOf(TestingStruct{
					Integer: 1,
					Float:   1.2,
					String:  "str",
					Struct: EmbeddedStruct{
						Slice: []int{1},
					},
					Webhook:           1,
					Reloadable:        1,
					SliceOfPrimitives: []int{1},
					SliceOfStructs: []EmbeddedStruct{
						{
							Primitive: 1,
							Slice:     []int{1},
						},
					},
				}),
			},
			wantErrors: []error{
				fmt.Errorf("non optional field Primitive in section Struct had zero value"),
			},
		},
		{
			name: "reloadable only",
			args: args{
				isReload:    true,
				parent:      "test",
				typeToCheck: reflect.TypeOf(TestingStruct{}),
				valueToCheck: reflect.ValueOf(TestingStruct{
					Reloadable: 1,
				}),
			},
			wantErrors: []error{},
		},
		{
			name: "empty slice in embedded",
			args: args{
				isReload:    false,
				parent:      "test",
				typeToCheck: reflect.TypeOf(TestingStruct{}),
				valueToCheck: reflect.ValueOf(TestingStruct{
					Integer: 1,
					Float:   1.2,
					String:  "str",
					Struct: EmbeddedStruct{
						Primitive: 1,
						Slice:     []int{},
					},
					Optional:          1,
					Webhook:           1,
					Reloadable:        1,
					SliceOfPrimitives: []int{1},
					SliceOfStructs: []EmbeddedStruct{
						{
							Primitive: 1,
							Slice:     []int{1},
						},
					},
				}),
			},
			wantErrors: []error{
				fmt.Errorf("non optional slice field Slice in Slice in Struct had zero length"),
			},
		},
		{
			name: "struct field missing in slice",
			args: args{
				isReload:    false,
				parent:      "test",
				typeToCheck: reflect.TypeOf(TestingStruct{}),
				valueToCheck: reflect.ValueOf(TestingStruct{
					Integer: 1,
					Float:   1.2,
					String:  "str",
					Struct: EmbeddedStruct{
						Primitive: 1,
						Slice:     []int{1},
					},
					Optional:          1,
					Webhook:           1,
					Reloadable:        1,
					SliceOfPrimitives: []int{1},
					SliceOfStructs: []EmbeddedStruct{
						{
							Slice: []int{1},
						},
					},
				}),
			},
			wantErrors: []error{
				fmt.Errorf("non optional field Primitive in section SliceOfStructs in test had zero value"),
			},
		},

		{
			name: "slice type field empty in one of the structs in slice",
			args: args{
				isReload:    false,
				parent:      "test",
				typeToCheck: reflect.TypeOf(TestingStruct{}),
				valueToCheck: reflect.ValueOf(TestingStruct{
					Integer: 1,
					Float:   1.2,
					String:  "str",
					Struct: EmbeddedStruct{
						Primitive: 1,
						Slice:     []int{1},
					},
					Optional:          1,
					Webhook:           1,
					Reloadable:        1,
					SliceOfPrimitives: []int{1},
					SliceOfStructs: []EmbeddedStruct{
						{
							Primitive: 1,
							Slice:     []int{},
						},
					},
				}),
			},
			wantErrors: []error{
				fmt.Errorf("non optional slice field Slice in Slice in SliceOfStructs in test had zero length"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErrors := checkStruct(tt.args.isReload, tt.args.parent, tt.args.typeToCheck, tt.args.valueToCheck); !reflect.DeepEqual(gotErrors, tt.wantErrors) {
				t.Errorf("checkStruct() = %v, want %v", gotErrors, tt.wantErrors)
			}
		})
	}
}
