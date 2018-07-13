package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strings"
)

// ErrNoAvatarURL 은 아바타 URL을 제공할 수 없을 때 발생하는 에러
var ErrNoAvatarURL = errors.New("Chat: Unable to get an avatar url")

// Avatar 는 사용자 프로필 사진을 표현할 수 있는 구조체
type Avatar interface {
	// GetAvatarURL은 지정된 클라이언트에 대한 아바타 URL을 가져옴
	// 지정된 클라이언트의 아바타 URL을 가져오지 못할 경우, ErrNoAvatarURL이 반환됨
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}

	return "", ErrNoAvatarURL
}

type GravatarAvatar struct{}

var UseGravatarAvatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if email, ok := c.userData["email"]; ok {
		if emailStr, ok := email.(string); ok {
			m := md5.New()
			io.WriteString(m, strings.ToLower(emailStr))
			return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil)), nil
		}
	}

	return "", ErrNoAvatarURL
}
