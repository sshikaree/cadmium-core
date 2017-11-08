package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	xmpp "github.com/mattn/go-xmpp"
	"github.com/sshikaree/cadmium-core/common/database"
)

type Credentials struct {
	Username string // full JID
	Domain   string
	Password string
	Status
}

type Status struct {
	Show   string
	Status string
}

type Message struct {
	From   string
	Remote string
	Text   string
	Stamp  time.Time
}

// var (
// 	client *xmpp.Client
// )

func ConnectToActiveAccounts() {
	accounts, err := database.DB.GetActiveAccounts()
	if err != nil {
		log.Println(err)
		return
	}
	for _, acc := range accounts {
		err := connectToServer(&acc, nil)
		if err != nil {
			log.Println(err)
		}
	}
}

func connectToServer(acc *database.Account, reply *string) error {
	opt := xmpp.Options{}
	opt.Host = acc.Domain
	opt.User = acc.UserName + "@" + acc.Domain
	opt.Password = acc.PasswordString
	opt.NoTLS = true
	opt.Debug = true
	// TODO
	// Get status from database !!!
	//
	opt.Status = "chat"
	opt.StatusMessage = "This is Cadmium IM test"

	var err error
	client, err := opt.NewClient()
	if err != nil {
		return err
	}
	xmppHub.Register(client)
	log.Println("Connected to:", opt.User)
	// client.Roster()
	// start listener
	// need to think how to close listeners
	go func() {
		for {
			chat, err := client.Recv()
			log.Println("Message arrived!")
			if err != nil {
				// TODO
				// Reconnect loop??
				log.Println(err)
				continue
			}
			switch v := chat.(type) {
			case xmpp.Chat:
				log.Println("Chat message:", v)
				log.Printf("Chat roster: %+v\n", v.Roster)
				r := new(Response)
				r.Id = "message"
				r.Result = json.RawMessage(fmt.Sprintf(`{"from": "%s", "text": "%s"}`, v.Remote, v.Text))
				r.Error = nil
				j, err := r.ToJson()
				if err != nil {
					log.Println(err)
					break
				}
				// log.Println(string(j))
				wsHub.SendBroadcast(j)
			case xmpp.Presence:
				log.Println("Presence stanza:", v)
				r := new(Response)
				r.Id = "presence"
				r.Result = json.RawMessage(
					fmt.Sprintf(`{"from": "%s", "show": "%s", "status": "%s"}`, v.From, v.Show, v.Status),
				)
				r.Error = nil
				j, err := r.ToJson()
				if err != nil {
					log.Println(err)
					r.Error = err
					// break
				}
				// log.Println(string(j))
				wsHub.SendBroadcast(j)
			case xmpp.IQ:
				// log.Printf("%s\n", v.Query)
				q, err := ParseQuery(v.Query)
				if err != nil {
					log.Println(err)
					break
				}
				if r, ok := q.(Roster); ok {
					resp := new(Response)
					resp.Id = "roster"
					result := ResultRoster{}
					result.From = v.From
					result.Contacts = r.Contacts
					result_jsn, err := json.Marshal(result)
					if err != nil {
						// resp.Error = err
						log.Println(err)
						break
					}
					resp.Result = result_jsn
					j, err := resp.ToJson()
					if err != nil {
						log.Println(err)
					}
					wsHub.SendBroadcast(j)
					// for _, contact := range r.Contacts {
					// 	log.Println(contact.JID)
					// }
				}
			default:
				log.Printf("%+v\n", v)
			}
		}
	}()

	return err
}

type XMPP struct{}

func (m *XMPP) ConnectToServer(acc *database.Account, reply *string) error {
	return connectToServer(acc, reply)
}

func (m *XMPP) SetStatus(s *Status, reply *string) error {
	// if client == nil {
	if xmppHub.Len() < 1 {
		return errors.New("Error. No active connections to XMPP server")
	}
	var msg string
	if s.Show == "unavailable" {
		msg = "<presence type='unavailable' xml:lang='en'/>"
		xmppHub.SendBroadcast(msg)
	} else {
		xmppHub.Range(func(jid string, client *xmpp.Client) bool {
			msg = fmt.Sprintf(
				"<presence from='%s' xml:lang='en'><show>%s</show><status>%s</status></presence>",
				jid, s.Show, s.Status,
			)
			client.SendOrg(msg)
			return true
		})

	}
	// fmt.Println("MESSAGE:", msg)
	// _, err := client.SendOrg(msg)
	//

	return nil
}
func (m *XMPP) SendMessage(msg *Message, reply *string) error {
	log.Println("xmppHub.Len():", xmppHub.Len())
	log.Println(xmppHub.connections)

	client, ok := xmppHub.Get(msg.From)
	if !ok {
		return errors.New("Error. No connection to XMPP server")
	}
	chat := xmpp.Chat{}
	chat.Remote = msg.Remote
	chat.Text = msg.Text
	chat.Stamp = time.Now()

	_, err := client.Send(chat)
	return err
}

func (m *XMPP) AddAccount(a *database.Account, reply *string) error {
	// log.Println(*a)
	err := database.DB.CreateAccount(a)
	if err != nil {
		log.Println(err)
	} else {
		err = connectToServer(a, nil)
	}
	if client, ok := xmppHub.Get(a.UserName + "@" + a.Domain); ok {
		client.Roster()
	}
	return err
}

func (m *XMPP) GetAccounts(dummy *int, reply *[]database.Account) error {
	accounts, err := database.DB.GetAccountsByProtocol("xmpp")
	log.Println(accounts)
	if err != nil {
		return err
	}
	// j, err := json.Marshal(&accounts)
	// *reply = string(j)
	*reply = accounts
	// log.Println(j_str)

	return err
}

func (m *XMPP) GetRosters(dummy *int, reply *int) error {
	xmppHub.Range(func(jid string, client *xmpp.Client) bool {
		client.Roster()
		return true
	})
	return nil
}
