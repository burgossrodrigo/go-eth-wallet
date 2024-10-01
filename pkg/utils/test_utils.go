package utils

import (
    "testing"
)

func TestDecriptIt(t *testing.T) {
    key := "this is a very secret key"
    plaintext := "this is some text to encrypt"

    ciphertext, err := encriptIt(key, plaintext)
    if err != nil {
        t.Fatalf("encriptIt failed: %v", err)
    }

    decrypted, err := decriptIt(key, ciphertext)
    if err != nil {
        t.Fatalf("decriptIt failed: %v", err)
    }

    if decrypted != plaintext {
        t.Errorf("Decrypted text doesn't match the original plaintext. got: %s, want: %s", decrypted, plaintext)
    }
}