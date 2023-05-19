package main

import (
	"github.com/gofiber/fiber/v2"
	"encoding/json"
	"fmt"

	gogpt "github.com/tapp-ai/go-openai"
	"go.uber.org/zap"
)

type GetDomains struct {
	Name string `json:"name"`
	ExtraContext string `json:"context"`
}

//	 GetDomains gets a list of available domains from a given business name
//	 {
//			"name": "name of business"
//	 }
func (a *App) GetDomains(c *fiber.Ctx) error {

	// parse the object into the struct
	input := GetDomains{}
	err := c.BodyParser(&input)
	if err != nil {
		a.Log.Error("error parsing POST body into the struct")
		return c.JSON(ErrorResponse("error parsing POST body into the struct"))
	}

	var contextStr string

	if len(input.ExtraContext) > 0 {
		contextStr = fmt.Sprintf(`(here is some context: %s)`, input.ExtraContext)
		} else {
		contextStr = ""
	}

	// write a prompt
	prompt := fmt.Sprintf(`Return to me an ordered list of available domain names for the following business name%s.
	Input: %s
	Output:`, contextStr, input.Name)

	fmt.Println("PROMPT HERRE ++++++\n\n", prompt, contextStr, "\n\n")

	// call the chat completion
	resp, err := a.GptClient.CreateChatCompletion(
		c.Context(),
		gogpt.ChatCompletionRequest{
			Model: gogpt.GPT3Dot5Turbo,
			Messages: []gogpt.ChatCompletionMessage{
				{
					Role:    gogpt.ChatMessageRoleSystem,
					Content: "You're a chat bot that's only role is to return domain names",
				},
				{
					Role:    gogpt.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: 400,
		},
	)

	if err != nil {
		a.Log.Error("Error in domain name chat bot", zap.Error(err))
		return c.JSON(ErrorResponse("Error in domain name chat bot"))
	}

	// this line is for debugging the response
	// choices, err := json.Marshal(resp)
	choices, err := json.Marshal(resp)
	if err != nil {
		a.Log.Error("Failed to marshal response", zap.Error(err))
		return c.JSON(ErrorResponse("Failed to marshal response"))
	}

	// fmt.Printf("Resp:  ++++\n", resp.Choices)
	// fmt.Println("\n\nHERE CHOICES++++\n", choices)

	a.Log.Info(string(choices))

	return c.JSON(SuccessResponse(resp.Choices[0].Message.Content))

}
