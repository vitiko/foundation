package reflectx

import (
	"fmt"
	"reflect"

	"github.com/hyperledger/fabric-chaincode-go/shim"
)

// Validator is an interface that can be implemented by types that can validate themselves.
type Validator interface {
	Validate(stub shim.ChaincodeStubInterface) error
}

// ValidateArguments validates the arguments for the specified method on the given value using reflection.
// It checks whether the specified method exists on the value 'v' and if the number of provided arguments matches
// the method's expected input parameters. Additionally, it attempts to convert the string arguments into the
// expected types using various unmarshaling or decoding interfaces such as JSON, proto.Message,
// encoding.TextUnmarshaler, encoding.BinaryUnmarshaler, and gob.GobDecoder. If an argument implements the Validator
// interface, its Validate method is called with the provided stub.
//
// The function returns an error if the method is not found, the number of arguments is incorrect, or if an error
// occurs during argument conversion or validation.
//
// Parameters:
//   - v: The value on which the method is to be validated.
//   - method: The name of the method to validate.
//   - stub: The ChaincodeStubInterface used for access control checks.
//   - args: A slice of strings representing the arguments for the method.
//
// Returns:
//   - error: An error if the method is not found, the number of arguments is incorrect, or any other issue during validation.
func ValidateArguments(v any, method string, stub shim.ChaincodeStubInterface, args ...string) error {
	inputVal := reflect.ValueOf(v)

	methodVal := inputVal.MethodByName(method)
	if !methodVal.IsValid() {
		return fmt.Errorf("%w: %s", ErrMethodNotFound, method)
	}

	methodType := methodVal.Type()
	if methodType.NumIn() != len(args) {
		return fmt.Errorf(
			"%w: found %d but expected %d: validate %s",
			ErrIncorrectArgumentCount,
			len(args),
			methodType.NumIn(),
			method,
		)
	}

	for i, arg := range args {
		value, err := valueOf(arg, methodType.In(i))
		if err != nil {
			return fmt.Errorf("%w: validate %s, argument %d", err, method, i)
		}

		if validator, ok := value.Interface().(Validator); ok {
			if err := validator.Validate(stub); err != nil {
				return fmt.Errorf(
					"%w: '%s': validation failed: '%v': validate %s, argument %d",
					ErrInvalidArgumentValue,
					arg,
					err.Error(),
					method,
					i,
				)
			}
		}
	}

	return nil
}