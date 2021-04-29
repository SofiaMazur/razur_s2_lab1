package gomodule

import (
	"fmt"
	"path"

	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
)

var (
	
	pctx = blueprint.NewPackageContext("github.com/FictProger/architecture2-lab-1/build/gomodule")

	
	goBuild = pctx.StaticRule("binaryBuild", blueprint.RuleParams{
		Command:     "cd $workDir && go build -o $outputPath $pkg",
		Description: "build go command $pkg",
	}, "workDir", "outputPath", "pkg")

	
	goVendor = pctx.StaticRule("vendor", blueprint.RuleParams{
		Command:     "cd $workDir && go mod vendor",
		Description: "vendor dependencies of $name",
	}, "workDir", "name")

	
	goTest = pctx.StaticRule("test", blueprint.RuleParams{
		Command:     "cd $workDir && go test -v $testPkg > $testOutput",
		Description: "testing $testPkg",
	}, "workDir", "testPkg", "testOutput")
)


type goBinaryModuleType struct {
	blueprint.SimpleName

	properties struct {
		Pkg string
		Srcs []string
		SrcsExclude []string
		VendorFirst bool
		TestPkg string

		Deps []string
	}
}

func (gb *goBinaryModuleType) DynamicDependencies(blueprint.DynamicDependerModuleContext) []string {
	return gb.properties.Deps
}

func isTest(fileName string) bool {
	endingIndex := len(fileName) - len("_test.go")
	if endingIndex < 1 {
		return false
	}
	return fileName[endingIndex:] == "_test.go"
}

func (gb *goBinaryModuleType) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	config := bood.ExtractConfig(ctx)
	config.Debug.Printf("Adding build actions for go binary module '%s'", name)

	outputPath := path.Join(config.BaseOutputDir, "bin", name)
	testOutputPath := path.Join(config.BaseOutputDir, "testOutput.txt")

	var inputs []string
	var testInputs []string
	inputErors := false

	for _, src := range gb.properties.Srcs {
		if matches, err := ctx.GlobWithDeps(src, gb.properties.SrcsExclude); err == nil {
			for _, file := range matches {
				if isTest(file) {
					testInputs = append(testInputs, file)
				} else {
					inputs = append(inputs, file)
				}
			}
			inputs = append(inputs, matches...)
		} else {
			ctx.PropertyErrorf("srcs", "Cannot resolve files that match pattern %s", src)
			inputErors = true
		}
	}
	if inputErors {
		return
	}

	if gb.properties.VendorFirst {
		vendorDirPath := path.Join(ctx.ModuleDir(), "vendor")
		ctx.Build(pctx, blueprint.BuildParams{
			Description: fmt.Sprintf("Vendor dependencies of %s", name),
			Rule:        goVendor,
			Outputs:     []string{vendorDirPath},
			Implicits:   []string{path.Join(ctx.ModuleDir(), "go.mod")},
			Optional:    true,
			Args: map[string]string{
				"workDir": ctx.ModuleDir(),
				"name":    name,
			},
		})
		inputs = append(inputs, vendorDirPath)
	}

	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Build %s as Go binary", name),
		Rule:        goBuild,
		Outputs:     []string{outputPath},
		Implicits:   inputs,
		Args: map[string]string{
			"outputPath": outputPath,
			"workDir":    ctx.ModuleDir(),
			"pkg":        gb.properties.Pkg,
		},
	})

	if len(gb.properties.TestPkg) != 0 {
		ctx.Build(pctx, blueprint.BuildParams{
			Description: fmt.Sprintf("Test package %s", gb.properties.TestPkg),
			Rule:        goTest,
			Outputs:     []string{testOutputPath},
			Implicits:   append(testInputs, inputs...),
			Args: map[string]string{
				"testOutput": testOutputPath,
				"workDir":    ctx.ModuleDir(),
				"testPkg":    gb.properties.TestPkg,
			},
		})
	}
}

func SimpleBinFactory() (blueprint.Module, []interface{}) {
	mType := &goBinaryModuleType{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
