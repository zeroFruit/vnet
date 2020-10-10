package link

import "bytes"

type multipleErr struct {
	happened bool
	causes   []error
}

func MultipleErr() multipleErr {
	return multipleErr{
		happened: false,
		causes:   make([]error, 0),
	}
}

func (be multipleErr) Error() string {
	var buf bytes.Buffer
	for _, err := range be.causes {
		buf.WriteString(err.Error())
	}
	buf.WriteString("\n")
	return buf.String()
}

func (be multipleErr) Happen(err error) multipleErr {
	if err != nil {
		be.causes = append(be.causes, err)
	}
	return multipleErr{
		causes: be.causes,
	}
}

func (be multipleErr) Return() error {
	if !be.happened {
		return nil
	}
	return be
}
