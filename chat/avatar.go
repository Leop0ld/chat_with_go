package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// ErrNoAvatarURL 은 아바타 URL을 제공할 수 없을 때 발생하는 에러
var ErrNoAvatarURL = errors.New("Chat: Unable to get an avatar url")

// Avatar 는 사용자 프로필 사진을 표현할 수 있는 구조체
type Avatar interface {
	// GetAvatarURL은 지정된 클라이언트에 대한 아바타 URL을 가져옴
	// 지정된 클라이언트의 아바타 URL을 가져오지 못할 경우, ErrNoAvatarURL이 반환됨
	GetAvatarURL(ChatUser) (string, error)
}

type TryAvatars []Avatar

func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, err
		}
	}
	return "", ErrNoAvatarURL
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if len(url) == 0 {
		return "", ErrNoAvatarURL
	}
	return u.AvatarURL(), nil
}

type GravatarAvatar struct{}

var UseGravatarAvatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	files, err := ioutil.ReadDir("avatars")
	if err != nil {
		return "", ErrNoAvatarURL
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fname := file.Name()
		if u.UniqueID() == strings.TrimSuffix(fname, filepath.Ext(fname)) {
			return "/avatars/" + fname, nil
		}
	}
	return "", ErrNoAvatarURL
}
