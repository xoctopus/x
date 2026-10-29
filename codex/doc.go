// Package codex identifies errors by typed integer codes.
//
// A code answers one question: which error is this. codex carries the code with
// the error and lets callers find it again. It does not define storage, encoding
// or presentation of codes. Capabilities such as an HTTP status code are declared
// on the code type and interpreted by whoever consumes them.
//
// # Code Types
//
// A code type is an integer type satisfying Code, whose values always fit in
// int64. It may optionally implement:
//
//   - ErrorMessage: describes the code; used as the error message.
//   - CanBeCode: lists all values of the code type.
//
// # Creating Errors
//
//   - New: an error carrying a code.
//   - Errorf: like New, with a formatted message.
//   - Wrap, Wrapf: an error carrying a code and caused by another error. They
//     return nil if the cause is nil.
//
// # Matching with a Known Code Type
//
//   - IsCode: whether the chain carries a specific code.
//   - Is: whether the chain carries any code of type C.
//   - As, AsCode: the first error or code of type C in the chain.
//
// # Consuming with an Unknown Code Type
//
// Callers such as transport layers or loggers see only `error`. They define
// their own capability interfaces and look them up on the carried code:
//
//   - errors.As, errors.AsType (recommended): an error carrying a code can be
//     viewed as its code, so the lookup goes through the whole chain. The target
//     must be an interface type, or implement error; errors.AsType requires the
//     latter, so the code type must implement error to be found by it.
//   - Valuer: returns the carried code of a single error, for callers that walk
//     the chain by themselves.
//
// See the package examples for usage.
package codex
