package errors

import "bytes"

type Multi struct {
	happened bool
	causes   []error
}

func Multiple() Multi {
	return Multi{
		happened: false,
		causes:   make([]error, 0),
	}
}

func (m Multi) Error() string {
	var buf bytes.Buffer
	for _, err := range m.causes {
		buf.WriteString(err.Error())
	}
	buf.WriteString("\n")
	return buf.String()
}

func (m Multi) Happen(err error) Multi {
	if err != nil {
		m.causes = append(m.causes, err)
	}
	return Multi{
		causes: m.causes,
	}
}

func (m Multi) Return() error {
	if !m.happened {
		return nil
	}
	return m
}
