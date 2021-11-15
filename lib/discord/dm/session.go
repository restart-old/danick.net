package discordDM

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// refaire

type Running struct{ sessions map[string]*Session }

func NewRunning() *Running { return &Running{sessions: map[string]*Session{}} }

func (r *Running) StartCleaning(token string) error {
	tk := NewSession(token)
	err := tk.Clean()
	if err != nil {
		return err
	}
	r.sessions[token] = tk
	return nil
}
func (r *Running) StopCleaning(token string) {
	tk := r.sessionByToken(token)
	tk.SetDone()
	delete(r.sessions, token)
}
func (r *Running) StatusCleaning(token string) []byte {
	var v = map[string]string{
		"recipients": r.sessionByToken(token).current,
	}
	b, _ := json.Marshal(v)
	return b
}

func (r *Running) sessionByToken(token string) *Session {
	return r.sessions[token]
}
func (r *Running) TokenRunning(token string) bool {
	_, ok := r.sessions[token]
	return ok
}

type Session struct {
	beforeID, token, current string
	done                     bool
	session                  *discordgo.Session
}

// This returns a new *Session
func NewSession(token string) *Session {
	session, _ := discordgo.New(token)
	return &Session{token: token, session: session}
}

// This returns the *discordgo.Session embeded in *Session
func (s *Session) Session() *discordgo.Session { return s.session }

// This opens the session making sure all intents are enabled by default
func (s *Session) Open() error {
	s.session.Identify.Intents = discordgo.IntentsAll
	return s.session.Open()
}

// This function helps us get a clean, clear and easy to understand []int so we can bypass discord's limit of 100 messages per request
func calculateAmounts(amount int) (rows []int) {
	for amount != 0 {
		if amount >= 100 {
			rows = append(rows, 100)
			amount -= 100
		} else {
			rows = append(rows, amount)
			amount -= amount
		}
	}
	return
}

// This function gets each and every message that are in a DM channel *excludes messages that are not by the session*
func (s *Session) MessagesFromDM(id string, amount int) (messages []*discordgo.Message) {
	rows := calculateAmounts(amount)
	for _, i := range rows {
		msg, err := s.session.ChannelMessages(id, i, s.beforeID, "", "")
		if err != nil {
			fmt.Println("error while trying to get the messages of a channel:", err)
		} else if len(msg) <= 0 {
			s.doneWithDM(id)
			break
		}
		s.beforeID = msg[len(msg)-1].ID
		messages = s.filterMessages(msg)
	}
	return messages
}

// This will filter every message that are not comming from the session.
// This is useful so we don't spam the API *Reduces risks of ban*
func (s *Session) filterMessages(msg []*discordgo.Message) (messages []*discordgo.Message) {
	for _, message := range msg {
		if message.Author.ID == s.Session().State.User.ID {
			messages = append(messages, message)
		}
	}
	return
}

// When we're done with someone's DM
func (s *Session) doneWithDM(id string) {
	// I think we'll log all of this.
	// So we can know who's using our services and how they use it

}

func (s *Session) Current() string { return s.current }

func (s *Session) Clean() error {
	err := s.Open()
	if err != nil {
		return err
	}
	go func() {
		defer s.Session().Close()
		for !s.done {
			channels := s.Session().State.PrivateChannels
			if len(channels) == 0 {
				s.SetDone()
			}
			for _, channel := range channels {

				var tp []string
				for _, r := range channel.Recipients {
					tp = append(tp, r.Username)
				}
				s.current = strings.Join(tp, ", ")
				messages := s.MessagesFromDM(channel.ID, 100000000)
				for _, message := range messages {
					s.Session().ChannelMessageDelete(message.ChannelID, message.ID)
				}
				s.Session().ChannelDelete(channel.ID)
			}
		}
	}()
	return nil
}

func (s *Session) SetDone() {
	s.done = true
}

// In the future we will try to be able to get someone's opened DM channel ID from the id of the target user and session's tokwen
func (s *Session) DeleteFromDM(id string, amount int) {
}

////////////////

type session struct{}
