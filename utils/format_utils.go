package utils

import (
	"bufio"
	"encoding/json"
	"errors"
	"html/template"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/golang/glog"
)

var templates = map[string]string{
	"utils.SingleVersionPackageDiffResult":    SingleVersionDiffOutput,
	"utils.MultiVersionPackageDiffResult":     MultiVersionDiffOutput,
	"utils.HistDiffResult":                    HistoryDiffOutput,
	"utils.DirDiffResult":                     FSDiffOutput,
	"utils.ListAnalyzeResult":                 ListAnalysisOutput,
	"utils.MultiVersionPackageAnalyzeResult":  MultiVersionPackageOutput,
	"utils.SingleVersionPackageAnalyzeResult": SingleVersionPackageOutput,
}

func JSONify(diff interface{}) error {
	diffBytes, err := json.MarshalIndent(diff, "", "  ")
	if err != nil {
		return err
	}
	f := bufio.NewWriter(os.Stdout)
	defer f.Flush()
	f.Write(diffBytes)
	return nil
}

func getTemplate(diff interface{}) (string, error) {
	diffType := reflect.TypeOf(diff).String()
	if template, ok := templates[diffType]; ok {
		return template, nil
	}
	return "", errors.New("No available template")
}

func TemplateOutput(diff interface{}) error {
	outputTmpl, err := getTemplate(diff)
	if err != nil {
		glog.Error(err)

	}
	funcs := template.FuncMap{"join": strings.Join}
	tmpl, err := template.New("tmpl").Funcs(funcs).Parse(outputTmpl)
	if err != nil {
		glog.Error(err)
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 8, 8, 8, ' ', 0)
	err = tmpl.Execute(w, diff)
	if err != nil {
		glog.Error(err)
		return err
	}
	w.Flush()
	return nil
}
