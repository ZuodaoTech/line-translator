package assistant

import (
	"context"
	"fmt"

	"github.com/zuodaotech/line-translator/common/ai"
)

func (a *Assistant) Translate(ctx context.Context, content, lang string) (string, string, error) {
	inst1 := fmt.Sprintf(`You are an expert linguist, specializing in translation to %s language.
Please provide the {target_lang} translation for above text.
Do not provide any explanations or text apart from the translation.
`, lang)

	// first round
	ret, err := a.aiInst.MultipleSteps(ctx, ai.ChainParams{
		Format: "text",
		Steps: []ai.ChainParamsStep{
			{Input: content},
			{Instruction: inst1},
		},
	})
	if err != nil {
		return "", "", err
	}

	round2Text := fmt.Sprintf(`
Your task is to carefully read a source text and a translation to %s, and then give constructive criticism and helpful suggestions to improve the translation.
The final style and tone of the translation should match the style of %s and source's original writing style.

The source text and initial translation, delimited by XML tags <SOURCE_TEXT></SOURCE_TEXT> and <TRANSLATION></TRANSLATION>, are as follows:
<SOURCE_TEXT>
%s
</SOURCE_TEXT>

<TRANSLATION>
%s
</TRANSLATION>

When writing suggestions, pay attention to whether there are ways to improve the translation's
(i) accuracy (by correcting errors of addition, mistranslation, omission, or untranslated text),
(ii) fluency (by applying {target_lang} grammar, spelling and punctuation rules, and ensuring there are no unnecessary repetitions),
(iii) style (by ensuring the translations reflect the style of the source text and takes into account any cultural context),
(iv) terminology (by ensuring terminology use is consistent and reflects the source text domain; and by only ensuring you use equivalent idioms {target_lang}).

Write a list of specific, helpful and constructive suggestions for improving the translation.
Each suggestion should address one specific part of the translation.
Output only the suggestions and nothing else.`, lang, lang, content, ret.Text)

	// 2nd round
	suggestion, err := a.aiInst.OneTimeRequest(ctx, round2Text)
	if err != nil {
		return "", "", err
	}

	round3Text := fmt.Sprintf(`
Your task is to carefully read, then edit, a translation to %s language, taking into
account a list of expert suggestions and constructive criticisms.

The source text, the initial translation, and the expert linguist suggestions are delimited by XML tags <SOURCE_TEXT></SOURCE_TEXT>, <TRANSLATION></TRANSLATION> and <EXPERT_SUGGESTIONS></EXPERT_SUGGESTIONS>
as follows:

<SOURCE_TEXT>
%s
</SOURCE_TEXT>

<TRANSLATION>
%s
</TRANSLATION>

<EXPERT_SUGGESTIONS>
%s
</EXPERT_SUGGESTIONS>

Please take into account the expert suggestions when editing the translation. Edit the translation by ensuring:

(i) accuracy (by correcting errors of addition, mistranslation, omission, or untranslated text),
(ii) fluency (by applying %s grammar, spelling and punctuation rules and ensuring there are no unnecessary repetitions),
(iii) style (by ensuring the translations reflect the style of the source text)
(iv) terminology (inappropriate for context, inconsistent use), or
(v) other errors.

Output only the new translation and nothing else`, lang, content, ret.Text, suggestion, lang)

	// 3rd round
	result, err := a.aiInst.OneTimeRequest(ctx, round3Text)
	if err != nil {
		return "", "", err
	}

	return ret.Text, result, nil
}
