package cryptography

// Hashable represents an object which can be hashed. Usually these objects
// are io.Serializable and signable. They tend to cache the hash inside for
// effectiveness, providing this accessor method. Anything that can be
// identified with a hash can then be signed and verified.
type Hashable interface {
	Hash() Hash
}
