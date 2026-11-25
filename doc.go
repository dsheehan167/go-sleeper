// Package sleeper is a client for the Sleeper fantasy sports API.
//
// Create a Client with NewClient and call its Get* and List* methods. Every
// request goes through a built-in rate limiter that keeps usage within the
// limits Sleeper documents. See https://docs.sleeper.com for the API reference.
package sleeper
