package labdoctemplate

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/apex/log"

	// Register all plugins.
	_ "github.com/drewstinnett/labdoc/internal/plugins/all"

	"github.com/drewstinnett/labdoc/pkg/labdoc"
	"github.com/spf13/cobra"
)

func Generate(inf, outf string) (bool, error) {
	var changed bool
	templateIn, err := ioutil.ReadFile(inf)
	if err != nil {
		return false, err
	}

	originalHash, err := labdoc.GetSha256(outf)
	if err != nil {
		log.Debugf("Couldn't get a hash on file %v: %v", outf, err)
	}
	log.Debug(fmt.Sprintf("%x", originalHash))

	allFunctions := template.FuncMap{}

	// Load All plugin data
	for name, c := range labdoc.Plugins {
		itemFunction, err := c().TemplateFunctions()
		if err != nil {
			return false, err
		}
		for k, v := range itemFunction {
			funcName := fmt.Sprintf("%v%v", name, strings.Title(k))
			allFunctions[funcName] = v
		}
	}

	tpl, err := template.New("tpl").Funcs(allFunctions).Parse(string(templateIn))
	if err != nil {
		return false, err
	}
	var out *os.File
	if outf == "" {
		out = os.Stdout
	} else {
		out, err = os.Create(outf)
		cobra.CheckErr(err)
	}
	err = tpl.Execute(out, nil)
	if err != nil {
		return false, err
	}

	newHash, err := labdoc.GetSha256(outf)
	if err != nil {
		return false, err
	}

	if originalHash != newHash {
		changed = true
	}
	return changed, nil
}
