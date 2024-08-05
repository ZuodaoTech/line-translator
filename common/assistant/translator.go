package assistant

import (
	"context"
	"fmt"
	"time"

	"github.com/zuodaotech/line-translator/common/ai"
)

func (a *Assistant) Translate(ctx context.Context, content, srcLang, dstLang string) (string, error) {
	zh2JaExample := `Example:
我喜欢学习新语言。
{"output": "私は新しい言語を学ぶのが好きです。"}
`

	ja2ZhExample := `Example:
私は新しい言語を学ぶのが好きです。
{"output": "我喜欢学习新语言。"}
`

	example := ""
	if srcLang == "zh" && dstLang == "ja" {
		example = zh2JaExample
	} else if srcLang == "ja" && dstLang == "zh" {
		example = ja2ZhExample
	}

	inst1 := fmt.Sprintf(`You are an expert linguist, specializing in translation to %s and %s language.
Please provide the %s translation for above text.

%s

Please always return JSON format with the translated text.
`, srcLang, dstLang, dstLang, example)

	// set timeout for 60s
	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()
	ret, err := a.aiInst.MultipleSteps(ctx, ai.ChainParams{
		Format: "json",
		Steps: []ai.ChainParamsStep{
			{Input: content},
			{Instruction: inst1},
		},
	})
	if err != nil {
		return "", err
	}

	result, ok := ret.Json["output"].(string)
	if !ok {
		result = "Failed to Translate"
	}

	return result, nil
}
