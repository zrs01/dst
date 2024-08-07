package yml

import (
	"fmt"
	"os"

	"github.com/samber/lo"
	"github.com/zrs01/dst/model"
	"github.com/ztrue/tracerr"
	yamlIn "gopkg.in/yaml.v3"
)

func loadFromYml(file string) (*model.DataDef, error) {
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	var d model.DataDef
	if err := yamlIn.Unmarshal(yamlFile, &d); err != nil {
		return nil, tracerr.Wrap(err)
	}

	updateReferenceTables(&d)
	expandFixColumns(&d)

	// validate	data
	validateResult := model.Verify(&d)
	if len(validateResult) > 0 {
		lo.ForEach(validateResult, func(v string, _ int) {
			fmt.Println(v)
		})
		return nil, tracerr.Errorf("invalid data")
	}

	return &d, err
}
