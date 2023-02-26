package bubblewrap

// CancelError is returned when the user indicates they want the input to be canceled,
// by using the escape or Control C inputs.
type CancelError error
