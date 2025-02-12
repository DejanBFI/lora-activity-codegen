package generator

import (
	"github.com/dave/jennifer/jen"
)

const (
	document  = "lora-partnership-ndf/internal/process/document"
	runtime   = "github.com/bfi-finance/lora-process-sdk/framework/runtime"
	framework = "github.com/bfi-finance/lora-process-sdk/framework"
	common    = "github.com/bfi-finance/lora-process-sdk/framework/defs/common"
	mapping   = "github.com/bfi-finance/lora-process-sdk/framework/defs/mapping"
)

var (
	importNames = map[string]string{
		document:  "document",
		runtime:   "runtime",
		framework: "framework",
		common:    "common",
		mapping:   "mapping",
	}
)

var (
	readSetVar = func() *jen.Statement {
		return jen.Id("readSet")
	}
	writeSetVar = func() *jen.Statement {
		return jen.Id("writeSet")
	}
	convVar = func() *jen.Statement {
		return jen.Id("conv")
	}
	convInVar = func() *jen.Statement {
		return jen.Id("convIn")
	}
	convOutVar = func() *jen.Statement {
		return jen.Id("convOut")
	}
	schemaVar = func() *jen.Statement {
		return jen.Id("schema")
	}
	processAndActivityNameVar = func() *jen.Statement {
		return jen.Id("processAndActivityName")
	}
)

func GenerateBoilerplate(out *jen.File, schemaName, processAndActivityName string, readSet, writeSet []string) {
	out.ImportNames(importNames)

	out.Const().Defs(
		schemaVar().Op("=").Lit("json-schema/proxy/inner/"+schemaName+".schema.json"),
		processAndActivityNameVar().Op("=").Lit(processAndActivityName),
	)

	generateReadSetCode(out, readSet)
	generateWriteSetCode(out, writeSet)
	generateConstructorLogic(out)
}

func generateReadSetCode(out *jen.File, readSet []string) {
	out.Comment("Read set.")
	readSetValues := make([]jen.Code, len(readSet))
	for _, v := range readSet {
		readSetValues = append(readSetValues, jen.Qual(document, v))
	}
	out.Var().Add(readSetVar()).Op("=").Index().Qual(common, "HString").Values(readSetValues...)

	readEnums := make([]jen.Code, len(readSet))
	for i, v := range readSet {
		e := jen.Id("readSet" + v)
		if i == 0 {
			e = e.Id("int").
				Op("=").
				Add(jen.Iota())
		}
		readEnums = append(readEnums, e)
	}
	out.Const().Defs(readEnums...)
}

func generateWriteSetCode(out *jen.File, writeSet []string) {
	out.Comment("Write set.")
	writeSetValues := make([]jen.Code, len(writeSet))
	for _, v := range writeSet {
		writeSetValues = append(writeSetValues, jen.Qual(document, v))
	}
	out.Var().Add(writeSetVar()).Op("=").Index().Qual(common, "HString").Values(writeSetValues...)

	writeEnums := make([]jen.Code, len(writeSet))
	for i, v := range writeSet {
		e := jen.Id("writeSet" + v)
		if i == 0 {
			e = e.Id("int").
				Op("=").
				Add(jen.Iota())
		}
		writeEnums = append(writeEnums, e)
	}
	out.Const().Defs(writeEnums...)
}

func generateConstructorLogic(out *jen.File) {
	out.Type().Id("Constructor").Struct(
		jen.Id("f").Op("*").Qual(runtime, "Function").Add(jen.Types(jen.Id("any"), jen.Id("any"))),
	)

	out.Func().Params(
		jen.Id("c").Op("*").Id("Constructor"),
	).Id("GenerateFunction").Params(
		jen.Id("apiLoader").
			Func().
			Params(jen.Id("url").String()).
			Params(
				jen.Op("*").Qual(framework, "APIFunction"),
				jen.Error(),
			),
		jen.Id("docFieldCheck").Func().Params(jen.Index().Qual(common, "HString")),
	).Error().Block(
		jen.Comment("forces a check at startup between the runtime document schema we parse and the fields used/generated in code"),
		jen.Id("docFieldCheck").Call(readSetVar()),
		jen.Id("docFieldCheck").Call(writeSetVar()),
		jen.Line(),

		jen.Comment("loads and parsed the API schema that we need"),
		jen.List(jen.Id("api"), jen.Err()).Op(":=").Id("apiLoader").Call(schemaVar()),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Return(jen.Err()),
		),
		jen.Line(),

		convInVar().
			Op(":=").
			Qual(mapping, "NewSimpleInputConverter").
			Types(jen.Map(jen.Qual(common, "HString")).Any()).
			Call(),
		jen.Id("fIn").Op(":=").Func().Params(jen.Id("m").Map(jen.Qual(common, "HString")).Any()).Params(jen.Op("*").Map(jen.Qual(common, "HString")).Any(), jen.Error()).Block(
			jen.Id("ret").Op(":=").Make(jen.Map(jen.Qual(common, "HString")).Any()),
			jen.Comment("TODO: implement the conversion logic here"),
			jen.Return(jen.Op("&").Id("ret"), jen.Nil()),
		),
		convInVar().Add(jen.Id(".SetInput")).Call(jen.Qual(common, "MakeReadSet").Call(readSetVar()), jen.Id("fIn")),

		jen.Line(),

		convOutVar().
			Op(":=").
			Qual(mapping, "NewSimpleOutputConverter").
			Types(jen.Map(jen.Qual(common, "HString")).Any()).
			Call(),
		jen.Id("fOut").Op(":=").Func().Params(jen.Id("data").Op("*").Map(jen.Qual(common, "HString")).Any()).Params(jen.Map(jen.Qual(common, "HString")).Any(), jen.Error()).Block(
			jen.Id("ret").Op(":=").Make(jen.Map(jen.Qual(common, "HString")).Any()),
			jen.Comment("TODO: implement the conversion logic here"),
			jen.Return(jen.Id("ret"), jen.Nil()),
		),
		convOutVar().Add(jen.Id(".SetOutput")).Call(jen.Qual(common, "MakeWriteSet").Call(writeSetVar()), jen.Id("fOut")),
		jen.Line(),

		convVar().Op(":=").Qual(mapping, "NewSimpleConverterFrom").Call(convInVar(), convOutVar()),
		jen.Id("c.f").Op("=").Id("api.ToRuntimeFunction").Call(processAndActivityNameVar(), jen.False(), convVar()),
		jen.Return(jen.Nil()),
	)

	out.Line()

	out.Func().Params(
		jen.Id("c").Op("*").Id("Constructor"),
	).Id("GenerateProcessStep").Params().Op("*").Qual(runtime, "ProcessStep").Block(
		jen.Id("p").Op(":=").Qual(runtime, "NewProcessStep").Call(processAndActivityNameVar(), jen.Id("c.f"), jen.Qual(runtime, "Eager"), jen.Index().Qual(runtime, "ProcessStepId").Values()),
		jen.Return(jen.Id("p")),
	)
}
