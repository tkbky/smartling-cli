package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/Smartling/api-sdk-go"
	"github.com/reconquest/hierr-go"
)

func doFilesImport(
	client *smartling.Client,
	config Config,
	args map[string]interface{},
) error {
	var (
		project     = config.ProjectID
		uri         = args["<uri>"].(string)
		file        = args["<file>"].(string)
		locale      = args["<locale>"].(string)
		fileType, _ = args["--type"].(string)
	)

	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return NewError(
			hierr.Errorf(err, "unable to read file for import"),
			"Check that specified file exists and you have permissions "+
				"to read it.",
		)
	}

	request := smartling.ImportRequest{}

	request.File = contents
	request.FileURI = uri

	if args["--published"].(bool) {
		request.TranslationState = smartling.TranslationStatePublished
	}

	if args["--post-translation"].(bool) {
		request.TranslationState = smartling.TranslationStatePostTranslation
	}

	if args["--overwrite"].(bool) {
		request.Overwrite = true
	}

	if fileType != "" {
		request.FileType = smartling.FileType(fileType)
	} else {
		request.FileType = smartling.GetFileTypeByExtension(
			filepath.Ext(file),
		)

		if request.FileType == smartling.FileTypeUnknown {
			return NewError(
				fmt.Errorf(
					"unable to deduce file type from extension: %q",
					filepath.Ext(file),
				),

				`You need to specify file type via --type option.`,
			)
		}
	}

	result, err := client.Import(project, locale, request)
	if err != nil {
		return hierr.Errorf(
			err,
			`unable to import file "%s" (original "%s")`,
			file,
			uri,
		)
	}

	fmt.Printf("XXXXXX do_files_import.go:57 string(result): %#v\n", string(result))

	return nil
}
