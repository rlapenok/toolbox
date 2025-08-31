package errors

// Details is a struct that contains the details of the error
type Details struct {
	service        string
	domain         string
	localeMessages map[string]string
}

// NewDetails creates a new Details struct
func NewDetails() *Details {
	return &Details{
		localeMessages: make(map[string]string),
	}
}

// WithService adds a service to the Details struct
func (d *Details) WithService(service string) *Details {
	d.service = service
	return d
}

// WithDomain adds a domain to the Details struct
func (d *Details) WithDomain(domain string) *Details {
	d.domain = domain
	return d
}

// WithLocaleMessage adds a locale message to the Details struct
func (d *Details) WithLocaleMessage(locale string, message string) *Details {
	d.localeMessages[locale] = message
	return d
}

// LocaleMessage returns the locale message for the given locale
func (d *Details) LocaleMessage(locale string) string {
	return d.localeMessages[locale]
}
