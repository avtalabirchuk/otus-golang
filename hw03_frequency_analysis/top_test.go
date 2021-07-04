package hw03frequencyanalysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Change to true if needed.
var taskWithAsteriskIsCompleted = false

var (
	text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

	enText = `The three little pigs are brothers. 
	They are going into the forest. 
	They want to build three houses. 
	"Let’s build our houses here," 
	says the first little pig, Percy. 
	"Yes," says the second little pig, Peter. 
	"That’s a good idea," says the third little pig, Patrick. 
	The first little pig, Percy, gets some straw he starts to build a house of straw. 
	He sings, "Hum de hum, dum de dum, hee de dum, dee de hum," when he works. 
	The second little pig, Peter, gets some wood and he starts to build a house of wood. 
	He sings, "Hum de hum, dum de dum, hee de dum, dee de hum," when he works. 
	The third pig, Patrick, is very clever. 
	He gets some bricks and he starts build a house of bricks. 
	He sings, "Hum de hum, dum de dum, hee de dum, dee de hum," when he works. 
	Now all the houses are ready. 
	The three little pigs make a fence and they paint it red. he
    `
	numeric = ` 18 20 12312 10 2 19 20 240 29 10 29 1029 229 229
	`
)

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		require.Len(t, Top10(""), 0)
	})

	t.Run("positive test", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{
				"а",         // 8
				"он",        // 8
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"в",         // 4
				"его",       // 4
				"если",      // 4
				"кристофер", // 4
				"не",        // 4
			}
			require.Equal(t, expected, Top10(text))
		} else {
			expected := []string{
				"он",        // 8
				"а",         // 6
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"-",         // 4
				"Кристофер", // 4
				"если",      // 4
				"не",        // 4
				"то",        // 4
			}
			require.ElementsMatch(t, expected, Top10(text))
		}
	})
	t.Run("English test", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{"de", "he", "little", "pig,", "dum,", "a", "the", "build", "The", "He"}
			require.Equal(t, expected, Top10(enText))
		} else {
			expected := []string{"de", "he", "little", "pig,", "dum,", "a", "the", "build", "The", "He"}
			require.ElementsMatch(t, expected, Top10(enText))
		}
	})

	t.Run("Numeric text", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{"20", "10", "229", "29", "1029", "240", "2", "18", "12312", "19"}
			require.Equal(t, expected, Top10(enText))
		} else {
			expected := []string{"229", "20", "10", "29", "1029", "240", "2", "18", "12312", "19"}
			require.ElementsMatch(t, expected, Top10(numeric))
		}
	})
}
