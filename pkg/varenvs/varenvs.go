package varenvs

import (
	"fmt"
	"os"
)

type varenv struct {
	name              string
	varp              *string
	envVarName        string
	defaultValue      string
	required          bool
	priorityToPointer bool
}

type Varenvs struct {
	varenvs []varenv
}

func New() *Varenvs {
	return &Varenvs{
		varenvs: make([]varenv, 0, 5),
	}
}

func (v *Varenvs) Add(name string, varp *string, envVarName string, defaultValue string, required bool, priorityToPointer bool) {
	v.varenvs = append(v.varenvs, varenv{
		name:              name,
		varp:              varp,
		envVarName:        envVarName,
		defaultValue:      defaultValue,
		required:          required,
		priorityToPointer: priorityToPointer,
	})
}

func (v *Varenvs) Parse() error {
	missings := make([]*varenv, 0, 0)
	for idx, ve := range v.varenvs {
		if ve.priorityToPointer && *ve.varp != "" {
			// Keep pointed value. Don't care about ENV value
		} else {
			v := os.Getenv(ve.envVarName)
			if v != "" {
				*ve.varp = v
			}
		}
		// Now, we are set. Check required error
		if ve.required && *ve.varp == "" {
			missings = append(missings, &v.varenvs[idx])
		}
	}
	if len(missings) > 0 {
		sep := ""
		l1 := ""
		l2 := ""
		for idx, _ := range missings {
			l1 = l1 + sep + missings[idx].name
			l2 = l2 + sep + missings[idx].envVarName
			sep = ", "
		}
		return fmt.Errorf("missing configuration variable: '%s' or ENV(%s)", l1, l2)
	}
	return nil
}
