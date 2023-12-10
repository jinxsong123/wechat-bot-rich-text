package sender

import (
	"time"

	"log"

	"github.com/eatmoreapple/openwechat"
)

type wxsender struct {
	friends openwechat.Friends
	groups  openwechat.Groups
	bot     *openwechat.Bot
}

func NewWXsender(friendsNickName, groupsNickName []string) (Sender, error) {
	// os.Remove("token.json")
	bot := openwechat.Default(openwechat.Desktop)
	bot.UUIDCallback = openwechat.PrintlnQrcodeUrl
	reloadStorage := openwechat.NewFileHotReloadStorage("token.json")
	defer reloadStorage.Close()
	err := bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption())

	if err != nil {
		log.Printf("Failed to login to wx, err: %s", err.Error())
		return nil, err
	}
	self, err := bot.GetCurrentUser()
	if err != nil {
		log.Printf("Failed to get current user, err: %s", err)
		return nil, err
	}
	var operatorFriends openwechat.Friends
	var operatorGroups openwechat.Groups
	for _, friend := range friendsNickName {
		if friends, err := self.Friends(); err == nil {
			searcdedFriend := friends.SearchByNickName(1, friend)
			if len(searcdedFriend) > 0 {
				operatorFriends = append(operatorFriends, searcdedFriend...)
			}
		} else {
			log.Printf("Failed to get friends, err: %s", err.Error())
		}
	}

	for _, group := range groupsNickName {
		if groups, err := self.Groups(); err == nil {
			searchededGroup := groups.SearchByNickName(1, group)
			if len(searchededGroup) > 0 {
				operatorGroups = append(operatorGroups, searchededGroup...)
			}

		} else {
			log.Printf("Failed to get groups, err: %s", err.Error())
		}
	}

	sender := wxsender{
		friends: operatorFriends,
		groups:  operatorGroups,
		bot:     bot,
	}
	go sender.alive()
	return &sender, nil
}

func (wx *wxsender) alive() {
	for {
		time.Sleep(10 * time.Second)
		if !wx.bot.Alive() {
			log.Println("The bot is not alive")
			reloadStorage := openwechat.NewFileHotReloadStorage("token.json")
			err := wx.bot.HotLogin(reloadStorage, openwechat.NewRetryLoginOption())
			if err != nil {
				log.Printf("Failed to login to wx, err: %s", err.Error())
			}
			if self, err := wx.bot.GetCurrentUser(); err == nil {
				self.Members(true)
			}
		}
	}

}

func (wx *wxsender) Send(messages []*string) error {

	if len(messages) == 0 {
		return nil
	}
	if len(wx.friends) > 0 {
		for _, friend := range wx.friends {
			for _, message := range messages {
				log.Printf("send message: %s", *message)
				_, err := friend.SendText(*message)
				if err != nil {
					log.Printf("Failed to send message to friend, err := %s \n", err.Error())
				}
			}
		}
	}
	if len(wx.friends) > 0 {
		for _, group := range wx.groups {
			for _, message := range messages {
				log.Printf("send message: %s", *message)
				if _, err := group.SendText(*message); err != nil {
					log.Printf("Failed to send message to group, err := %s \n", err.Error())
				}
				time.Sleep(1 * time.Second)
			}

		}
	}

	return nil
}
