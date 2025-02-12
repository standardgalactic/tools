// Copyright 2019 Istio Authors
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package main

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
)

func main() {
	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generateFile(gen, f)
		}
		return nil
	})
}

func generateFile(gen *protogen.Plugin, file *protogen.File) {
	filename := file.GeneratedFilenamePrefix + "_deepcopy.gen.go"
	p := gen.NewGeneratedFile(filename, file.GoImportPath)

	protoIdent := protogen.GoIdent{
		GoName:       "Clone",
		GoImportPath: "github.com/golang/protobuf/proto",
	}
	p.P("// Code generated by protoc-gen-deepcopy. DO NOT EDIT.")
	p.P("package ", file.GoPackageName)
	var process func([]*protogen.Message)

	process = func(messages []*protogen.Message) {
		for _, message := range messages {
			// skip maps in protos.
			if message.Desc.Options().(*descriptorpb.MessageOptions).GetMapEntry() {
				continue
			}
			typeName := message.GoIdent.GoName
			// Generate DeepCopyInto() method for this type
			p.P(`// DeepCopyInto supports using `, typeName, ` within kubernetes types, where deepcopy-gen is used.`)
			p.P(`func (in *`, typeName, `) DeepCopyInto(out *`, typeName, `) {`)
			p.P(`p := `, protoIdent, `(in).(*`, typeName, `)`)
			p.P(`*out = *p`)
			p.P(`}`)

			// Generate DeepCopy() method for this type
			p.P(`// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new `, typeName, `. Required by controller-gen.`)
			p.P(`func (in *`, typeName, `) DeepCopy() *`, typeName, ` {`)
			p.P(`if in == nil { return nil }`)
			p.P(`out := new(`, typeName, `)`)
			p.P(`in.DeepCopyInto(out)`)
			p.P(`return out`)
			p.P(`}`)

			// Generate DeepCopyInterface() method for this type
			p.P(`// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new `, typeName, `. Required by controller-gen.`)
			p.P(`func (in *`, typeName, `) DeepCopyInterface() interface{} {`)
			p.P(`return in.DeepCopy()`)
			p.P(`}`)
			process(message.Messages)
		}
	}
	process(file.Messages)
}
