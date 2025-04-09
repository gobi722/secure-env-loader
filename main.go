package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const userName = "salispay"

// func getIV(iv string) []byte {
// 	ivBytes := []byte(iv)
// 	if len(ivBytes) >= 16 {
// 		return ivBytes[:16]
// 	}
// 	paddedIV := make([]byte, 16)
// 	copy(paddedIV, ivBytes)
// 	return paddedIV
// }

// func hashData(data string) string {
// 	hash := md5.Sum([]byte(data)) // Mimicking hashing used in JS
// 	return hex.EncodeToString(hash[:])
// }

// func pkcs7Pad(data []byte, blockSize int) []byte {
// 	padLen := blockSize - (len(data) % blockSize)
// 	padText := bytes.Repeat([]byte{byte(padLen)}, padLen)
// 	return append(data, padText...)
// }

// func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
// 	length := len(data)
// 	if length == 0 {
// 		return nil, errors.New("empty data")
// 	}
// 	padLen := int(data[length-1])
// 	if padLen > blockSize || padLen == 0 {
// 		return nil, errors.New("invalid padding size")
// 	}
// 	return data[:length-padLen], nil
// }
// func pkcs7Unpad(data []byte) ([]byte, error) {
// 	length := len(data)
// 	if length == 0 {
// 		return nil, errors.New("invalid padding")
// 	}
// 	padLen := int(data[length-1])
// 	if padLen > length || padLen > aes.BlockSize {
// 		return nil, errors.New("invalid padding")
// 	}
// 	return data[:length-padLen], nil
// }

// func getIV(iv string) []byte {
// 	ivBuffer := []byte(iv)
// 	if len(ivBuffer) >= 16 {
// 		return ivBuffer[:16]
// 	}
// 	padding := make([]byte, 16-len(ivBuffer))
// 	return append(ivBuffer, padding...)
// }

//	func hashData(data string) string {
//		hash := sha256.Sum256([]byte(data))
//		return hex.EncodeToString(hash[:])[:32] // First 32 bytes
//	}
func getIV(iv string) []byte {
	ivBytes := []byte(iv)
	if len(ivBytes) >= 16 {
		return ivBytes[:16]
	}
	paddedIV := make([]byte, 16)
	copy(paddedIV, ivBytes)
	return paddedIV
}

func hashData(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])[:32]
}
func decryptData(data string, key string) (string, error) {
	tokenFromUI := hashData(userName)
	keyBytes := []byte(tokenFromUI)
	ivBytes := getIV(key)
	encryptedData, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, ivBytes)
	decrypted := make([]byte, len(encryptedData))
	mode.CryptBlocks(decrypted, encryptedData)

	// Remove possible padding
	decrypted = removePadding(decrypted)
	decryptedKey1 := strings.Trim(string(decrypted), "\"")
	return decryptedKey1, nil
}

func removePadding(data []byte) []byte {
	padLen := int(data[len(data)-1])
	if padLen > len(data) {
		return data
	}
	return data[:len(data)-padLen]
}

func loadKeyFromJSON(filename string) (string, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	var data map[string]string
	if err := json.Unmarshal(file, &data); err != nil {
		return "", err
	}
	return data["key"], nil
}

func loadEncryptedTextFromEnv(filename string) (string, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ENCRYPTION_KEY=") {
			return strings.TrimPrefix(line, "ENCRYPTION_KEY="), nil
		}
	}
	return "", errors.New("ENCRYPTION_KEY not found in mem.env")
}

func decryptEnvFile(filename, key string) (map[string]string, error) {

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	// config := make(map[string]string)
	decryptedData:= make(map[string]string)
	for _, line := range lines {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			// config[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			decryptedValue, err := decryptData( strings.TrimSpace(parts[1]), key)
			if err != nil {
				return nil, err
			}
			decryptedData[strings.TrimSpace(parts[0])] = decryptedValue
		}
	}
	fmt.Println(decryptedData)
	// file, err := ioutil.ReadFile(filename)
	// if err != nil {
	// 	return nil, err
	// }

	// lines := strings.Split(string(file), "\n")
	// fmt.Println("---------",lines)
	// decryptedData := make(map[string]string)

	// for _, line := range lines {
	// 	if strings.TrimSpace(line) == "" {
	// 		continue
	// 	}
	// 	parts := strings.SplitN(line, "=", 2)
	// 	if len(parts) != 2 {
	// 		continue
	// 	}
	// 	keyName := parts[0]
	// 	encryptedValue := parts[1]
	// 	// fmt.Println(keyName,"----", encryptedValue,key)
		// decryptedValue, err := decryptData(encryptedValue, key)
		// if err != nil {
		// 	return nil, err
		// }
		// decryptedData[keyName] = decryptedValue
	// }

	 return decryptedData, nil
}
func BuildMongoDBURL(username, password, clusterURL, database string) string {

	clusterURL = strings.TrimPrefix(clusterURL, "mongodb+srv://")
	
	
	return fmt.Sprintf("mongodb+srv://%s:%s@%s/%s", username, password, clusterURL, database)
}
func main() {
	key, err := loadKeyFromJSON("encryption_key.json")
	if err != nil {
		fmt.Println("Error loading key:", err)
		return
	}
	fmt.Println("key", key)
	encryptedText, err := loadEncryptedTextFromEnv("MEK.env")
	if err != nil {
		fmt.Println("Error loading encrypted text:", err)
		return
	}
	fmt.Println("encryptedText", encryptedText)
	decryptedKey, err := decryptData(encryptedText, hashData(key))
	if err != nil {
		fmt.Println("Error decrypting key:", err)
		return
	}
	fmt.Printf("Raw Decrypted Bytes: %q\n", decryptedKey)

	decryptedEnv, err := decryptEnvFile("db-config.env", "2dc54acc8fe7bf6f44407413c9359524")
	if err != nil {
		fmt.Println("Error decrypting db-config.env:", err)
		return
	}
	fmt.Println("Decrypted Environment Variables:")
	for k, v := range decryptedEnv {
		fmt.Printf("%s=%s\n", k, v)
	}
	fmt.Println(BuildMongoDBURL(decryptedEnv["USERNAME"], decryptedEnv["PASSWORD"], decryptedEnv["CONNECTION_STRING"], "SWITCH_MMS"))
	// fmt.Println("Decrypted Environment Variables:")
	// for k, v := range decryptedEnv {
	// 	fmt.Printf("%s=%s\n", k, v)
	// }
}
