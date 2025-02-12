package main

import (
	"fmt"
	"strings"

	"lora-activity-codegen/generator"

	"github.com/dave/jennifer/jen"
)

func main() {
	schemaName := "mobile-phone-usage-age-v1"
	processAndActivityName := "check_phone_usage_age"

	readSet := []string{
		"DocCustomerKtpNik",
		"DocCustomerContactMobileNumber",
	}

	writeSet := []string{
		"DocProcessPhoneCheckResultIsVerifiedPhoneNikUsageAge",
		"DocProcessPhoneCheckResultUsageAge",
	}

	// DO NOT EDIT BELOW THIS LINE
	packageName := strings.ReplaceAll(processAndActivityName, "_", "")
	out := jen.NewFile(packageName)
	generator.GenerateBoilerplate(out, schemaName, processAndActivityName, readSet, writeSet)

	fmt.Printf("%#v", out)
}
