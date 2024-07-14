package utils

import (
	"sync"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
)

var (
	cooldowns  = make(map[string]map[string]int64)
	globalCD   = make(map[string]int64)
	cooldownMu sync.Mutex
)

func CooldownCanContinue(user twitch.User, cmdName string, cmdCooldown, globalCDTime int) bool {
	/* if isUserPermitted(user, []string{"broadcaster", "moderator"}) {
		return true
	} */

	if !globalCooldown(cmdName, globalCDTime) {
		return false
	}

	if !addCooldown(user.Name, cmdName, cmdCooldown) {
		return false
	}

	return true
}

func addCooldown(user string, commandName string, cdTime int) bool {
	cooldownMu.Lock()
	defer cooldownMu.Unlock()

	if _, exists := cooldowns[commandName]; !exists {
		cooldowns[commandName] = make(map[string]int64)
	}

	currentTime := time.Now().UnixMilli()
	cooldownAmount := int64(cdTime) * 1000

	if expireTime, exists := cooldowns[commandName][user]; exists {
		if currentTime < expireTime {
			return false
		}
	}

	cooldowns[commandName][user] = currentTime + cooldownAmount
	time.AfterFunc(time.Duration(cooldownAmount)*time.Millisecond, func() {
		cooldownMu.Lock()
		delete(cooldowns[commandName], user)
		cooldownMu.Unlock()
	})

	return true
}

func globalCooldown(commandName string, cdTime int) bool {
	cooldownMu.Lock()
	defer cooldownMu.Unlock()

	currentTime := time.Now().UnixMilli()
	cooldownAmount := int64(cdTime) * 1000

	if expireTime, exists := globalCD[commandName]; exists {
		if currentTime < expireTime {
			return false
		}
	}

	globalCD[commandName] = currentTime + cooldownAmount
	time.AfterFunc(time.Duration(cooldownAmount)*time.Millisecond, func() {
		cooldownMu.Lock()
		delete(globalCD, commandName)
		cooldownMu.Unlock()
	})

	return true
}
