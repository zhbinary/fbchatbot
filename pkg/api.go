// Created by zhbinary on 2023/8/20.
package pkg

import "fmt"

func SaveReviewByOrder(feedback *Feedback, sendId string) bool {
	fmt.Println("saveReviewByOrder")
	// send msg to kafka, then another service will consume the msg and save into DB

	// Failed to send to kafka, return false, it occurs while network busy

	// success to send to kafka
	return true
}

func ResendReply(response string, senderPsid string) {
	fmt.Println("resendReply")
	// Send delay message to kafka or pulsar, then another service will consumer the msg and retry
}
