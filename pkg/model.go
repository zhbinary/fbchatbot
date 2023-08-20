// Created by zhbinary on 2023/8/19.
package pkg

type WebhookEvent struct {
	Object string   `json:"object"`
	Entry  []*Entry `json:"entry"`
}

type Entry struct {
	ID        string      `json:"id"`
	Time      int64       `json:"time"`
	Messaging []Messaging `json:"messaging"`
}

type Messaging struct {
	Sender    *Sender    `json:"sender"`
	Recipient *Recipient `json:"recipient"`
	Message   *Message   `json:"message"`
	FeedBack  *Feedback  `json:"messaging_feedback"`
}

type Sender struct {
	ID string `json:"id"`
}

type Recipient struct {
	ID string `json:"id"`
}

type SendRequest struct {
	Recipient     *Recipient `json:"recipient"`
	MessagingType string     `json:"messaging_type"`
	Message       *Message   `json:"message"`
	AccessToken   string     `json:"access_token"`
}

type Message struct {
	Text         string        `json:"text,omitempty"`
	Attachment   *Attachment   `json:"attachment,omitempty"`
	QuickReplies []*QuickReply `json:"quick_replies,omitempty"`
	Nlp          *Nlp          `json:"nlp"`
}

type Nlp struct {
	Traits *Traits `json:"traits"`
}

type Traits struct {
	Sentiment []*Sentiment `json:"wit$sentiment"`
}

type Sentiment struct {
	Id         string
	Value      string
	Confidence float32
}

type Attachment struct {
	Type    string   `json:"type"`
	Payload *Payload `json:"payload,omitempty"`
}

type Payload struct {
	TemplateType string `json:"template_type,omitempty"`
	// other template or file content fields
}

type QuickReply struct {
	// quick reply fields
}

type Feedback struct {
}
