// Created by zhbinary on 2023/8/20.
package pkg

import "testing"
import . "github.com/smartystreets/goconvey/convey"

func TestHandleMessage(t *testing.T) {
	Convey("", t, func() {
		msg := &Message{}
		So(HandleMessage("1", msg), ShouldBeNil)

		msg.Text = "xxx"
		tmp := HandleMessage("123", msg).(map[string]interface{})
		So(tmp["text"], ShouldNotBeNil)

		msg.Text = "xxx review sdij32"
		tmp = HandleMessage("123", msg).(map[string]interface{})
		So(tmp["text"], ShouldBeNil)
		So(tmp["attachment"], ShouldNotBeNil)

		msg.Text = "hi"
		msg.Nlp = &Nlp{
			Traits: &Traits{
				Sentiment: []*Sentiment{
					{
						Id:         "23",
						Value:      "negative",
						Confidence: 0.97,
					},
				},
			},
		}
		tmp = HandleMessage("123", msg).(map[string]interface{})
		So(tmp["text"], ShouldContainSubstring, "Sorry sir")
	})
}
