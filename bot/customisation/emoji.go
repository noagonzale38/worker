package customisation

import (
	"fmt"

	"github.com/rxdn/gdl/objects"
	"github.com/rxdn/gdl/objects/guild/emoji"
)

type CustomEmoji struct {
	Name     string
	Id       uint64
	Animated bool
}

func NewCustomEmoji(name string, id uint64, animated bool) CustomEmoji {
	return CustomEmoji{
		Name: name,
		Id:   id,
	}
}

func (e CustomEmoji) String() string {
	if e.Animated {
		return fmt.Sprintf("<a:%s:%d>", e.Name, e.Id)
	} else {
		return fmt.Sprintf("<:%s:%d>", e.Name, e.Id)
	}
}

func (e CustomEmoji) BuildEmoji() *emoji.Emoji {
	return &emoji.Emoji{
		Id:       objects.NewNullableSnowflake(e.Id),
		Name:     e.Name,
		Animated: e.Animated,
	}
}

var (
	EmojiId         = NewCustomEmoji("ticketID", 1342682254274596904, false)
	EmojiOpen       = NewCustomEmoji("open", 1342682257831497760, false)
	EmojiOpenTime   = NewCustomEmoji("time", 1342682259018223637, false)
	EmojiClose      = NewCustomEmoji("close", 1342682256728264816, false)
	EmojiCloseTime  = NewCustomEmoji("time", 1342682259018223637, false)
	EmojiReason     = NewCustomEmoji("reason", 1342682259974787112, false)
	EmojiSubject    = NewCustomEmoji("subject", 1327350205896458251, false)
	EmojiTranscript = NewCustomEmoji("transcripts", 1342688904393916488, false)
	EmojiClaim      = NewCustomEmoji("claimedby", 1342682255587278969, false)
	EmojiPanel      = NewCustomEmoji("panel", 1306922982815301673, false)
	EmojiRating     = NewCustomEmoji("rating", 1342688907325734954, false)
	EmojiStaff      = NewCustomEmoji("support", 1342688908902797323, false)
	EmojiThread     = NewCustomEmoji("ticketThread", 1342688905987756042, false)
	EmojiBulletLine = NewCustomEmoji("arrowRight", 1342688910211547217, false)
	EmojiPatreon    = NewCustomEmoji("patreon", 1342688911729758250, false)
	EmojiDiscord    = NewCustomEmoji("discord", 1342688913214672977, false)
	//EmojiTime       = NewCustomEmoji("time", 974006684622159952, false)
)

// PrefixWithEmoji Useful for whitelabel bots
func PrefixWithEmoji(s string, emoji CustomEmoji, includeEmoji bool) string {
	if includeEmoji {
		return fmt.Sprintf("%s %s", emoji, s)
	} else {
		return s
	}
}
