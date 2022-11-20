// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: pb/v1/resource.proto

package v1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on PodsRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *PodsRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PodsRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in PodsRequestMultiError, or
// nil if none found.
func (m *PodsRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *PodsRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return PodsRequestMultiError(errors)
	}

	return nil
}

// PodsRequestMultiError is an error wrapping multiple validation errors
// returned by PodsRequest.ValidateAll() if the designated constraints aren't met.
type PodsRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PodsRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PodsRequestMultiError) AllErrors() []error { return m }

// PodsRequestValidationError is the validation error returned by
// PodsRequest.Validate if the designated constraints aren't met.
type PodsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PodsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PodsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PodsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PodsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PodsRequestValidationError) ErrorName() string { return "PodsRequestValidationError" }

// Error satisfies the builtin error interface
func (e PodsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPodsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PodsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PodsRequestValidationError{}

// Validate checks the field values on PodsResponse with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *PodsResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PodsResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in PodsResponseMultiError, or
// nil if none found.
func (m *PodsResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *PodsResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return PodsResponseMultiError(errors)
	}

	return nil
}

// PodsResponseMultiError is an error wrapping multiple validation errors
// returned by PodsResponse.ValidateAll() if the designated constraints aren't met.
type PodsResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PodsResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PodsResponseMultiError) AllErrors() []error { return m }

// PodsResponseValidationError is the validation error returned by
// PodsResponse.Validate if the designated constraints aren't met.
type PodsResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PodsResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PodsResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PodsResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PodsResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PodsResponseValidationError) ErrorName() string { return "PodsResponseValidationError" }

// Error satisfies the builtin error interface
func (e PodsResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPodsResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PodsResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PodsResponseValidationError{}
