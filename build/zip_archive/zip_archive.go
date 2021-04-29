package zip_archive

import (
	"fmt"
	"path"
	"strings"

	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
)

var (
	pctx = blueprint.NewPackageContext("github.com/FictProger/architecture2-lab-1/build/zip_archive")

	zipRule = pctx.StaticRule("zipArchive", blueprint.RuleParams{
		Command:     "mkdir $outputPath && zip $outputFile $files",
		Description: "zipping into $outputFile",
	}, "workDir", "outputPath", "outputFile", "files")
)
type zipArchiveType struct {
	blueprint.SimpleName

	properties struct {
		Name string
		Srcs []string
	}
}

func (zipper *zipArchiveType) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	config := bood.ExtractConfig(ctx)
	outputPath := path.Join(config.BaseOutputDir, "archives")
	outputFile := "./" + path.Join(outputPath, zipper.properties.Name) + ".zip"

	var inputs []string

	for _, src := range zipper.properties.Srcs {
		if matches, err := ctx.GlobWithDeps(src, nil); err == nil {
			inputs = append(inputs, matches...)
		}
	}
	for i, _ := range inputs {
		inputs[i] = "./" + inputs[i]
	}
	filesStr := strings.Join(inputs, " ")

	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Archiving into '%s'", name),
		Rule:        zipRule,
		Outputs:     []string{outputPath},
		Args: map[string]string{
			"workDir":    ctx.ModuleDir(),
			"outputPath": outputPath,
			"outputFile": outputFile,
			"files":      filesStr,
		},
	})
}

func SimpleZipArchiveFactory() (blueprint.Module, []interface{}) {
	mType := &zipArchiveType{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
