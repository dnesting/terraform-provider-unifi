package provider

import (
	"errors"
	"fmt"

	"github.com/paultyng/go-unifi/unifi"
)

// improveError tries to find a more descriptive error message than the
// plain apiErr that we got from the server.
func improveError(err error) error {
	apiErr := &unifi.APIError{}
	if errors.As(err, &apiErr) {
		switch apiErr.Message {

		case "api.err.NoEdit":
			return errors.New("Resource cannot be modified. Consider using 'terraform taint' to force it to be recreated.")

		case "api.err.InvalidPayload":
			// The controller didn't like a field we supplied.  This is probably
			// a bug in go-unifi but we should give enough information here to
			// let someone address it.
			if v := apiErr.ValidationError; v != nil {
				// XXX: These field names may not match the names in Terraform.
				return fmt.Errorf("API field %q must match pattern %q", v.Field, v.Pattern)
			}

		case "api.err.ObjectReferredByDevice":
			return fmt.Errorf("Resource is in use by %s %q (%s). Consider using 'lifecycle { create_before_destroy = true }'.",
				apiErr.Type, apiErr.AliasOrMAC, apiErr.DeviceID)
		default:
			return fmt.Errorf("%s: an unexpected API error occurred from the UniFi controller", apiErr.Message)
		}
	}
	return err
}
