package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// EncryptAES 使用AES加密算法对明文进行加密
// 参数key是加密密钥，plaintext是待加密的明文
// 返回加密后的密文字符串和可能出现的错误
func EncryptAES(key []byte, plaintext string) (string, error) {
	// 创建一个新的AES密码块，使用提供的密钥
	block, err := aes.NewCipher(key)
	if err != nil {
		// 如果创建密码块失败，返回空字符串和错误信息
		return "", err
	}

	// 创建一个字节切片用于存储密文，长度为AES块大小加上明文的长度
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	// iv是初始化向量，使用密码块大小的字节切片
	iv := ciphertext[:aes.BlockSize]
	// 从随机数生成器中读取数据填充初始化向量
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		// 如果读取随机数失败，返回空字符串和错误信息
		return "", err
	}

	// 创建一个新的CFB模式的加密流，使用密码块和初始化向量
	stream := cipher.NewCFBEncrypter(block, iv)
	// 使用加密流对明文进行加密，结果存储在ciphertext的适当位置
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	// 将加密后的密文使用base64编码转换为字符串
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// DecryptAES 使用AES加密算法对密文进行解密
func DecryptAES(key []byte, cryptoText string) (string, error) {
	// 将密文使用base64解码为字节切片
	ciphertext, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	// 创建一个新的AES密码块，使用提供的密钥
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 检查密文长度是否足够包含初始化向量和加密数据
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	// 从密文中提取初始化向量
	iv := ciphertext[:aes.BlockSize]
	// 从密文中提取加密数据
	ciphertext = ciphertext[aes.BlockSize:]

	// 创建一个新的CFB模式的解密流，使用密码块和初始化向量
	stream := cipher.NewCFBDecrypter(block, iv)
	// 使用解密流对加密数据进行解密，结果存储在ciphertext的适当位置
	stream.XORKeyStream(ciphertext, ciphertext)

	// 返回解密后的明文字符串
	return string(ciphertext), nil
}
