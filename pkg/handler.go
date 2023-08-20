// Created by zhbinary on 2023/8/19.
package pkg

import "strings"

func HandleFeedback(senderPsid string, feedback *Feedback) bool {
	// Get rate and comment
	return SaveReviewByOrder(feedback, senderPsid)
}

func HandleMessage(senderPsid string, receivedMessage *Message) interface{} {
	var response interface{}

	// Checks if the message contains text
	if receivedMessage.Text != "" {
		// Create the payload for a basic text message, which
		// will be added to the body of your request to the Send API
		if strings.Contains(receivedMessage.Text, "review") {
			response = map[string]interface{}{
				"attachment": map[string]interface{}{
					"type": "template",
					"payload": map[string]interface{}{
						"template_type": "customer_feedback",
						"title":         "Rate your experience with Original Coast Clothing.",                             // Business needs to define.
						"subtitle":      "Let Original Coast Clothing know how they are doing by answering two questions", // Business needs to define.
						"button_title":  "Rate Experience",                                                                // Business needs to define.
						"feedback_screens": []map[string]interface{}{
							{
								"questions": []map[string]interface{}{
									{
										"id":           "hauydmns8", // Unique id for question that business sets
										"type":         "csat",
										"title":        "How would you rate your experience with Original Coast Clothing?", // Optional. If business does not define, we show standard text. Standard text based on question type ("csat", "nps", "ces" >>> "text")
										"score_label":  "neg_pos",                                                          // Optional
										"score_option": "five_stars",                                                       // Optional
										"follow_up": map[string]interface{}{ // Optional. Inherits the title and id from the previous question on the same page.  Only free-from input is allowed. No other title will show.
											"type":        "free_form",
											"placeholder": "Give additional feedback", // Optional
										},
									},
								},
							},
						},
						"business_privacy": map[string]interface{}{
							"url": "https://www.example.com",
						},
						"expires_in_days": 3, // Optional, default 1 day, business defines 1-7 days
					},
				},
			}
		} else {
			response = map[string]interface{}{
				"text": `Input 'review' to give a comment to order`,
			}
			if receivedMessage.Nlp != nil && len(receivedMessage.Nlp.Traits.Sentiment) > 0 {
				sentimentTrait := receivedMessage.Nlp.Traits.Sentiment[0]
				if sentimentTrait.Value == "negative" {
					response = map[string]interface{}{
						"text": `Sorry sir, what can I do for you.`,
					}
				}
			}
		}
	}
	return response
}
