package vcode

import "github.com/mojocn/base64Captcha"

var (
	store = base64Captcha.DefaultMemStore
	mtp   = base64Captcha.DriverMath{Width: 200, Height: 100}
)

func Get() (id, b64s string, err error) {
	driver := mtp.ConvertFonts()
	c := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err = c.Generate()
	return
}

func Verify(id string, answer string) (match bool) {
	if id == "" || answer == "" {
		return
	}
	driver := mtp.ConvertFonts()
	c := base64Captcha.NewCaptcha(driver, store)
	return c.Verify(id, answer, true)
}
