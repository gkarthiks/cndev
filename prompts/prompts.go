package prompts

import (
	"fmt"
	"github.com/gkarthiks/cndev/utils"
	"github.com/manifoldco/promptui"
	"os"
)

// PromptYesNo prompts Yes/No question
func PromptYesNo(label interface{}) (promptResult string) {
	items := []string{utils.StringYes, utils.StringNo}
	index := -1
	var err error
	for index < 0 {
		prompt := promptui.Select{
			Label: label,
			Items: items,
		}
		index, promptResult, err = prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			os.Exit(1)
		}
	}
	return
}

// PromptUser prompts the user with the provided questions
func PromptUser(label interface{}, validate func(input string) error) string {
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}
	prompt := promptui.Prompt{
		Label:     label,
		Templates: templates,
		Validate:  validate,
	}
	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Input: %s\n", result)
	return result
}
