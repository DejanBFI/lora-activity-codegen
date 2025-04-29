package main

import (
	"fmt"
	"strings"

	"lora-activity-codegen/generator"

	"github.com/dave/jennifer/jen"
)

func main() {
	schemaName := ""
	processAndActivityName := "set_survey_data_by_assignment_rule"

	readSet := []string{
		"DocLoanStructureRiskLevel",
		"DocLoanStructureProductId",
	}

	writeSet := []string{
		"DocSurveyAppointmentSurveyType",
		"DocSurveyAppointmentSurveyLocationType",
		"DocSurveyAppointmentSurveyResourceType",
	}

	// DO NOT EDIT BELOW THIS LINE
	packageName := strings.ReplaceAll(processAndActivityName, "_", "")
	out := jen.NewFile(packageName)
	generator.GenerateBoilerplate(out, schemaName, processAndActivityName, readSet, writeSet)

	fmt.Printf("%#v", out)
}
