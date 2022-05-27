package cmd

import (
	"fmt"
	"strings"
)

// CustomHeader represents a single HTTP header name/value pair to be attached to
// outbound HTTP requests against the Connect REST API.
type CustomHeader struct {
	Name  string
	Value string
}

func (h *CustomHeader) Parse(raw string) error {
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid header flag format, expected '<name>:<value>'")
	}

	h.Name = parts[0]
	h.Value = strings.TrimLeft(parts[1], " ")
	return nil
}

func (h *CustomHeader) String() string {
	return fmt.Sprintf("%s: %s", h.Name, h.Value)
}

// HeadersFlag implements the Value interface from pflag in order to support parsing
// flag values of the form '<name>:<value>' directly into a list of CustomHeader
// instances.
type HeadersFlag struct {
	Headers []CustomHeader
}

func (h *HeadersFlag) String() string {
	var pairs []string
	for _, header := range h.Headers {
		pairs = append(pairs, header.String())
	}
	return strings.Join(pairs, ", ")
}

func (h *HeadersFlag) Set(raw string) error {
	header := CustomHeader{}
	if err := header.Parse(raw); err != nil {
		return err
	}
	h.Headers = append(h.Headers, header)
	return nil
}

func (h *HeadersFlag) Type() string {
	return "<name>:<value>"
}
