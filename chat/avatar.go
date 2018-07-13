package main

import "errors"

// ErrNoAvatarURL 은 아바타 URL을 제공할 수 없을 때 발생하는 에러
var ErrNoAvatarURL = errors.New("Chat: Unable to get an avatar url")

// Avatar 는 사용자 프로필 사진을 표현할 수 있는 구조체
type Avatar interface {
	// GetAvatarURL은 지정된 클라이언트에 대한 아바타 URL을 가져옴
	// 지정된 클라이언트의 아바타 URL을 가져오지 못할 경우, ErrNoAvatarURL이 반환됨
	GetAvatarURL(c *client) (string, error)
}

// AuthAvatar 는 Avatar 인터페이스 구현체
type AuthAvatar struct{}

// UseAuthAvatar 는 Avatar 인터페이스 타입을 찾는 용도
var UseAuthAvatar AuthAvatar

// GetAvatarURL 구현
func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}

	return "", ErrNoAvatarURL
}
