package assistant

import (
	"context"

	"github.com/zuodaotech/line-translator/common/ai"
)

func (a *Assistant) Translate(ctx context.Context, content string) (string, error) {
	inst1 := `You are an expert linguist, specializing in translation to Chinese and Japanese language.
Please provide the translation for above text.
If the text is Chinese, translate it to Japanese. If the text is Japanese, translate it to Chinese.
Do not provide any explanations or text apart from the translation.
`

	ret, err := a.aiInst.MultipleSteps(ctx, ai.ChainParams{
		Format: "text",
		Steps: []ai.ChainParamsStep{
			{Input: content},
			{Instruction: inst1},
		},
	})
	if err != nil {
		return "", err
	}

	return ret.Text, nil
}
