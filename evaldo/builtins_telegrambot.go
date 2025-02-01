//go:build !no_telegram
// +build !no_telegram

package evaldo

import (
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/refaktor/rye/env"

	tgm "github.com/go-telegram-bot-api/telegram-bot-api"
)

/*

	new-telegram-bot "MyRyeBot" :bot
	|on-updates-do {
		-> 'Message :msg |if {
			print msg -> 'From -> 'UserName
			print concat " " msg -> Text

			new-telegram-message msg -> 'Chat -> 'ID  msg -> 'Text
			|send* bot
		}
	}

	new-telegram-bot 1
	telegra-bot//on-updates-do
	new-telegram-message
	telegram-bot//send

*/

func TelegramUpdateToRyeDict(update_ tgm.Update) env.Dict {
	update := make(map[string]any)
	if update_.Message != nil {
		msg := make(map[string]any)
		msg["Text"] = update_.Message.Text
		msg["MessageID"] = update_.Message.MessageID
		if update_.Message.From != nil {
			from := make(map[string]any)
			from["UserName"] = update_.Message.From.UserName
			from["FirstName"] = update_.Message.From.FirstName
			from["ID"] = update_.Message.From.ID
			from["LastName"] = update_.Message.From.LastName
			from["IsBot"] = update_.Message.From.IsBot

			msg["From"] = *env.NewDict(from)
		}
		if update_.Message.Chat != nil {
			chat := make(map[string]any)
			chat["ID"] = update_.Message.Chat.ID
			msg["Chat"] = *env.NewDict(chat)
		}
		update["Message"] = *env.NewDict(msg)
	}
	return *env.NewDict(update)
}

var Builtins_telegrambot = map[string]*env.Builtin{

	"bot": {
		Argsn: 1,
		Doc:   "Create new telegram bot using API value.",
		Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			switch arg := arg0.(type) {
			case env.String:
				bot, err := tgm.NewBotAPI(arg.Value)
				if err != nil {
					return MakeBuiltinError(ps, "Error in NewBotAPI function.", "new-telegram-bot")
				}
				return *env.NewNative(ps.Idx, bot, "telegram-bot")
			default:
				return MakeArgError(ps, 1, []env.Type{env.StringType}, "new-telegram-bot")
			}
		},
	},

	"telegram-bot//on-update": {
		Argsn: 2,
		Doc:   "Get telegram update and add to Rye dictionary",
		Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			switch bot := arg0.(type) {
			case env.Native:
				switch code := arg1.(type) {
				case env.Block:
					u := tgm.NewUpdate(0)
					u.Timeout = 60

					updates, _ := bot.Value.(*tgm.BotAPI).GetUpdatesChan(u)
					ser := ps.Ser
					ps.Ser = code.Series

					for update := range updates {
						dict := TelegramUpdateToRyeDict(update)
						ps = EvalBlockInj(ps, dict, true)
						if ps.ErrorFlag {
							return ps.Res
						}
						ps.Ser.Reset()
					}
					ps.Ser = ser
					return env.NewInteger(1)
				case env.Function:
					u := tgm.NewUpdate(0)
					u.Timeout = 60

					updates, _ := bot.Value.(*tgm.BotAPI).GetUpdatesChan(u)

					for update := range updates {
						dict := TelegramUpdateToRyeDict(update)
						ps.FailureFlag = false
						ps.ErrorFlag = false
						ps.ReturnFlag = false
						psTemp := env.ProgramState{}
						err := copier.Copy(&psTemp, &ps)
						if err != nil {
							fmt.Println(err.Error())
							// TODO return makeError(ps, err.Error())
						}
						CallFunctionArgsN(code, ps, nil, dict)
					}
					return env.NewInteger(1)

				default:
					return MakeArgError(ps, 2, []env.Type{env.BlockType}, "telegram-bot//on-update")
				}
			default:
				return MakeArgError(ps, 1, []env.Type{env.NativeType}, "telegram-bot//on-update")
			}
		},
	},

	"telegram-message//send": {
		Argsn: 2,
		Doc:   "Send message in telegram bot.",
		Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			switch msg := arg0.(type) {
			case env.Native:
				switch bot := arg1.(type) {
				case env.Native:
					_, err := bot.Value.(*tgm.BotAPI).Send(msg.Value.(*tgm.MessageConfig))
					if err != nil {
						return makeError(ps, err.Error())
					}
					return arg0
				default:
					return MakeArgError(ps, 2, []env.Type{env.NativeType}, "telegram-message//send")
				}
			default:
				return MakeArgError(ps, 1, []env.Type{env.NativeType}, "telegram-message//send")
			}

		},
	},

	"telegram-message//parse-mode!": {
		Argsn: 2,
		Doc:   "Set the parse mode of the message",
		Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			switch msg := arg0.(type) {
			case env.Native:
				switch mode := arg1.(type) {
				case env.String:
					message, _ := msg.Value.(tgm.MessageConfig)
					message.ParseMode = mode.Value
					return msg
				default:
					return MakeArgError(ps, 2, []env.Type{env.NativeType}, "telegram-message//send")
				}
			default:
				return MakeArgError(ps, 1, []env.Type{env.NativeType}, "telegram-message//send")
			}

		},
	},

	"message": {
		Argsn: 2,
		Doc:   "Create new telegram bot message.",
		Fn: func(ps *env.ProgramState, arg0 env.Object, arg1 env.Object, arg2 env.Object, arg3 env.Object, arg4 env.Object) env.Object {
			switch cid := arg0.(type) {
			case env.Integer:
				switch txt := arg1.(type) {
				case env.String:
					msg := tgm.NewMessage(cid.Value, txt.Value)
					msg.ParseMode = "HTML"
					return *env.NewNative(ps.Idx, &msg, "telegram-message")
				default:
					return MakeArgError(ps, 2, []env.Type{env.StringType}, "new-telegram-message")
				}
			default:
				return MakeArgError(ps, 1, []env.Type{env.IntegerType}, "new-telegram-message")
			}
		},
	},
}
