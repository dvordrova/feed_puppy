package layout

import (
	tele "gopkg.in/telebot.v3"
)

type LocaleFunc func(*tele.User) string

func (lt *Layout) Middleware(defaultLocale string, localeFunc ...LocaleFunc) tele.MiddlewareFunc {
	var f LocaleFunc
	if len(localeFunc) > 0 {
		f = localeFunc[0]
	}

	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			locale := defaultLocale
			if f != nil {
				if l := f(c.Sender()); l != "" {
					locale = l
				}
			}

			lt.SetLocale(c, locale)

			defer func() {
				lt.mu.Lock()
				delete(lt.ctxs, c)
				lt.mu.Unlock()
			}()

			return next(c)
		}
	}
}
